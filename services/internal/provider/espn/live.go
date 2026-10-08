package espn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Live streams a match as ESPN reports it: it reads the summary every Every
// and sends what changed as events. The channel is closed at the final
// whistle or when ctx is done.
type Live struct {
	Client *Client
	League string // "eng.1"
	Every  time.Duration
}

func (l *Live) Stream(ctx context.Context, matchID string) (<-chan event.Event, error) {
	out := make(chan event.Event)
	go func() {
		defer close(out)
		t := tracker{matchID: matchID}
		for {
			if over := l.poll(ctx, &t, out); over {
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(l.Every):
			}
		}
	}()
	return out, nil
}

// poll reads the summary once, sends its events and reports whether the
// match is over. A failed read is only logged: the next poll tries again.
func (l *Live) poll(ctx context.Context, t *tracker, out chan<- event.Event) bool {
	body, err := l.Client.Summary(ctx, l.League, t.matchID)
	if err == nil {
		var s Summary
		if s, err = ParseSummary(body); err == nil {
			events, over := t.next(s)
			for _, e := range events {
				select {
				case out <- e:
				case <-ctx.Done():
					return true
				}
			}
			return over
		}
	}
	log.Printf("espn %s: %v", t.matchID, err)
	return false
}

// The second half starts at 45:00 on ESPN's clock, as on ours
const secondHalf = 45 * time.Minute

// tracker turns successive summaries of one match into events: where the
// clock is, and a goal for every goal the score went up by. It knows nothing
// of time itself, so it is tested without a network or a clock.
type tracker struct {
	matchID string
	score   Score // what has been sent so far
}

// next returns the events of a new summary and whether the match is over.
func (t *tracker) next(s Summary) ([]event.Event, bool) {
	var out []event.Event
	if s.State == "pre" {
		return nil, false
	}

	period, at, stopped := s.Period, s.Clock, false
	if s.Status == "STATUS_HALFTIME" {
		// The break belongs to the second half, with its clock stopped at 45:00
		period, at, stopped = 2, secondHalf, true
	}

	out = append(out, t.goals(event.Home, &t.score.Home, s.Score.Home, period, at)...)
	out = append(out, t.goals(event.Away, &t.score.Away, s.Score.Away, period, at)...)
	out = append(out, event.Event{
		ID:      fmt.Sprintf("%s:clock", t.matchID),
		MatchID: t.matchID,
		Type:    event.Clock,
		Period:  period,
		Elapsed: at,
		Stopped: stopped,
	})
	return out, s.State == "post"
}

// goals sends a goal for each one the score went up by. A goal taken back
// (VAR) only lowers the count: the event is out already.
func (t *tracker) goals(side event.Side, sent *int, now, period int, at time.Duration) []event.Event {
	var out []event.Event
	for *sent < now {
		*sent++
		out = append(out, event.Event{
			ID:      fmt.Sprintf("%s:goal:%s:%d", t.matchID, side, *sent),
			MatchID: t.matchID,
			Type:    event.Goal,
			Side:    side,
			Period:  period,
			Elapsed: at,
		})
	}
	*sent = min(*sent, now)
	return out
}
