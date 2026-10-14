package live

import (
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
)

func TestLineupsFollowTheMatch(t *testing.T) {
	store := NewStore()
	p := NewPublisher(NewHub(), store, "m1", flow.DefaultParams())
	off := &event.Moment{Period: 2, Elapsed: 70*time.Minute + 26*time.Second}
	p.UseLineups(&event.Lineups{Home: event.Lineup{Formation: "4-4-2", Players: []event.Player{
		{ID: "72", Name: "Kelechi Iheanacho", Number: "72", Starter: true, Pos: &event.Position{X: 58, Y: 67}, Off: off},
		{ID: "7", Name: "Raheem Sterling", Number: "7", On: off},
	}}})

	// Before kick-off there are no lineups in the store yet; Schedule puts them in
	p.Schedule(event.Match{ID: "m1"})
	at := func(status string, period int, elapsed time.Duration) TeamLineup {
		store.update("m1", func(s *Snapshot) { s.Status, s.Period, s.At = status, period, elapsed })
		snap, _ := store.Get("m1")
		l, ok := LineupsOf(snap)
		if !ok {
			t.Fatal("no lineups")
		}
		return l.Home
	}

	before := at("scheduled", 1, 0)
	if *before.Formation != "4-4-2" || len(before.Starters) != 1 || len(before.Bench) != 1 || before.Starters[0].Off != nil {
		t.Fatalf("before kick-off %+v", before)
	}
	// At 60' the substitution of the 70th minute is still to come
	if l := at("live", 2, 60*time.Minute); l.Starters[0].Off != nil || l.Bench[0].On != nil {
		t.Errorf("at 60' %+v", l)
	}
	l := at("live", 2, 75*time.Minute)
	if got := l.Starters[0].Off; got == nil || got.Minute != 70 || got.AddedTime != nil {
		t.Errorf("at 75' Iheanacho off at %+v, want 70'", got)
	}
	if got := l.Bench[0].On; got == nil || got.Minute != 70 {
		t.Errorf("at 75' Sterling on at %+v, want 70'", got)
	}

	if _, ok := LineupsOf(Snapshot{}); ok {
		t.Error("lineups of a match without them")
	}
}
