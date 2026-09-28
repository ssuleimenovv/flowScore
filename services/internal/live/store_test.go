package live

import (
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
)

func TestStoreGetReturnsCopy(t *testing.T) {
	store := NewStore()
	store.Start(event.Match{ID: "m1"})
	store.update("m1", func(s *Snapshot) {
		s.Events = append(s.Events, MatchEvent{ID: "e1"})
	})

	snap, _ := store.Get("m1")
	snap.Events[0].ID = "changed"

	again, _ := store.Get("m1")
	if again.Events[0].ID != "e1" {
		t.Errorf("store changed through a copy: event ID = %q", again.Events[0].ID)
	}
}

func TestPublisherKeepsSnapshot(t *testing.T) {
	hub := NewHub()
	store := NewStore()
	p := NewPublisher(hub, store, "m1", flow.DefaultParams())

	goalAt := 7*time.Minute + 30*time.Second
	updates := make(chan flow.Update, 2)
	updates <- flow.Update{
		At:    goalAt,
		Home:  74,
		Cause: &event.Event{ID: "g1", Type: event.Goal, Side: event.Home, Period: 1, Elapsed: goalAt},
	}
	updates <- flow.Update{At: 9 * time.Minute, Home: 70} // a tick
	close(updates)

	p.Run(event.Match{ID: "m1"}, updates)

	snap, ok := store.Get("m1")
	if !ok {
		t.Fatal("no snapshot")
	}
	if snap.Score != (Score{Home: 1}) {
		t.Errorf("score = %+v, want 1:0", snap.Score)
	}
	if len(snap.Events) != 1 || snap.Events[0].Minute != 8 {
		t.Errorf("events = %+v, want one goal at 8'", snap.Events)
	}
	// goal: match.event + flow.update, tick: flow.update
	if snap.Seq != 3 {
		t.Errorf("seq = %d, want 3", snap.Seq)
	}
	if snap.Status != "finished" {
		t.Errorf("status = %q, want finished", snap.Status)
	}
}

func TestUpsertPointReplacesSameMinute(t *testing.T) {
	points := []FlowPoint{{Minute: 45, Home: 10}, {Minute: 46, Home: 20}}

	points = upsertPoint(points, FlowPoint{Minute: 46, Home: 30}) // second half restarts at 46
	points = upsertPoint(points, FlowPoint{Minute: 47, Home: 40})

	if len(points) != 3 || points[1].Home != 30 || points[2].Minute != 47 {
		t.Errorf("points = %+v", points)
	}
}
