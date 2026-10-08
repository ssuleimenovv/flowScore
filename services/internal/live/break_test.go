package live

import (
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
)

// A live source reports the break with the clock stopped: the match is at
// halftime for as long as it lasts, though its updates are of the second half.
func TestLiveBreakKeepsHalftime(t *testing.T) {
	store := NewStore()
	p := NewPublisher(NewHub(), store, "m1", flow.DefaultParams())
	p.Start(event.Match{ID: "m1"})

	clock := func(period int, at time.Duration, stopped bool) flow.Update {
		return flow.Update{At: at, Cause: &event.Event{Type: event.Clock, Period: period, Elapsed: at, Stopped: stopped}}
	}
	updates := make(chan flow.Update)
	done := make(chan struct{})
	go func() {
		p.Run(updates)
		close(done)
	}()
	status := func() string {
		s, _ := store.Get("m1")
		return s.Status
	}

	updates <- clock(1, 46*time.Minute, false)
	updates <- clock(2, 45*time.Minute, true)
	updates <- flow.Update{At: 45 * time.Minute} // a tick during the break
	updates <- flow.Update{At: 45 * time.Minute}
	if got := status(); got != "halftime" {
		t.Errorf("during the break: %s, want halftime", got)
	}
	updates <- clock(2, 46*time.Minute, false)
	updates <- flow.Update{At: 46 * time.Minute}
	if got := status(); got != "live" {
		t.Errorf("in the second half: %s, want live", got)
	}
	close(updates)
	<-done
}
