package flow

import (
	"math"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func ev(t event.Type, side event.Side, period, minute int) event.Event {
	return event.Event{Type: t, Side: side, Period: period, Elapsed: time.Duration(minute) * time.Minute}
}

func near(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.01 {
		t.Errorf("%s = %.4f, want %.4f", name, got, want)
	}
}

func TestGoal(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.Goal, event.Home, 1, 10))

	home, away := s.Flow()
	near(t, "home", home, 74.15)
	near(t, "away", away, 0)
}

func TestDecayAfterTau(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.Goal, event.Home, 1, 10))
	s.Advance(15 * time.Minute) // one tau later the impulse is 1/e of its value

	home, _ := s.Flow()
	near(t, "home", home, 39.21)
}

func TestYellowCardHelpsOpponent(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.YellowCard, event.Away, 1, 20))

	home, away := s.Flow()
	if home <= 0 || away != 0 {
		t.Errorf("home = %.2f, away = %.2f; want home > 0, away = 0", home, away)
	}
}

func TestHalftimeKeepsPartOfImpulse(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.Goal, event.Home, 1, 45))
	s.Apply(ev(event.Foul, event.Away, 2, 45)) // zero weight, only switches the period

	home, _ := s.Flow()
	near(t, "home", home, 61.21)
}

func TestNegativeImpulseShowsZero(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.Substitution, event.Home, 1, 60))

	home, _ := s.Flow()
	near(t, "home", home, 0)
}
