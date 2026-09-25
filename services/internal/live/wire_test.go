package live

import (
	"testing"
	"time"
)

func TestMinuteOf(t *testing.T) {
	tests := []struct {
		period  int
		elapsed time.Duration
		minute  int
		added   int // 0 means no added time
	}{
		{1, 7*time.Minute + 30*time.Second, 8, 0},
		{1, 44*time.Minute + 59*time.Second, 45, 0},
		{1, 45*time.Minute + 49*time.Second, 45, 1},
		{2, 45*time.Minute + 44*time.Second, 46, 0},
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
