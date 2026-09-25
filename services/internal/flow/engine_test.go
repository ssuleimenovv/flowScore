package flow

import (
	"context"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func fastEngine() Engine {
	return Engine{Params: DefaultParams(), Speed: 600, Tick: 5 * time.Second}
}

func TestEngineUpdatesOnEvent(t *testing.T) {
	events := make(chan event.Event)
	updates := fastEngine().Run(context.Background(), "m1", events)

	events <- ev(event.Goal, event.Home, 1, 10)
	u := next(t, updates, func(u Update) bool { return u.Cause != nil })

	if u.Cause.Type != event.Goal || u.Home < 74 {
		t.Errorf("got %+v, want a goal update with home ≥ 74", u)
	}

	close(events)
	for range updates {
	}
}

func TestEngineDecaysOnTicks(t *testing.T) {
	events := make(chan event.Event)
	updates := fastEngine().Run(context.Background(), "m1", events)
	defer close(events)

	events <- ev(event.Goal, event.Home, 1, 10)
	afterGoal := next(t, updates, func(u Update) bool { return u.Cause != nil })

	later := next(t, updates, func(u Update) bool { return u.Cause == nil && u.Home < afterGoal.Home })
	if later.At <= afterGoal.At {
		t.Errorf("tick at %v, want after %v", later.At, afterGoal.At)
	}
}

func TestEngineStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	updates := fastEngine().Run(ctx, "m1", make(chan event.Event))
	cancel()

	deadline := time.After(time.Second)
	for {
		select {
		case _, open := <-updates:
			if !open {
				return
			}
		case <-deadline:
			t.Fatal("updates not closed after cancel")
		}
	}
}

// next reads updates until match returns true, failing after a second.
func next(t *testing.T, updates <-chan Update, match func(Update) bool) Update {
	t.Helper()
	deadline := time.After(time.Second)
	for {
		select {
		case u, open := <-updates:
			if !open {
				t.Fatal("updates closed")
			}
			if match(u) {
				return u
			}
		case <-deadline:
			t.Fatal("no matching update within a second")
		}
	}
}
