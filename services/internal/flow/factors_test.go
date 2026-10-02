package flow

import (
	"math"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func factor(t *testing.T, s *State, side event.Side, g Group) Factor {
	t.Helper()
	for _, f := range s.Factors() {
		if f.Side == side && f.Group == g {
			return f
		}
	}
	t.Fatalf("no %s factor for %s in %+v", g, side, s.Factors())
	return Factor{}
}

func TestFactorsAddUpToImpulse(t *testing.T) {
	s := NewState(DefaultParams())
	share := 0.3
	for _, e := range []event.Event{
		ev(event.ShotOnTarget, event.Home, 1, 5),
		ev(event.Corner, event.Home, 1, 6),
		{Type: event.Possession, Period: 1, Elapsed: 7 * time.Minute, HomeShare: &share},
		ev(event.KeyPass, event.Away, 1, 8),
		ev(event.RedCard, event.Home, 1, 9), // hits the floor: the parts must still add up
		ev(event.Goal, event.Away, 1, 30),
		ev(event.Substitution, event.Home, 2, 50),
	} {
		s.Apply(e)
	}
	s.Advance(55 * time.Minute)

	for _, side := range []event.Side{event.Home, event.Away} {
		sum := 0.0
		for _, v := range s.parts[side] {
			sum += v
		}
		if math.Abs(sum-s.impulse[side]) > 1e-9 {
			t.Errorf("%s: parts add up to %.6f, impulse is %.6f", side, sum, s.impulse[side])
		}
	}
}

func TestFactorCountsRecentEvents(t *testing.T) {
	s := NewState(DefaultParams())
	for _, minute := range []int{2, 10, 12, 15} {
		s.Apply(ev(event.ShotOffTarget, event.Home, 1, minute))
	}
	s.Advance(16 * time.Minute)

	// 2′ is out of the 10-minute window, 10′ to 15′ are in
	f := factor(t, s, event.Home, Shots)
	if f.Count != 3 || f.Since != 6*time.Minute {
		t.Errorf("shots = %+v, want 3 in the last 6 minutes", f)
	}
	if f.Value <= 0 {
		t.Errorf("shots add %.2f, want a positive value", f.Value)
	}
}

func TestRedCardCancelsWhatTheTeamBuilt(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.ShotOnTarget, event.Home, 1, 10))
	s.Apply(ev(event.RedCard, event.Home, 1, 10))

	// The card takes away exactly what the shot had built, and no more
	shots, red := factor(t, s, event.Home, Shots), factor(t, s, event.Home, RedCards)
	near(t, "shots + red card", shots.Value+red.Value, 0)
}

func TestPossessionFactorHasTheShare(t *testing.T) {
	s := NewState(DefaultParams())
	for m, share := range []float64{0.7, 0.6, 0.8} {
		share := share
		s.Apply(event.Event{Type: event.Possession, Period: 1, Elapsed: time.Duration(m+1) * time.Minute, HomeShare: &share})
	}

	home := factor(t, s, event.Home, Possession)
	if home.Share == nil || math.Abs(*home.Share-0.7) > 1e-9 || home.Count != 3 {
		t.Errorf("home possession = %+v, want 70%% over 3 minutes", home)
	}
}

func TestBreakStartsTheWindowAgain(t *testing.T) {
	s := NewState(DefaultParams())
	s.Apply(ev(event.Corner, event.Home, 1, 44))
	s.Apply(ev(event.Foul, event.Away, 2, 46))

	if f := factor(t, s, event.Home, Corners); f.Count != 0 {
		t.Errorf("corners = %+v, want the first-half corner left out of the count", f)
	}
}
