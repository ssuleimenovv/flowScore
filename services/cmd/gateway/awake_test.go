package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestKeepAwakePingsHealthz(t *testing.T) {
	var pings atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/healthz" {
			pings.Add(1)
		}
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		keepAwake(ctx, srv.URL, 10*time.Millisecond)
		close(done)
	}()
	deadline := time.Now().Add(2 * time.Second)
	for pings.Load() < 2 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	cancel()
	<-done // and it stops with the server
	if pings.Load() < 2 {
		t.Errorf("%d pings, want at least 2", pings.Load())
	}
}
