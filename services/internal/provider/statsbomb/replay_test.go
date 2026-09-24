package statsbomb

import (
	"context"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func events(minutes ...int) []event.Event {
	out := make([]event.Event, len(minutes))
	for i, m := range minutes {
		out[i] = event.Event{ID: string(rune('a' + i)), Elapsed: time.Duration(m) * time.Minute}
	}
	return out
}

func TestPlayDeliveryAllInOrder(t *testing.T) {
	in := events(1, 2, 5)

	var got []string
	for e := range play(context.Background(), in, 1e6) {
		got = append(got, e.ID)
	}

	if len(got) != 3 || got[0] != "a" || got[2] != "c" {
		t.Fatalf("got %v, want [a b c]", got)
	}
}

func TestPlayStopsOnCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := play(ctx, events(0, 60, 90), 1) // real time: the second event is an hour away

	<-ch
	cancel()

	select {
	case _, open := <-ch:
		if open {
			t.Fatal("got an event after cancel")
		}
	case <-time.After(time.Second):
		t.Fatal("channel not closed after cancel")
	}
}
