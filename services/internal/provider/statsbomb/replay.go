package statsbomb

import (
	"context"
	"errors"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Replay plays a historical StatsBomb match as if it were live.
type Replay struct {
	MatchesPath string
	EventsDir   string
	Speed       float64 // 1 = real time, 60 = one match minute per second
}

func (r *Replay) Stream(ctx context.Context, matchID string) (<-chan event.Event, error) {
	if r.Speed <= 0 {
		return nil, errors.New("speed must be positive")
	}
	id, err := strconv.Atoi(matchID)
	if err != nil {
		return nil, err
	}

	eventsPath := filepath.Join(r.EventsDir, "events-"+matchID+".json")
	evs, err := LoadMatch(r.MatchesPath, eventsPath, id)
	if err != nil {
		return nil, err
	}
	return play(ctx, evs, r.Speed), nil
}

// play sends events with the same gaps as in the real match, scaled by speed.
func play(ctx context.Context, evs []event.Event, speed float64) <-chan event.Event {
	out := make(chan event.Event)

	go func() {
		defer close(out)

		var prev time.Duration
		for _, e := range evs {
			// At the start of the second half Elapsed goes back from 47:xx to 45:00,
			// so the gap is negative and we don't wait
			if gap := e.Elapsed - prev; gap > 0 {
				timer := time.NewTimer(time.Duration(float64(gap) / speed))
				select {
				case <-timer.C:
				case <-ctx.Done():
					timer.Stop()
					return
				}
			}
			prev = e.Elapsed

			select {
			case out <- e:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
