package live

import (
	"math"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// tally is one team's side of the match stats.
type tally struct {
	shots, onTarget, keyPasses, tackles int
	xg                                  float64
}

// Stats counts the "Статистика матча" rows from the events as they come,
// so the numbers always agree with the timeline.
type Stats struct {
	teams     map[event.Side]*tally
	homeShare float64 // sum of the home team's share of the ball, one per minute
	minutes   int
}

func NewStats() *Stats {
	return &Stats{teams: map[event.Side]*tally{event.Home: {}, event.Away: {}}}
}

// Add counts the event and reports whether any row changed.
func (s *Stats) Add(e event.Event) bool {
	if e.Type == event.Possession {
		if e.HomeShare == nil {
			return false
		}
		s.homeShare += *e.HomeShare
		s.minutes++
		return true
	}

	t, ok := s.teams[e.Side]
	if !ok || e.OwnGoal { // an own goal is on the scoreboard, but it was no shot
		return false
	}
	switch e.Type {
	case event.Goal, event.ShotOnTarget:
		t.shots++
		t.onTarget++
	case event.ShotOffTarget, event.ShotBlocked:
		t.shots++
	case event.KeyPass:
		t.keyPasses++
	case event.Tackle:
		t.tackles++
	default:
		return false
	}
	if e.XG != nil {
		t.xg += *e.XG
	}
	return true
}

// Rows returns the stats in the order of the Match board.
func (s *Stats) Rows() []StatRow {
	home, away := s.teams[event.Home], s.teams[event.Away]

	possession := 50.0 // nobody has had the ball yet
	if s.minutes > 0 {
		possession = math.Round(s.homeShare / float64(s.minutes) * 100)
	}

	return []StatRow{
		{Key: "possession", Home: possession, Away: 100 - possession},
		{Key: "shots", Home: float64(home.shots), Away: float64(away.shots)},
		{Key: "shots_on_target", Home: float64(home.onTarget), Away: float64(away.onTarget)},
		{Key: "xg", Home: round2(home.xg), Away: round2(away.xg)},
		{Key: "key_passes", Home: float64(home.keyPasses), Away: float64(away.keyPasses)},
		{Key: "pressure_tackles", Home: float64(home.tackles), Away: float64(away.tackles)},
	}
}

// round2 keeps two decimals, as xG is shown: 1.84
func round2(x float64) float64 {
	return math.Round(x*100) / 100
}
