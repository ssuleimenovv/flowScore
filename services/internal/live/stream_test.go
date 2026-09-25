package live

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

// TestStreamDeliversToSubscriber goes through the real route, HTTP upgrade
// and WebSocket frames, the way a browser does.
func TestStreamDeliversToSubscriber(t *testing.T) {
	hub := NewHub()
	mux := http.NewServeMux()
	Register(mux, hub, nil)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/matches/m1/stream"
	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.CloseNow()

	// The handler subscribes right after the handshake; wait for it,
	// otherwise the message could be published before anyone listens.
	select {
	case <-hub.Joined():
	case <-ctx.Done():
		t.Fatal("handler never subscribed")
	}
	if n := hub.Count("m1"); n != 1 {
		t.Fatalf("subscribers of m1 = %d, want 1 (is the path parameter read correctly?)", n)
	}

	hub.Publish("m1", []byte(`{"type":"flow.update"}`))

	_, msg, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if got := string(msg); got != `{"type":"flow.update"}` {
		t.Errorf("got %s", got)
	}
}
