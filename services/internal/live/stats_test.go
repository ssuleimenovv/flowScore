package live

import (
	"testing"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func TestStatsCountsEvents(t *testing.T) {
	s := NewStats()
	xg := func(v float64) *float64 { return &v }
	share := func(v float64) *float64 { return &v }

	s.Add(event.Event{Type: event.Goal, Side: event.Home, XG: xg(0.41)})
	s.Add(event.Event{Type: event.ShotOffTarget, Side: event.Home, XG: xg(0.12)})
	s.Add(event.Event{Type: event.ShotBlocked, Side: event.Away, XG: xg(0.05)})
	s.Add(event.Event{Type: event.KeyPass, Side: event.Home})
	s.Add(event.Event{Type: event.Tackle, Side: event.Away})
	s.Add(event.Event{Type: event.Possession, HomeShare: share(0.7)})
	s.Add(event.Event{Type: event.Possession, HomeShare: share(0.6)})

	want := []StatRow{
		{Key: "possession", Home: 65, Away: 35},
		{Key: "shots", Home: 2, Away: 1},
		{Key: "shots_on_target", Home: 1, Away: 0},
		{Key: "xg", Home: 0.53, Away: 0.05},
		{Key: "key_passes", Home: 1, Away: 0},
		{Key: "pressure_tackles", Home: 0, Away: 1},
	}
	got := s.Rows()
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("row %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestStatsIgnoresTimelineOnlyEvents(t *testing.T) {
	s := NewStats()
	if s.Add(event.Event{Type: event.Foul, Side: event.Home}) {
		t.Error("a foul changed the stats")
	}
	if s.Add(event.Event{Type: event.Goal, Side: event.Home, OwnGoal: true}) {
		t.Error("an own goal counted as a shot")
	}
	if row := s.Rows()[0]; row.Home != 50 || row.Away != 50 {
		t.Errorf("possession before any minute = %+v, want 50:50", row)
	}
}
