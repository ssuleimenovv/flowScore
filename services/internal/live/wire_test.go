package live

import (
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
)

func TestMinuteOf(t *testing.T) {
	tests := []struct {
		period  int
		elapsed time.Duration
		minute  int
		added   int // 0 means no added time
	}{
		{1, 7*time.Minute + 30*time.Second, 7, 0},
		{1, 44*time.Minute + 59*time.Second, 44, 0},
		{1, 45 * time.Minute, 45, 1},
		{1, 45*time.Minute + 49*time.Second, 45, 1},
		{2, 45*time.Minute + 44*time.Second, 45, 0},
		{2, 92*time.Minute + 34*time.Second, 90, 3},
	}

	for _, tt := range tests {
		minute, added := minuteOf(tt.period, tt.elapsed)
		gotAdded := 0
		if added != nil {
			gotAdded = *added
		}
		if minute != tt.minute || gotAdded != tt.added {
			t.Errorf("minuteOf(%d, %v) = %d+%d, want %d+%d",
				tt.period, tt.elapsed, minute, gotAdded, tt.minute, tt.added)
		}
	}
}

func TestToFactors(t *testing.T) {
	share := 0.68
	got := toFactors([]flow.Factor{
		{Side: event.Home, Group: flow.Shots, Value: 13.96, Count: 4, Since: 6*time.Minute + 10*time.Second},
		{Side: event.Home, Group: flow.Possession, Value: 9.04, Count: 10, Since: 9 * time.Minute, Share: &share},
		{Side: event.Away, Group: flow.Corners, Value: 1.2}, // the corner left the window but still counts
	})

	want := []FlowFactor{
		{Side: event.Home, Key: "shots", Value: 14, Count: 4, Minutes: 7},
		{Side: event.Home, Key: "possession", Value: 9, Count: 10, Minutes: 9, Share: &share},
		{Side: event.Away, Key: "corners", Value: 1.2},
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("factor %d = %+v, want %+v", i, got[i], want[i])
		}
	}

	if empty := toFactors(nil); empty == nil || len(empty) != 0 {
		t.Errorf("no factors = %#v, want an empty slice for JSON []", empty)
	}
}
