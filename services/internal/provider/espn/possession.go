package espn

import (
	"fmt"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// A share of the ball that has not moved may be one ESPN has not counted
// again yet; it is taken as it is only after this long
const stale = 5 * time.Minute

// ball turns ESPN's share of the ball, which is of the whole match so far,
// into one Possession event per match minute, as the Flow engine takes it
// (docs/FLOW.md, section 2). The minutes between two readings get the share
// that moved the whole match from the first to the second
type ball struct {
	at     time.Duration // match time of the last reading taken
	share  float64       // the home share of the match up to then
	minute int           // the last minute sent
}

func (b *ball) next(matchID string, period int, at time.Duration, share *float64) []event.Event {
	// A reading is taken once a minute has ended since the last one, so
	// that a 30-second step of a share rounded to 0.1% is not all it is
	if share == nil || at < time.Duration(b.minute+1)*time.Minute ||
		*share == b.share && at-b.at < stale {
		return nil
	}
	home := (*share*at.Minutes() - b.share*b.at.Minutes()) / (at - b.at).Minutes()
	home = min(max(home, 0), 1)
	b.at, b.share = at, *share

	var out []event.Event
	for ; time.Duration(b.minute+1)*time.Minute <= at; b.minute++ {
		m := b.minute + 1
		p := period
		if m <= 45 {
			p = 1 //  the reading came in the second half, the minute was in the first
		}
		out = append(out, event.Event{
			ID:        fmt.Sprintf("%s:possession:%d", matchID, m),
			MatchID:   matchID,
			Type:      event.Possession,
			Period:    p,
			Elapsed:   time.Duration(m) * time.Minute,
			HomeShare: &home,
		})
	}
	return out
}
