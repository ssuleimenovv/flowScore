package flow

import (
	"context"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Update is one flow value for subscribers
type Update struct {
	MatchID string
	At      time.Duration // match time
	Home    float64
	Away    float64
	Cause   *event.Event // nil when the update comes from a tick
}

// Engine runs the flow actor of the match
type Engine struct {
	Params Params
	Speed  float64       // match second per real second: 1 live, 60 replay
	Tick   time.Duration // recompute interval, in match time
}

// Run starts the actor goroutine. Only this goroutine touches the match state
// so no mutex is needed. The returned channel is closed when events is closed
// or ctx is cancelled
func (en Engine) Run(ctx context.Context, matchID string, events <-chan event.Event) <-chan Update {
	out := make(chan Update, 16)

	go func() {
		defer close(out)

		state := NewState(en.Params)
		clock := matchClock{speed: en.Speed}

		ticker := time.NewTicker(time.Duration(float64(en.Tick) / en.Speed))
		defer ticker.Stop()

		for {
			var u Update
			select {
			case e, ok := <-events:
				if !ok {
					return
				}
				state.Apply(e)
				clock.sync(e.Elapsed, time.Now())
				u = snapshot(matchID, state, &e)

			case now := <-ticker.C:
				state.Advance(clock.at(now))
				u = snapshot(matchID, state, nil)

			case <-ctx.Done():
				return
			}

			select {
			case out <- u:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

func snapshot(matchID string, s *State, cause *event.Event) Update {
	home, away := s.Flow()
	return Update{MatchID: matchID, At: s.At(), Home: home, Away: away, Cause: cause}
}

// matchClock estimates match time between events:
// the last event's time plus the real time since it, scaled by speed
type matchClock struct {
	speed   float64
	elapsed time.Duration
	wall    time.Time
}

func (c *matchClock) sync(elapsed time.Duration, now time.Time) {
	c.elapsed = elapsed
	c.wall = now
}

func (c *matchClock) at(now time.Time) time.Duration {
	if c.wall.IsZero() {
		return 0
	}
	return c.elapsed + time.Duration(float64(now.Sub(c.wall))*c.speed)
}
