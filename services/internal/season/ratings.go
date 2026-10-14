package season

import "github.com/ssuleimenovv/flowscore/services/internal/event"

// A team's rating counts this many average games before its own, as in
// ai/prediction/dataset.py: after 5 matches it is half the team's record
// and half "an average team"
const priorGames = 5

// ratings are the team's strength after the matches read so far, by team
// ID: each team's xG difference summed over its matches, over the matches
// plus priorGames. It is what the outcome model was fitted on
func ratings(games []Game) map[string]float64 {
	type record struct {
		diff float64
		n    int
	}
	teams := map[string]*record{}
	add := func(id string, diff float64) {
		r := teams[id]
		if r == nil {
			r = &record{}
			teams[id] = r
		}
		r.diff += diff
		r.n++
	}
	for _, g := range games {
		if g.Flow == nil {
			continue
		}
		diff := g.Flow.XGHome - g.Flow.XGAway
		add(g.Home.Team.ID, diff)
		add(g.Away.Team.ID, -diff)
	}
	out := make(map[string]float64, len(teams))
	for id, r := range teams {
		out[id] = r.diff / float64(r.n+priorGames)
	}
	return out
}

// Rating is the home team's rating minus the away team's, from the league
// both play in, and false while no such league has been read. A team with
// no match read yet counts as average, 0.
func (s *Seasons) Rating(home, away event.Team) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, l := range s.leagues {
		if l.plays(home.ID) && l.plays(away.ID) {
			return l.Ratings[home.ID] - l.Ratings[away.ID], true
		}
	}
	return 0, false
}

// plays reports whether the team is in the league's table
func (l League) plays(teamID string) bool {
	for _, r := range l.Table.Rows {
		if r.Team.ID == teamID {
			return true
		}
	}
	return false
}
