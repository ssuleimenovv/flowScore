package live

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
)

// Publisher turns Flow updates of one match into contract messages for the hub.
// It runs in a single goroutine, so seq and history need no locking.
type Publisher struct {
	hub     *Hub
	matchID string
	params  flow.Params
	seq     int64
	history map[int]FlowValues // Flow at the end of each match minute
}

func NewPublisher(hub *Hub, matchID string, params flow.Params) *Publisher {
	return &Publisher{hub: hub, matchID: matchID, params: params, history: map[int]FlowValues{}}
}

// Run publishes until updates is closed.
func (p *Publisher) Run(updates <-chan flow.Update) {
	for u := range updates {
		if u.Cause != nil && u.Cause.Type != event.Possession {
			p.send("match.event", p.matchEvent(*u.Cause))
		}
		p.send("flow.update", p.flowUpdate(u))
	}
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
	}
}

func (p *Publisher) matchEvent(e event.Event) MatchEvent {
	minute, added := minuteOf(e.Period, e.Elapsed)
	_, impact := p.params.Weight(e)

	me := MatchEvent{
		ID:         e.ID,
		Type:       e.Type,
		Side:       e.Side,
		Minute:     minute,
		AddedTime:  added,
		XG:         e.XG,
		Position:   e.Pos,
		FlowImpact: &impact,
	}
	if e.Player != "" {
		me.Player = &PersonRef{Name: e.Player}
	}
	return me
}

func (p *Publisher) send(kind string, data any) {
	p.seq++
	msg, err := json.Marshal(Message{
		Type:    kind,
		MatchID: p.matchID,
		Seq:     p.seq,
		SentAt:  time.Now().UTC(),
		Data:    data,
	})
	if err != nil {
		log.Printf("marshal %s: %v", kind, err)
		return
	}
	p.hub.Publish(p.matchID, msg)
}
