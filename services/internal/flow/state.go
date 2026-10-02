package flow

import (
	"math"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// State holds the impulse S of both teams (docs/FLOW.md, section 3).
// Decay is exponential, so instead of summing every past event
// we keep one number per team and shrink it as match time passes.
type State struct {
	p       Params
	impulse map[event.Side]float64
	at      time.Duration
	period  int

	// The same impulse split by group of events, for the explanation (section 6).
	// Every change of the impulse goes to its group too, so the parts of a
	// team always add up to its impulse.
	parts  map[event.Side]map[Group]float64
	recent []mark
}

func NewState(p Params) *State {
	return &State{
		p:       p,
		impulse: map[event.Side]float64{},
		parts:   map[event.Side]map[Group]float64{event.Home: {}, event.Away: {}},
	}
}

// StartPeriod moves on to the next period, which kicks off at match time at.
// Part of the impulse is lost over the break (docs/FLOW.md, section 3.1).
// Calling it again for the period already under way changes nothing.
func (s *State) StartPeriod(period int, at time.Duration) {
	if s.period != 0 && period > s.period {
		s.scale(s.p.HalftimeKeep)
		s.at = at // the second half clock restarts at 45:00
		s.recent = nil
	}
	s.period = period
}

// Apply moves match time to the event and adds its weight.
func (s *State) Apply(e event.Event) {
	s.StartPeriod(e.Period, e.Elapsed)
	s.Advance(e.Elapsed)

	if e.Type == event.Possession && e.HomeShare != nil {
		lead := *e.HomeShare - 0.5 // +0.2 means home had the ball 70% of the minute
		s.add(event.Home, Possession, s.p.Possession*lead)
		s.add(event.Away, Possession, -s.p.Possession*lead)
		s.remember(mark{side: event.Home, group: Possession, at: e.Elapsed, share: *e.HomeShare})
		s.remember(mark{side: event.Away, group: Possession, at: e.Elapsed, share: 1 - *e.HomeShare})
		return
	}

	side, w := s.p.Weight(e)
	if w == 0 {
		return // fouls, tackles: nothing to add and nothing to explain
	}
	g := groupOf(e.Type)
	s.add(side, g, w)
	s.remember(mark{side: side, group: g, at: e.Elapsed})
}

// add puts a weight into the team's impulse. A negative weight only dampens
// what the team has built up and never leaves a debt: otherwise a long spell
// without the ball would hide the team's next shots until the debt is paid off.
// The group gets the change that really happened, the weight after the floor.
func (s *State) add(side event.Side, g Group, w float64) {
	before := s.impulse[side]
	s.impulse[side] = math.Max(0, before+w)
	s.parts[side][g] += s.impulse[side] - before
}

// Advance lets the impulse decay up to match time t. Time never goes back.
func (s *State) Advance(t time.Duration) {
	if t <= s.at {
		return
	}
	s.scale(math.Exp(-(t - s.at).Minutes() / s.p.Tau))
	s.at = t
}

// scale shrinks the impulse and every part of it by the same factor.
func (s *State) scale(f float64) {
	for side := range s.impulse {
		s.impulse[side] *= f
	}
	for _, parts := range s.parts {
		for g := range parts {
			parts[g] *= f
		}
	}
}

// At returns the match time the state has been advanced to
func (s *State) At() time.Duration {
	return s.at
}

// Flow returns both teams' Flow on the 0–100 scale. A team on the pitch always
// has the base impulse, so its Flow never falls to zero.
func (s *State) Flow() (home, away float64) {
	return s.flow(event.Home), s.flow(event.Away)
}

func (s *State) flow(side event.Side) float64 {
	return 100 * (1 - math.Exp(-(s.p.Base+s.impulse[side])/s.p.K))
}
