package season

import (
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func TestPlayers(t *testing.T) {
	at := func(minute int) *event.Moment {
		return &event.Moment{Period: 2, Elapsed: time.Duration(minute) * time.Minute}
	}
	l := &event.Lineups{
		Home: event.Lineup{Players: []event.Player{
			{ID: "9", Name: "Haaland", Starter: true, Off: at(70)},
			{ID: "47", Name: "Foden", Starter: true},
			{ID: "19", Name: "Marmoush", On: at(70)},
			{ID: "31", Name: "Ederson"}, // stayed on the bench
		}},
		Away: event.Lineup{Players: []event.Player{
			{ID: "4", Name: "van Dijk", Starter: true},
		}},
	}
	// As in ESPN's commentary: the name, no ID
	goal := event.Event{ID: "g1", Type: event.Goal, Side: event.Home, Player: "Haaland", Elapsed: 20 * time.Minute}
	events := []event.Event{
		{ID: "g1:key", Type: event.KeyPass, Side: event.Home, PlayerID: "47", Elapsed: 20 * time.Minute},
		goal,
		{ID: "s2", Type: event.ShotOffTarget, Side: event.Home, PlayerID: "9", Elapsed: 30 * time.Minute},
		// A key pass to a shot that missed is no assist
		{ID: "s3:key", Type: event.KeyPass, Side: event.Home, PlayerID: "47", Elapsed: 50 * time.Minute},
		{ID: "s3", Type: event.ShotBlocked, Side: event.Home, PlayerID: "19", Elapsed: 80 * time.Minute},
		// van Dijk's own goal: not a goal of his, not a shot
		{ID: "og", Type: event.Goal, OwnGoal: true, Side: event.Home, PlayerID: "4", Elapsed: 85 * time.Minute},
		{ID: "end", Type: event.Fulltime, Elapsed: 94 * time.Minute},
	}

	players := Players(events, l, nil)
	if len(players) != 4 {
		t.Fatalf("%d players, want 4: the bench is left out", len(players))
	}
	haaland, foden, marmoush, vanDijk := players[0], players[1], players[2], players[3]
	if haaland.Minutes != 70 || haaland.Goals != 1 || haaland.Shots != 2 || haaland.OnTarget != 1 {
		t.Errorf("Haaland %+v", haaland)
	}
	if foden.Minutes != 94 || foden.Assists != 1 {
		t.Errorf("Foden %+v", foden)
	}
	if marmoush.Starter || marmoush.Minutes != 24 || marmoush.Shots != 1 || marmoush.OnTarget != 0 {
		t.Errorf("Marmoush %+v", marmoush)
	}
	if vanDijk.Side != event.Away || vanDijk.Goals != 0 || vanDijk.Shots != 0 {
		t.Errorf("van Dijk %+v", vanDijk)
	}

	if Players(events, nil, nil) != nil {
		t.Error("players without lineups")
	}
}
