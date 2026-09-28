package live

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
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
	at      time.Duration      // match time of the last update
	history map[int]FlowValues // Flow at the end of each match minute
}

func NewPublisher(hub *Hub, store *Store, matchID string, params flow.Params) *Publisher {
	return &Publisher{hub: hub, store: store, matchID: matchID, params: params, history: map[int]FlowValues{}}
}

// Start puts the match at kick-off. seq keeps growing across replays, so a
// message of the previous replay that a client still holds counts as old.
func (p *Publisher) Start(match event.Match) {
	clear(p.history) // minutes of the previous replay
	p.period = 1
	p.at = 0
	p.store.Start(match, p.seq)
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

		if u.Cause != nil && u.Cause.Type != event.Possession {
			me := p.matchEvent(*u.Cause)
			p.emit("match.event", me, func(s *Snapshot) {
				s.Events = append(s.Events, me)
				if me.Type == event.Goal {
					addGoal(&s.Score, *me.Side)
				}
			})
		}

		fu := p.flowUpdate(u)
		p.emit("flow.update", fu, func(s *Snapshot) {
			s.Flow = fu.Current
			s.Delta10 = fu.Delta10
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
