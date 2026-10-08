package flow

import (
	"testing"
	"time"
)

func TestClockStopsForTheBreak(t *testing.T) {
	c := matchClock{speed: 1}
	start := time.Now()
	c.sync(45*time.Minute, start, true)
	if at := c.at(start.Add(10 * time.Minute)); at != 45*time.Minute {
		t.Errorf("10 minutes into the break: %s, want 45m", at)
	}
	c.sync(45*time.Minute, start.Add(15*time.Minute), false)
	if at := c.at(start.Add(16 * time.Minute)); at != 46*time.Minute {
		t.Errorf("a minute into the second half: %s, want 46m", at)
	}
}
