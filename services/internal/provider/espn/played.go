package espn

import (
	"cmp"
	"slices"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Played reads a finished match from its summary at once: its plays and its
// share of the ball, in match order, as the live stream would have sent them.
// ESPN keeps only the share of the whole match, so every minute gets that.
func Played(body []byte) (Summary, []event.Event, error) {
	s, err := ParseSummary(body)
	if err != nil {
		return Summary{}, nil, err
	}
	plays, err := ParsePlays(body, s.Match)
	if err != nil {
		return Summary{}, nil, err
	}
	end := 2 * secondHalf
	for _, p := range plays {
		end = max(end, p.Elapsed)
	}
	var b ball
	events := append(plays, b.next(s.Match.ID, 2, end, s.Possession)...)
	slices.SortStableFunc(events, func(a, b event.Event) int {
		return cmp.Or(cmp.Compare(a.Period, b.Period), cmp.Compare(a.Elapsed, b.Elapsed))
	})
	return s, events, nil
}
