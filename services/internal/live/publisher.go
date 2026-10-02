package live

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
)

// Publisher turns Flow updates of one match into contract messages for the hub
// and keeps the match snapshot in the store up to date.
// It runs in a single goroutine, so seq and history need no locking.
type Publisher struct {
	hub     *Hub
	store   *Store
	matchID string
	params  flow.Params
	seq     int64
	period  int
	at      time.Duration // match time of the last update
	stats   *Stats
	history map[int]FlowValues // Flow at the end of each match minute

	// The outcome model, nil without one, and what it knows about the match
	model     *predict.Model
	rating    float64
	situation predict.Situation
	sent      Probabilities // the chances viewers have, to send only a change

}

func NewPublisher(hub *Hub, store *Store, matchID string, params flow.Params) *Publisher {
	return &Publisher{hub: hub, store: store, matchID: matchID, params: params, history: map[int]FlowValues{}}
}

// UseModel makes the Publisher predict the outcome with m (docs/PREDICTION.md).
func (p *Publisher) UseModel(m *predict.Model) {
	p.model = m
}

// Schedule registers the match before kick-off, with the pre-match chances,
// so the page can load it before anyone starts the replay.
func (p *Publisher) Schedule(match event.Match) {
	p.store.Schedule(match)
	p.preMatch(match)
}

// Start puts the match at kick-off. seq keeps growing across replays, so a
// message of the previous replay that a client still holds counts as old.
func (p *Publisher) Start(match event.Match) {
	clear(p.history) // minutes of the previous replay
	p.period = 1
	p.at = 0
	p.stats = NewStats()
	p.store.Start(match, p.seq)
	p.preMatch(match)
}

// preMatch puts the chances before kick-off into the snapshot: they are both
// the current chances and the "before the match" line of the card.
func (p *Publisher) preMatch(match event.Match) {
	if p.model == nil {
		return
	}
	p.rating = p.model.Rating(match.Home.Name, match.Away.Name)
	p.situation = predict.Situation{Period: 1}
	p.sent = toPercents(p.model.Predict(p.rating, p.situation))
	prediction := Prediction{Current: p.sent, PreMatch: p.sent, Model: p.model.Name}
	p.store.update(p.matchID, func(s *Snapshot) { s.Prediction = &prediction })
}

// Run publishes the match started with Start until updates is closed.
func (p *Publisher) Run(updates <-chan flow.Update) {
	for u := range updates {
		if u.Cause != nil && u.Cause.Period != p.period {
			if p.period == 1 {
				p.whistle(event.Halftime, "halftime")
			}
			p.period = u.Cause.Period
		}

		if u.Cause != nil && !u.Cause.Type.Internal() {
			me := p.matchEvent(*u.Cause)
			p.emit("match.event", me, func(s *Snapshot) {
				s.Events = append(s.Events, me)
				if me.Type == event.Goal {
					addGoal(&s.Score, *me.Side)
				}
			})
		}

		if u.Cause != nil && p.stats.Add(*u.Cause) {
			ms := MatchStats{Stats: p.stats.Rows()}
			p.emit("match.stats", ms, func(s *Snapshot) { s.Stats = ms.Stats })
		}

		if u.Cause != nil {
			p.count(*u.Cause)
		}
		p.situation.Period = p.period
		p.situation.At = u.At
		p.predict()

		fu := p.flowUpdate(u)
		p.emit("flow.update", fu, func(s *Snapshot) {
			s.Flow = fu.Current
			s.Delta10 = fu.Delta10
			s.Factors = fu.Factors
			s.Points = upsertPoint(s.Points, fu.Point)
			s.At = u.At
			s.Period = p.period
			if s.Status == "halftime" && p.period > 1 {
				s.Status = "live" // the second half has kicked off
			}
		})
		p.at = u.At
	}

	p.whistle(event.Fulltime, "finished")
	p.situation.Finished = true // the result is known: 100% for it
	p.predict()
}

// count keeps the goals and red cards the outcome model needs.
func (p *Publisher) count(e event.Event) {
	s := &p.situation
	switch {
	case e.Type == event.Goal && e.Side == event.Home:
		s.HomeGoals++
	case e.Type == event.Goal && e.Side == event.Away:
		s.AwayGoals++
	case e.Type == event.RedCard && e.Side == event.Home:
		s.HomeReds++
	case e.Type == event.RedCard && e.Side == event.Away:
		s.AwayReds++
	}
}

// predict sends the chances when a whole percent of them has changed. The
// model is cheap, but a message every tick with the same numbers is not.
func (p *Publisher) predict() {
	if p.model == nil {
		return
	}
	now := toPercents(p.model.Predict(p.rating, p.situation))
	if now == p.sent {
		return
	}
	p.sent = now
	p.emit("prediction.update", now, func(s *Snapshot) {
		if s.Prediction != nil {
			s.Prediction.Current = now
		}
	})
}

// whistle ends a half or the match: a timeline event without a side that
// also changes the match status, so viewers learn it without a reload.
func (p *Publisher) whistle(t event.Type, status string) {
	minute, added := minuteOf(p.period, p.at)
	me := MatchEvent{ID: p.matchID + ":" + string(t), Type: t, Minute: minute, AddedTime: added}
	p.emit("match.event", me, func(s *Snapshot) {
		s.Events = append(s.Events, me)
		s.Status = status
		if t == event.Halftime {
			score := s.Score
			s.Halftime = &score
		}
	})
}

func (p *Publisher) flowUpdate(u flow.Update) FlowUpdate {
	minute := int(u.At.Minutes()) + 1
	current := FlowValues{Home: u.Home, Away: u.Away}
	p.history[minute] = current

	before := p.history[minute-10] // zero before the 11th minute
	return FlowUpdate{
		Current: current,
		Delta10: FlowValues{Home: current.Home - before.Home, Away: current.Away - before.Away},
		Point:   FlowPoint{Minute: minute, Home: current.Home, Away: current.Away},
		Clock:   ClockOf(p.period, u.At, time.Now().UTC()),
		Factors: toFactors(u.Factors),
	}
}

func (p *Publisher) matchEvent(e event.Event) MatchEvent {
	minute, added := minuteOf(e.Period, e.Elapsed)
	_, impact := p.params.Weight(e)

	me := MatchEvent{
		ID:         e.ID,
		Type:       e.Type,
		Side:       &e.Side,
		Minute:     minute,
		AddedTime:  added,
		XG:         e.XG,
		Position:   e.Pos,
		FlowImpact: &impact,
	}
	if e.Player != "" {
		me.Player = &PersonRef{ID: e.PlayerID, Name: e.Player}
	}
	return me
}

// emit records the change in the store first and only then sends the message,
// so a snapshot is never behind a message a client has already received.
func (p *Publisher) emit(kind string, data any, apply func(*Snapshot)) {
	p.seq++
	seq := p.seq
	p.store.update(p.matchID, func(s *Snapshot) {
		apply(s)
		s.Seq = seq
	})

	msg, err := json.Marshal(Message{
		Type:    kind,
		MatchID: p.matchID,
		Seq:     seq,
		SentAt:  time.Now().UTC(),
		Data:    data,
	})
	if err != nil {
		log.Printf("marshal %s: %v", kind, err)
		return
	}
	p.hub.Publish(p.matchID, msg)
}

func addGoal(score *Score, side event.Side) {
	if side == event.Home {
		score.Home++
	} else {
		score.Away++
	}
}
