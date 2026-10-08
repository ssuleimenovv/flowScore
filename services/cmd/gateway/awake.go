package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

// keepAwake asks the service's own public address for /healthz every so
// often. Render's free plan puts a service to sleep after 15 minutes without
// requests from outside, and a sleeping gateway misses the real matches it
// follows. The request goes out and back through Render's proxy, so it
// counts as one from outside. A cron job on GitHub was meant to do this, but
// GitHub runs scheduled jobs hours late.
func keepAwake(ctx context.Context, url string, every time.Duration) {
	client := &http.Client{Timeout: 30 * time.Second}
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(every):
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url+"/healthz", nil)
		if err != nil {
			log.Printf("keep awake: %v", err)
			return
		}
		res, err := client.Do(req)
		if err != nil {
			log.Printf("keep awake: %v", err)
			continue
		}
		res.Body.Close()
	}
}
