package api

import (
	"cmp"
	"context"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
	"github.com/ssuleimenovv/flowscore/services/internal/season"
)

// Squads is where a team's players come from: espn.Client
type Squads interface {
	Squad(ctx context.Context, league, team string) ([]espn.Player, error)
}

// RegisterTeams adds the team routes of api/openapi.yaml. A team is found in
// the seasons the gateway follows: the table of its league and its matches
// in every competition, cups too. Its squad is ESPN's.
func RegisterTeams(mux *http.ServeMux, seasons Seasons, squads Squads, ids []string) {
	kept := &squadCache{squads: squads, kept: map[string]keptSquad{}}
	mux.HandleFunc("GET /api/v1/teams/{teamId}", teamHandler(seasons, ids))
	mux.HandleFunc("GET /api/v1/teams/{teamId}/squad", squadHandler(seasons, kept, ids))
}

type teamLeague struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Rank int    `json:"rank"` // its place in the table
}

type teamMatch struct {
	ID        string           `json:"id"`
	LeagueID  string           `json:"leagueId"`
	League    string           `json:"league"` // "Лига чемпионов"
	KickoffAt time.Time        `json:"kickoffAt"`
	State     string           `json:"state"` // pre, in or post
	Clock     string           `json:"clock"` // "72'" while in play
	Home      team             `json:"home"`
	Away      team             `json:"away"`
	Score     *live.Score      `json:"score"`  // null before kick-off
	Flow      *live.FlowValues `json:"flow"`   // null until it is over and read
	Result    string           `json:"result"` // win, draw or loss for the team; "" until it is over
}

type teamPage struct {
	Team    team        `json:"team"`
	League  *teamLeague `json:"league"`  // the first league whose table has the team; null for none
	Matches []teamMatch `json:"matches"` // the whole season in every competition, by kick-off
}

// teamHandler answers 404 for a team none of the seasons has, and 503 while
// one of them is still being read: the team may be in it.
func teamHandler(seasons Seasons, ids []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("teamId")
		p := teamPage{Matches: []teamMatch{}}
		found, reading := false, false
		for _, lid := range ids {
			l, ok := seasons.League(lid)
			if !ok {
				reading = true
				continue
			}
			// FLOWSCORE_ESPN lists the leagues before the cups, so the first
			// table with the team is its league's
			if p.League == nil {
				for _, row := range l.Table.Rows {
					if row.Team.ID == id {
						p.League = &teamLeague{ID: l.ID, Name: nameOf(l.ID), Rank: row.Rank}
					}
				}
			}
			for _, g := range l.Games {
				if m, us, ok := teamMatchOf(g, l.ID, id); ok {
					p.Team, found = toTeam(us), true
					p.Matches = append(p.Matches, m)
				}
			}
		}
		switch {
		case found:
			slices.SortStableFunc(p.Matches, func(a, b teamMatch) int { return a.KickoffAt.Compare(b.KickoffAt) })
			writeJSON(w, http.StatusOK, p)
		case reading:
			w.Header().Set("Retry-After", "30")
			writeProblem(w, http.StatusServiceUnavailable, "The season is still being read")
		default:
			writeProblem(w, http.StatusNotFound, "Team not found")
		}
	}
}

// teamMatchOf is a game as the team's page lists it, and the team as the
// game names it; false for a game the team does not play
func teamMatchOf(g season.Game, leagueID, teamID string) (teamMatch, event.Team, bool) {
	home := g.Home.Team.ID == teamID
	if !home && g.Away.Team.ID != teamID {
		return teamMatch{}, event.Team{}, false
	}
	us, them := g.Home, g.Away
	if !home {
		us, them = them, us
	}
	m := teamMatch{
		ID: g.ID, LeagueID: leagueID, League: nameOf(leagueID), KickoffAt: g.Start, State: g.State,
		Home: toTeam(g.Home.Team), Away: toTeam(g.Away.Team), Score: scoreOf(g),
	}
	if g.State == "in" {
		m.Clock = g.Clock
	}
	if g.Flow != nil {
		m.Flow = &live.FlowValues{Home: g.Flow.Home, Away: g.Flow.Away}
	}
	if g.State == "post" && m.Score != nil {
		m.Result = resultOf(*m.Score, home, us.Winner, them.Winner)
	}
	return m, us.Team, true
}

// resultOf is how a match that is over went for one of its teams. Level
// after extra time, the shootout decides, and ESPN names its winner.
func resultOf(s live.Score, home, won, lost bool) string {
	us, them := s.Home, s.Away
	if !home {
		us, them = them, us
	}
	switch {
	case us > them, us == them && won:
		return "win"
	case us < them, us == them && lost:
		return "loss"
	}
	return "draw"
}

type squadPlayer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Number   string `json:"number"`   // "" when ESPN gives none
	Position string `json:"position"` // goalkeeper, defender, midfielder or forward; "" for none
	Age      int    `json:"age"`
	Country  string `json:"country"`
}

type squadList struct {
	TeamID string        `json:"teamId"`
	Items  []squadPlayer `json:"items"`
}

func squadHandler(seasons Seasons, kept *squadCache, ids []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("teamId")
		league, reading := leagueOfTeam(seasons, ids, id)
		switch {
		case league == "" && reading:
			w.Header().Set("Retry-After", "30")
			writeProblem(w, http.StatusServiceUnavailable, "The season is still being read")
			return
		case league == "":
			writeProblem(w, http.StatusNotFound, "Team not found")
			return
		}
		players, err := kept.get(r.Context(), league, id)
		if err != nil {
			writeProblem(w, http.StatusBadGateway, "ESPN did not give the squad")
			return
		}
		items := make([]squadPlayer, 0, len(players))
		for _, p := range players {
			items = append(items, squadPlayer{
				ID: p.ID, Name: p.Name, Number: p.Jersey, Position: p.Position, Age: p.Age, Country: p.Country,
			})
		}
		// As on the board: the goalkeepers first, each line by number
		slices.SortStableFunc(items, func(a, b squadPlayer) int {
			return cmp.Or(cmp.Compare(lineOf(a.Position), lineOf(b.Position)), cmp.Compare(numberOf(a), numberOf(b)))
		})
		writeJSON(w, http.StatusOK, squadList{TeamID: id, Items: items})
	}
}

// leagueOfTeam is a league whose season has the team, to ask ESPN under;
// "" for none, with reading true while a season is still being read
func leagueOfTeam(seasons Seasons, ids []string, teamID string) (league string, reading bool) {
	for _, lid := range ids {
		l, ok := seasons.League(lid)
		if !ok {
			reading = true
			continue
		}
		for _, g := range l.Games {
			if g.Home.Team.ID == teamID || g.Away.Team.ID == teamID {
				return lid, false
			}
		}
	}
	return "", reading
}

var lines = []string{"goalkeeper", "defender", "midfielder", "forward"}

// lineOf orders the positions back to front; a player with none goes last
func lineOf(position string) int {
	if i := slices.Index(lines, position); i >= 0 {
		return i
	}
	return len(lines)
}

// numberOf is a player's number to sort by; one without goes last
func numberOf(p squadPlayer) int {
	if n, err := strconv.Atoi(p.Number); err == nil {
		return n
	}
	return 100
}

// A squad changes in the transfer windows, and ESPN's roster is 200 KB
const squadFor = 6 * time.Hour

// squadCache keeps each team's squad for squadFor. When ESPN does not
// answer, the squad it kept last stands in, however old.
type squadCache struct {
	squads Squads
	mu     sync.Mutex
	kept   map[string]keptSquad // by team ID
}

type keptSquad struct {
	players []espn.Player
	at      time.Time
}

func (c *squadCache) get(ctx context.Context, league, team string) ([]espn.Player, error) {
	c.mu.Lock()
	k, ok := c.kept[team]
	c.mu.Unlock()
	if ok && time.Since(k.at) < squadFor {
		return k.players, nil
	}
	players, err := c.squads.Squad(ctx, league, team)
	if err != nil {
		if ok {
			return k.players, nil
		}
		return nil, err
	}
	c.mu.Lock()
	c.kept[team] = keptSquad{players: players, at: time.Now()}
	c.mu.Unlock()
	return players, nil
}
