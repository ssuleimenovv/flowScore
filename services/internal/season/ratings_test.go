package season

import (
	"math"
	"testing"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
)

func TestRatings(t *testing.T) {
	game := func(home, away string, xgHome, xgAway float64) Game {
		g := Game{Game: espn.Game{Home: espn.Side{Team: event.Team{ID: home}}, Away: espn.Side{Team: event.Team{ID: away}}}}
		g.Flow = &MatchFlow{XGHome: xgHome, XGAway: xgAway}
		return g
	}
	// A beat B 2.0 to 0.5 and C 1.5 to 1.0; C–D is still to come
	games := []Game{game("A", "B", 2.0, 0.5), game("C", "A", 1.0, 1.5), {Game: espn.Game{
		Home: espn.Side{Team: event.Team{ID: "C"}}, Away: espn.Side{Team: event.Team{ID: "D"}}}}}
	r := ratings(games)
	// A: (1.5 + 0.5) / (2 + 5); B: -1.5 / (1 + 5); D has played nothing
	for id, want := range map[string]float64{"A": 2.0 / 7, "B": -1.5 / 6, "C": -0.5 / 6, "D": 0} {
		if math.Abs(r[id]-want) > 1e-9 {
			t.Errorf("rating of %s = %.4f, want %.4f", id, r[id], want)
		}
	}
}

func TestSeasonsRating(t *testing.T) {
	row := func(id string) espn.Standing { return espn.Standing{Team: event.Team{ID: id}} }
	s := &Seasons{leagues: map[string]League{"eng.1": {
		Table:   espn.Table{Rows: []espn.Standing{row("A"), row("B"), row("C")}},
		Ratings: map[string]float64{"A": 0.4, "B": -0.1},
	}}}
	if r, ok := s.Rating(event.Team{ID: "A"}, event.Team{ID: "B"}); !ok || math.Abs(r-0.5) > 1e-9 {
		t.Errorf("A – B = %.2f, %v; want 0.50", r, ok)
	}
	// C has not played yet: an average team
	if r, ok := s.Rating(event.Team{ID: "C"}, event.Team{ID: "A"}); !ok || math.Abs(r+0.4) > 1e-9 {
		t.Errorf("C – A = %.2f, %v; want -0.40", r, ok)
	}
	// X is in no league read so far
	if _, ok := s.Rating(event.Team{ID: "A"}, event.Team{ID: "X"}); ok {
		t.Error("a rating for a team of no league read")
	}
}
