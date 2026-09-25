package statsbomb

import (
	"math"
	"testing"
)

func at(second int, teamID int) rawEvent {
	return rawEvent{Period: 1, Minute: second / 60, Second: second % 60, PossessionTeam: ref{ID: teamID}}
}

func TestPossessionShare(t *testing.T) {
	// Home holds the ball 0–30 s and 45–60 s, away 30–45 s: 45 of 60 seconds.
	raw := []rawEvent{at(0, homeID), at(30, 1), at(45, homeID), at(60, homeID)}

	got := possessionEvents(raw, "m1", homeID)

	if len(got) != 1 {
		t.Fatalf("got %d events, want 1", len(got))
	}
	if share := *got[0].HomeShare; math.Abs(share-0.75) > 1e-9 {
		t.Errorf("home share = %.3f, want 0.75", share)
	}
}

func TestPossessionSkipsPeriodChange(t *testing.T) {
	first := at(47*60+10, homeID)
	second := rawEvent{Period: 2, Minute: 45, PossessionTeam: ref{ID: 1}}

	if got := possessionEvents([]rawEvent{first, second}, "m1", homeID); len(got) != 0 {
		t.Errorf("got %d events across the break, want 0", len(got))
	}
}
