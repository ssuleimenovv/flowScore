package statsbomb

import (
	"testing"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

const homeID = 43

func TestMapEvent(t *testing.T) {
	tests := []struct {
		name  string
		raw   rawEvent
		types []event.Type
		side  event.Side
	}{
		{
			name:  "goal keeps xG",
			raw:   shot(homeID, "Goal", 0.41),
			types: []event.Type{event.Goal},
			side:  event.Home,
		},
		{
			name:  "blocked shot by away team",
			raw:   shot(1, "Blocked", 0.05),
			types: []event.Type{event.ShotBlocked},
			side:  event.Away,
		},
		{
			name: "foul with yellow becomes two events",
			raw: rawEvent{
				ID:   "f1",
				Type: ref{Name: "Foul Committed"},
				Team: ref{ID: homeID},
				FoulCommitted: &struct {
					Card *ref `json:"card"`
				}{Card: &ref{Name: "Yellow Card"}},
			},
			types: []event.Type{event.Foul, event.YellowCard},
			side:  event.Home,
		},
		{
			name:  "ordinary pass is ignored",
			raw:   rawEvent{Type: ref{Name: "Pass"}, Team: ref{ID: homeID}},
			types: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapEvent(tt.raw, "m1", homeID)

			if len(got) != len(tt.types) {
				t.Fatalf("got %d events, want %d", len(got), len(tt.types))
			}
			for i, e := range got {
				if e.Type != tt.types[i] {
					t.Errorf("event %d: type %q, want %q", i, e.Type, tt.types[i])
				}
				if e.Side != tt.side {
					t.Errorf("event %d: side %q, want %q", i, e.Side, tt.side)
				}
			}
		})
	}
}

func TestMapEventScalesPosition(t *testing.T) {
	r := shot(homeID, "Saved", 0.1)
	r.Location = []float64{120, 40}

	got := mapEvent(r, "m1", homeID)

	if got[0].Pos.X != 100 || got[0].Pos.Y != 50 {
		t.Errorf("pos = %+v, want {X:100 Y:50}", *got[0].Pos)
	}
}

func shot(teamID int, outcome string, xg float64) rawEvent {
	r := rawEvent{
		ID:   "s1",
		Type: ref{Name: "Shot"},
		Team: ref{ID: teamID},
	}
	r.Shot = &struct {
		XG      float64 `json:"statsbomb_xg"`
		Outcome ref     `json:"outcome"`
	}{XG: xg, Outcome: ref{Name: outcome}}
	return r
}
