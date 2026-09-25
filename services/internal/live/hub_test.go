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

func TestFirstSubscriber(t *testing.T) {
	hub := NewHub()
	select {
	case <-hub.FirstSubscriber():
		t.Fatal("closed before anyone subscribed")
	default:
	}

	hub.Subscribe("m1")
	hub.Subscribe("m1") // a second call must not close the channel twice

	select {
	case <-hub.FirstSubscriber():
	default:
		t.Fatal("not closed after the first subscriber")
	}
}
