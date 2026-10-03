package predict

import (
	"slices"
	"time"
)

// Red is a sending-off the simulator adds later in the match.
type Red struct {
	Home   bool // the home team loses a player, else the away team
	Minute int  // on the match clock; one already past counts from now
}

// Simulate returns the chances if the reds happen. The rest of the match is
// cut at each of them, and every piece has the goal rates of the players then
// on the pitch: 10 against 11 from the 78th minute on is 12 minutes of it, not
// the whole match. Without reds it is Predict.
func (m Model) Simulate(rating float64, s Situation, reds []Red) Outcome {
	lead := s.HomeGoals - s.AwayGoals
	reds = slices.SortedFunc(slices.Values(reds), func(a, b Red) int { return a.Minute - b.Minute })

	var home, away float64
	left := remaining(s)
	now := s
	for _, r := range reds {
		span := max(left-remaining(s.at(r.Minute)), 0)
		home += m.rate(1, rating, now.AwayReds-now.HomeReds, lead) * span
		away += m.rate(0, -rating, now.HomeReds-now.AwayReds, -lead) * span
		left -= span
		if r.Home {
			now.HomeReds++
		} else {
			now.AwayReds++
		}
	}
	home += m.rate(1, rating, now.AwayReds-now.HomeReds, lead) * left
	away += m.rate(0, -rating, now.HomeReds-now.AwayReds, -lead) * left
	return m.chances(home, away, lead)
}

// at is the same match at a later minute on the clock. A minute past the 45th
// is in the second half: the simulator does not add reds in first-half
// stoppage time.
func (s Situation) at(minute int) Situation {
	if s.Period == 1 && minute >= 45 {
		s.Period = 2
	}
	s.At = time.Duration(minute) * time.Minute
	return s
}
