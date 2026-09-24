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
}

func NewState(p Params) *State {
	return &State{p: p, impulse: map[event.Side]float64{}}
}

// Apply moves match time to the event and adds its weight.
func (s *State) Apply(e event.Event) {
	if s.period != 0 && e.Period > s.period {
		for side := range s.impulse {
			s.impulse[side] *= s.p.HalftimeKeep
		}
		s.at = e.Elapsed // the second half clock restarts at 45:00
	}
	s.period = e.Period
	s.Advance(e.Elapsed)

	side, w := s.p.weight(e)
	s.impulse[side] += w
}

// Advance lets the impulse decay up to match time t. Time never goes back.
func (s *State) Advance(t time.Duration) {
	if t <= s.at {
		return
	}
	decay := math.Exp(-(t - s.at).Minutes() / s.p.Tau)
	for side := range s.impulse {
		s.impulse[side] *= decay
	}
	s.at = t
}

// Flow returns both teams' Flow on the 0–100 scale.
func (s *State) Flow() (home, away float64) {
	return s.flow(event.Home), s.flow(event.Away)
}

func (s *State) flow(side event.Side) float64 {
	return 100 * (1 - math.Exp(-math.Max(0, s.impulse[side])/s.p.K))
}
