package statsbomb

import (
	"fmt"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// possessionEvents turns StatsBomb's possession_team into one Possession event
// per match minute with the home team's share of the ball (docs/FLOW.md, section 2).
// The time between two consecutive events belongs to the team in possession at the first one.
func possessionEvents(raw []rawEvent, matchID string, homeTeamID int) []event.Event {
	type minuteKey struct{ period, minute int }
	type share struct{ home, total, end time.Duration }

	shares := map[minuteKey]*share{}
	var order []minuteKey

	for i := 0; i+1 < len(raw); i++ {
		cur, next := raw[i], raw[i+1]
		if cur.Period != next.Period {
			continue
		}
		dt := elapsed(next) - elapsed(cur)
		if dt <= 0 {
			continue
		}

		key := minuteKey{cur.Period, cur.Minute}
		s, ok := shares[key]
		if !ok {
			s = &share{}
			shares[key] = s
			order = append(order, key)
		}
		s.total += dt
		if cur.PossessionTeam.ID == homeTeamID {
			s.home += dt
		}
		s.end = elapsed(next)
	}

	out := make([]event.Event, 0, len(order))
	for _, key := range order {
		s := shares[key]
		homeShare := s.home.Seconds() / s.total.Seconds()
		out = append(out, event.Event{
			ID:        fmt.Sprintf("%s:possession:%d:%d", matchID, key.period, key.minute),
			MatchID:   matchID,
			Type:      event.Possession,
			Period:    key.period,
			Elapsed:   s.end,
			HomeShare: &homeShare,
		})
	}
	return out
}

func elapsed(r rawEvent) time.Duration {
	return time.Duration(r.Minute)*time.Minute + time.Duration(r.Second)*time.Second
}
