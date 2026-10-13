package api

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/season"
)

// Seasons is where the league routes read from: season.Seasons.
type Seasons interface {
	League(id string) (season.League, bool)
}

// leagueNames are the names the screens give the leagues, by ESPN's ID
var leagueNames = map[string]string{
	"eng.1": "Премьер-лига",
	"esp.1": "Ла Лига",
	"ita.1": "Серия A",
	"ger.1": "Бундеслига",
	"fra.1": "Лига 1",
}

// RegisterLeagues adds the league routes of api/openapi.yaml for the leagues
// the gateway follows, in the order the screens show them.
func RegisterLeagues(mux *http.ServeMux, seasons Seasons, ids []string) {
	mux.HandleFunc("GET /api/v1/leagues", leagueListHandler(seasons, ids))
	mux.HandleFunc("GET /api/v1/leagues/{leagueId}", withLeague(seasons, ids, leagueResponse))
	mux.HandleFunc("GET /api/v1/leagues/{leagueId}/standings", withLeague(seasons, ids, standingsResponse))
	mux.HandleFunc("GET /api/v1/leagues/{leagueId}/rounds/{round}", roundHandler(seasons, ids))
}

func nameOf(id string) string {
	return cmp.Or(leagueNames[id], id)
}

type leagueItem struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Ready bool   `json:"ready"`
}

type leagueItems struct {
	Items []leagueItem `json:"items"`
}

func leagueListHandler(seasons Seasons, ids []string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		items := make([]leagueItem, 0, len(ids))
		for _, id := range ids {
			_, ready := seasons.League(id)
			items = append(items, leagueItem{ID: id, Name: nameOf(id), Ready: ready})
		}
		writeJSON(w, http.StatusOK, leagueItems{Items: items})
	}
}

// withLeague answers 404 for a league the gateway does not follow and 503
// for one whose season is still being read, which takes about a minute.
func withLeague(seasons Seasons, ids []string, render func(season.League) any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if l, ok := league(w, r, seasons, ids); ok {
			writeJSON(w, http.StatusOK, render(l))
		}
	}
}

func league(w http.ResponseWriter, r *http.Request, seasons Seasons, ids []string) (season.League, bool) {
	id := r.PathValue("leagueId")
	if !slices.Contains(ids, id) {
		writeProblem(w, http.StatusNotFound, "League not found")
		return season.League{}, false
	}
	l, ok := seasons.League(id)
	if !ok {
		w.Header().Set("Retry-After", "30")
		writeProblem(w, http.StatusServiceUnavailable, "The season is still being read")
		return season.League{}, false
	}
	return l, true
}

// A league's records are its matches where one team's Flow was highest on
// average: the peak alone is 97–99 in most matches
const records = 3

type flowRecord struct {
	MatchID string     `json:"matchId"`
	Round   int        `json:"round"`
	Home    team       `json:"home"`
	Away    team       `json:"away"`
	Score   live.Score `json:"score"`
	Side    event.Side `json:"side"` // the team the record is of
	Flow    float64    `json:"flow"` // its Flow averaged over the match
	Peak    float64    `json:"peak"` // the highest Flow of the match, of either team
}

type leagueHeader struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Season    string       `json:"season"` // "2026/27"
	Round     int          `json:"round"`  // the round being played, or the next one
	Rounds    int          `json:"rounds"`
	Records   []flowRecord `json:"records"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

func leagueResponse(l season.League) any {
	h := leagueHeader{
		ID:        l.ID,
		Name:      nameOf(l.ID),
		Season:    fmt.Sprintf("%d/%02d", l.Table.Season, (l.Table.Season+1)%100),
		Round:     currentRound(l.Games),
		Records:   []flowRecord{},
		UpdatedAt: l.UpdatedAt,
	}
	if n := len(l.Games); n > 0 {
		h.Rounds = l.Games[n-1].Round
	}
	for _, g := range l.Games {
		score := scoreOf(g)
		if g.Flow == nil || score == nil {
			continue
		}
		side, f := event.Home, g.Flow.Home
		if g.Flow.Away > f {
			side, f = event.Away, g.Flow.Away
		}
		h.Records = append(h.Records, flowRecord{
			MatchID: g.ID, Round: g.Round, Home: toTeam(g.Home.Team), Away: toTeam(g.Away.Team),
			Score: *score, Side: side, Flow: f, Peak: g.Flow.Peak,
		})
	}
	slices.SortStableFunc(h.Records, func(a, b flowRecord) int { return cmp.Compare(b.Flow, a.Flow) })
	h.Records = h.Records[:min(records, len(h.Records))]
	return h
}

// currentRound is the round of the first match not over yet, or the last
// round once the season is.
func currentRound(games []season.Game) int {
	for _, g := range games {
		if g.State != "post" {
			return g.Round
		}
	}
	if n := len(games); n > 0 {
		return games[n-1].Round
	}
	return 0
}

type standing struct {
	Rank         int      `json:"rank"`
	Team         team     `json:"team"`
	Played       int      `json:"played"`
	Won          int      `json:"won"`
	Drawn        int      `json:"drawn"`
	Lost         int      `json:"lost"`
	GoalsFor     int      `json:"goalsFor"`
	GoalsAgainst int      `json:"goalsAgainst"`
	Points       int      `json:"points"`
	Zone         *string  `json:"zone"` // champions, europa, conference, relegation; null for none
	Flow         *float64 `json:"flow"` // the team's Flow averaged over its matches; null before any is read
}

type standingList struct {
	LeagueID string     `json:"leagueId"`
	Items    []standing `json:"items"`
}

func standingsResponse(l season.League) any {
	// Each team's average over the matches it played, by team ID
	type sum struct {
		flow float64
		n    int
	}
	flows := map[string]*sum{}
	add := func(id string, f float64) {
		s := flows[id]
		if s == nil {
			s = &sum{}
			flows[id] = s
		}
		s.flow += f
		s.n++
	}
	for _, g := range l.Games {
		if g.Flow != nil {
			add(g.Home.Team.ID, g.Flow.Home)
			add(g.Away.Team.ID, g.Flow.Away)
		}
	}

	items := make([]standing, 0, len(l.Table.Rows))
	for _, r := range l.Table.Rows {
		s := standing{
			Rank: r.Rank, Team: toTeam(r.Team),
			Played: r.Played, Won: r.Won, Drawn: r.Drawn, Lost: r.Lost,
			GoalsFor: r.For, GoalsAgainst: r.Against, Points: r.Points,
			Zone: zoneOf(r.Zone),
		}
		if f := flows[r.Team.ID]; f != nil {
			avg := f.flow / float64(f.n)
			s.Flow = &avg
		}
		items = append(items, s)
	}
	return standingList{LeagueID: l.ID, Items: items}
}

// zoneOf names what ESPN says a place leads to: "Champions League",
// "Europa League", "Conference League qualifying", "Relegation playoff"
func zoneOf(note string) *string {
	n := strings.ToLower(note)
	var z string
	switch {
	case n == "":
		return nil
	case strings.Contains(n, "relegation"):
		z = "relegation"
	case strings.Contains(n, "champions"):
		z = "champions"
	case strings.Contains(n, "europa"):
		z = "europa"
	case strings.Contains(n, "conference"):
		z = "conference"
	default:
		return nil
	}
	return &z
}

type roundMatch struct {
	ID        string           `json:"id"`
	KickoffAt time.Time        `json:"kickoffAt"`
	State     string           `json:"state"` // pre, in or post
	Clock     string           `json:"clock"` // "72'" while in play, as ESPN shows it
	Home      team             `json:"home"`
	Away      team             `json:"away"`
	Score     *live.Score      `json:"score"` // null before kick-off
	Flow      *live.FlowValues `json:"flow"`  // averaged over the match; null until it is over and read
}

type roundList struct {
	LeagueID string       `json:"leagueId"`
	Round    int          `json:"round"`
	Items    []roundMatch `json:"items"`
}

func roundHandler(seasons Seasons, ids []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ok := league(w, r, seasons, ids)
		if !ok {
			return
		}
		n, err := strconv.Atoi(r.PathValue("round"))
		if err != nil || n < 1 || len(l.Games) == 0 || n > l.Games[len(l.Games)-1].Round {
			writeProblem(w, http.StatusNotFound, "Round not found")
			return
		}
		items := []roundMatch{}
		for _, g := range l.Games {
			if g.Round != n {
				continue
			}
			m := roundMatch{
				ID: g.ID, KickoffAt: g.Start, State: g.State, Home: toTeam(g.Home.Team), Away: toTeam(g.Away.Team),
				Score: scoreOf(g),
			}
			if g.State == "in" {
				m.Clock = g.Clock
			}
			if g.Flow != nil {
				m.Flow = &live.FlowValues{Home: g.Flow.Home, Away: g.Flow.Away}
			}
			items = append(items, m)
		}
		writeJSON(w, http.StatusOK, roundList{LeagueID: l.ID, Round: n, Items: items})
	}
}

// scoreOf is the score of a match under way or over; nil before kick-off
func scoreOf(g season.Game) *live.Score {
	if g.State == "pre" {
		return nil
	}
	home, err1 := strconv.Atoi(g.Home.Score)
	away, err2 := strconv.Atoi(g.Away.Score)
	if err1 != nil || err2 != nil {
		return nil
	}
	return &live.Score{Home: home, Away: away}
}
