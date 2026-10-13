package espn

import (
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func TestPlayed(t *testing.T) {
	body := boxscore(summary("post", "STATUS_FULL_TIME", 0, 0, "1", "0",
		play("7", "goal", 1, 600, "Arsenal"), play("8", "shot-blocked", 2, 3000, "Leeds United")), "60", "40")
	s, events, err := Played([]byte(body))
	if err != nil {
		t.Fatal(err)
	}
	if s.State != "post" {
		t.Errorf("summary %+v", s)
	}
	// Two plays and a minute of possession for each of the 90
	if len(events) != 92 {
		t.Fatalf("%d events", len(events))
	}
	// In match order; at the same moment the play goes first, as the live tracker sends it
	if g, p := events[9], events[10]; g.Type != event.Goal || g.Elapsed != 10*time.Minute ||
		p.Type != event.Possession || p.Elapsed != 10*time.Minute {
		t.Errorf("around the goal: %+v, %+v", g, p)
	}
	if last := events[len(events)-1]; last.Type != event.Possession || last.Elapsed != 90*time.Minute || *last.HomeShare != 0.6 {
		t.Errorf("last %+v", last)
	}
}
