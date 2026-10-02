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

// floor is the Flow of a team with only the base impulse: 100 · (1 − e^(−3/18.7))
const floor = 14.82

func TestIdleTeamsStayAtFloor(t *testing.T) {
	s := NewState(DefaultParams())

	home, away := s.Flow()
	near(t, "home", home, floor)
	near(t, "away", away, floor)
}

func TestMockupCalibration(t *testing.T) {
	// The Explainability factors on the Match board add up to S = 29 and City shows 82
	s := NewState(DefaultParams())
	s.Apply(ev(event.Goal, event.Home, 1, 10))
	s.impulse[event.Home] = 29

	home, _ := s.Flow()
	near(t, "home", math.Round(home), 82)
}

func TestGoal(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.Goal, event.Home, 1, 10))

	home, away := s.Flow()
	near(t, "home", home, 75.10)
	near(t, "away", away, floor)
}

func TestDecayAfterTau(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.Goal, event.Home, 1, 10))
	s.Advance(15 * time.Minute) // one tau later the impulse is 1/e of its value

	home, _ := s.Flow()
	near(t, "home", home, 45.82)
}

func TestYellowCardHelpsOpponent(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.YellowCard, event.Away, 1, 20))

	home, away := s.Flow()
	near(t, "home", home, 34.81)
	near(t, "away", away, floor)
}

func TestHalftimeKeepsPartOfImpulse(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.Goal, event.Home, 1, 45))
	s.Apply(ev(event.Foul, event.Away, 2, 45)) // zero weight, only switches the period

	home, _ := s.Flow()
	near(t, "home", home, 63.99)
}

func TestNegativeWeightStopsAtFloor(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.RedCard, event.Home, 1, 60))

	home, _ := s.Flow()
	near(t, "home", home, floor)
}

func TestPossessionLeadHelpsHolder(t *testing.T) {
	s := NewState(DefaultParams())
	share := 0.7
	s.Apply(event.Event{Type: event.Possession, Period: 1, Elapsed: time.Minute, HomeShare: &share})

	home, away := s.Flow()
	near(t, "home", home, 22.64)
	near(t, "away", away, floor)
}

func TestNoDebtAfterLongSpellWithoutBall(t *testing.T) {
	// Away has the ball 30% of the time for 20 minutes, then shoots.
	// The shot must count in full, as if the spell had never happened.
	s := NewState(DefaultParams())
	share := 0.7
	for m := 1; m <= 20; m++ {
		s.Apply(event.Event{Type: event.Possession, Period: 1, Elapsed: time.Duration(m) * time.Minute, HomeShare: &share})
	}
	s.Apply(ev(event.ShotOnTarget, event.Away, 1, 20))

	fresh := NewState(DefaultParams())
	fresh.Apply(ev(event.ShotOnTarget, event.Away, 1, 20))

	_, away := s.Flow()
	_, want := fresh.Flow()
	near(t, "away", away, want)
}
