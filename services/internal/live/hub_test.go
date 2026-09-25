package live

import "testing"

func TestHubFanOut(t *testing.T) {
	hub := NewHub()
	a := hub.Subscribe("m1")
	b := hub.Subscribe("m1")
	other := hub.Subscribe("m2")

	hub.Publish("m1", []byte("goal"))

	for name, c := range map[string]*Client{"a": a, "b": b} {
		if got := string(<-c.Messages()); got != "goal" {
			t.Errorf("%s got %q, want goal", name, got)
		}
	}
	select {
	case msg := <-other.Messages():
		t.Errorf("subscriber of m2 got %q", msg)
	default:
	}
}

func TestHubDropsSlowClient(t *testing.T) {
	hub := NewHub()
	slow := hub.Subscribe("m1")

	for range sendBuffer + 1 {
		hub.Publish("m1", []byte("tick"))
	}

	for range sendBuffer {
		<-slow.Messages()
	}
	if _, open := <-slow.Messages(); open {
		t.Fatal("slow client still subscribed, want its channel closed")
	}
}

func TestJoinedCollapsesSignals(t *testing.T) {
	hub := NewHub()
	hub.Subscribe("m1")
	hub.Subscribe("m1") // must not block although nobody reads Joined yet

	<-hub.Joined()
	select {
	case <-hub.Joined():
		t.Fatal("two subscriptions produced two signals, want one")
	default:
	}
	if n := hub.Count("m1"); n != 2 {
		t.Errorf("Count = %d, want 2", n)
	}
}
