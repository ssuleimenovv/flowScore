package api

import (
	"context"
	"sync"
	"time"
)

// cache keeps what ESPN said about a team or a player for keepFor, by its
// ID. When ESPN does not answer, what it said last stands in, however old.
type cache[T any] struct {
	keepFor time.Duration
	fetch   func(ctx context.Context, league, id string) (T, error)
	mu      sync.Mutex
	kept    map[string]keptValue[T]
}

type keptValue[T any] struct {
	value T
	at    time.Time
}

func newCache[T any](keepFor time.Duration, fetch func(ctx context.Context, league, id string) (T, error)) *cache[T] {
	return &cache[T]{keepFor: keepFor, fetch: fetch, kept: map[string]keptValue[T]{}}
}

func (c *cache[T]) get(ctx context.Context, league, id string) (T, error) {
	c.mu.Lock()
	k, ok := c.kept[id]
	c.mu.Unlock()
	if ok && time.Since(k.at) < c.keepFor {
		return k.value, nil
	}
	v, err := c.fetch(ctx, league, id)
	if err != nil {
		if ok {
			return k.value, nil
		}
		var zero T
		return zero, err
	}
	c.mu.Lock()
	c.kept[id] = keptValue[T]{value: v, at: time.Now()}
	c.mu.Unlock()
	return v, nil
}
