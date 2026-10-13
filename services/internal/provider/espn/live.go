package espn

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"slices"
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
		t := tracker{matchID: matchID, sent: map[string]bool{}}
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
	s, plays, err := l.read(ctx, t.matchID)
	if err != nil {
		log.Printf("espn %s: %v", t.matchID, err)
		return false
	}
	events, over := t.next(s, plays)
	for _, e := range events {
		select {
		case out <- e:
		case <-ctx.Done():
			return true
		}
	}
	return over
}

// read reads the summary once: the match as a whole and its plays
func (l *Live) read(ctx context.Context, matchID string) (Summary, []event.Event, error) {
	body, err := l.Client.Summary(ctx, l.League, matchID)
	if err != nil {
		return Summary{}, nil, err
	}
	s, err := ParseSummary(body)
	if err != nil {
		return Summary{}, nil, err
	}
	plays, err := ParsePlays(body, s.Match)
	return s, plays, err
}

// The second half starts at 45:00 on ESPN's clock, as on ours
const secondHalf = 45 * time.Minute

// tracker turns successive summaries of one match into events: the plays
// not sent yet, the minutes of possession, then where the clock is. It knows
// nothing of time itself, so it is tested without a network or a clock.
type tracker struct {
	matchID string
	sent    map[string]bool // the plays sent so far, by ID
	period  int             // the period of the last clock sent
	ball    ball
}

// next returns the events of a new summary and whether the match is over.
func (t *tracker) next(s Summary, plays []event.Event) ([]event.Event, bool) {
	if s.State == "pre" {
		return nil, false
	}

	var out []event.Event
	for _, p := range plays {
		// A play of a half that is over came too late: the publisher has
		// moved on to the next one, and would go back
		if t.sent[p.ID] || p.Period < t.period {
			continue
		}
		t.sent[p.ID] = true
		out = append(out, p)
	}
	// A finished match comes without a clock or a period: a clock at 0:00 of
	// period 0 would take the match back to before kick-off
	if s.State == "post" {
		return out, true
	}

	period, at, stopped := s.Period, s.Clock, false
	if s.Status == "STATUS_HALFTIME" {
		// The break belongs to the second half, with its clock stopped at 45:00
		period, at, stopped = 2, secondHalf, true
	}

	for _, e := range t.ball.next(t.matchID, period, at, s.Possession) {
		if e.Period >= t.period { // late, as a play can be
			out = append(out, e)
		}
	}
	// The minutes of possession go between the plays: joining a match under
	// way, its first hour must not reach the engine after all its shots
	slices.SortStableFunc(out, func(a, b event.Event) int {
		return cmp.Or(cmp.Compare(a.Period, b.Period), cmp.Compare(a.Elapsed, b.Elapsed))
	})

	t.period = period

	out = append(out, event.Event{
		ID:      fmt.Sprintf("%s:clock", t.matchID),
		MatchID: t.matchID,
		Type:    event.Clock,
		Period:  period,
		Elapsed: at,
		Stopped: stopped,
	})
	return out, false
}
