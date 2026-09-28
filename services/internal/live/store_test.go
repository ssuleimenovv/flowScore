package live

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
)

func TestStoreGetReturnsCopy(t *testing.T) {
	store := NewStore()
	store.Start(event.Match{ID: "m1"}, 0)
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

	p.Start(event.Match{ID: "m1"})
	p.Run(updates)

	snap, ok := store.Get("m1")
	if !ok {
		t.Fatal("no snapshot")
	}
	if snap.Score != (Score{Home: 1}) {
		t.Errorf("score = %+v, want 1:0", snap.Score)
	}
	// The goal at 7:30 is 7′, as the live badge shows it; then the final whistle
	if len(snap.Events) != 2 || snap.Events[0].Minute != 7 || snap.Events[1].Type != event.Fulltime {
		t.Errorf("events = %+v, want a goal at 7' and the final whistle", snap.Events)
	}
	// goal: match.event + match.stats + flow.update, tick: flow.update, fulltime: match.event
	if snap.Seq != 5 {
		t.Errorf("seq = %d, want 5", snap.Seq)
	}
	if len(snap.Stats) == 0 || snap.Stats[1] != (StatRow{Key: "shots", Home: 1}) {
		t.Errorf("stats = %+v, want one home shot", snap.Stats)
	}
	if snap.Status != "finished" {
		t.Errorf("status = %q, want finished", snap.Status)
	}
}

func TestStartKeepsSeq(t *testing.T) {
	store := NewStore()
	p := NewPublisher(NewHub(), store, "m1", flow.DefaultParams())
	p.Start(event.Match{ID: "m1"})
	p.Run(closed(flow.Update{At: time.Minute, Home: 50}))

	// The replay starts again: the score is gone, seq is not
	p.Start(event.Match{ID: "m1"})

	snap, _ := store.Get("m1")
	if snap.Status != "live" || len(snap.Events) != 0 || snap.Seq != 2 {
		t.Errorf("snapshot = %+v, want live, no events, seq 2", snap)
	}
}

func TestHalftimeWhistle(t *testing.T) {
	hub := NewHub()
	client := hub.Subscribe("m1")
	store := NewStore()
	p := NewPublisher(hub, store, "m1", flow.DefaultParams())

	goalAt := 30 * time.Minute
	p.Start(event.Match{ID: "m1"})
	p.Run(closed(
		flow.Update{At: goalAt, Cause: &event.Event{Type: event.Goal, Side: event.Away, Period: 1, Elapsed: goalAt}},
		flow.Update{At: 46*time.Minute + 20*time.Second}, // a tick in added time
		flow.Update{At: 45 * time.Minute, Cause: &event.Event{Type: event.Foul, Side: event.Home, Period: 2, Elapsed: 45 * time.Minute}},
	))

	ht := read[MatchEvent](t, client, "match.event", func(e MatchEvent) bool { return e.Type == event.Halftime })
	if ht.Minute != 45 || ht.AddedTime == nil || *ht.AddedTime != 2 || ht.Side != nil {
		t.Errorf("halftime = %+v, want 45+2 without a side", ht)
	}

	snap, _ := store.Get("m1")
	if snap.Halftime == nil || *snap.Halftime != (Score{Away: 1}) {
		t.Errorf("halftime score = %+v, want 0:1", snap.Halftime)
	}
}

func TestFlowUpdateCarriesClock(t *testing.T) {
	hub := NewHub()
	client := hub.Subscribe("m1")
	p := NewPublisher(hub, NewStore(), "m1", flow.DefaultParams())

	kickoff := 45 * time.Minute
	p.Start(event.Match{ID: "m1"})
	p.Run(closed(flow.Update{
		At:    kickoff + 12*time.Second,
		Cause: &event.Event{Type: event.Foul, Period: 2, Elapsed: kickoff},
	}))

	got := read[FlowUpdate](t, client, "flow.update", nil).Clock
	want := Clock{ElapsedSeconds: 2712, Period: "second_half"}
	if got.ElapsedSeconds != want.ElapsedSeconds || got.Period != want.Period {
		t.Errorf("clock = %+v, want %+v", got, want)
	}
	if got.ObservedAt.IsZero() {
		t.Error("observedAt is empty")
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

// closed returns a channel that yields the updates and is then closed,
// the way the Flow Engine ends a match.
func closed(updates ...flow.Update) <-chan flow.Update {
	ch := make(chan flow.Update, len(updates))
	for _, u := range updates {
		ch <- u
	}
	close(ch)
	return ch
}

// read returns the data of the first message of the given type that match
// accepts (nil accepts any). Run has already finished, so every message
// is waiting in the client's buffer.
func read[T any](t *testing.T, c *Client, kind string, match func(T) bool) T {
	t.Helper()
	for {
		select {
		case raw := <-c.Messages():
			var msg struct {
				Type string
				Data T
			}
			if err := json.Unmarshal(raw, &msg); err != nil {
				t.Fatal(err)
			}
			if msg.Type == kind && (match == nil || match(msg.Data)) {
				return msg.Data
			}
		default:
			t.Fatalf("no %s message", kind)
		}
	}
}
