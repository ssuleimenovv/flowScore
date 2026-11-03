package api

import (
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
	"github.com/ssuleimenovv/flowscore/services/internal/season"
)

// stageNames are the names the screens give a cup's knockout stages, by ESPN's
var stageNames = map[string]string{
	"preliminary-round":       "Предварительный раунд",
	"qualifying-round":        "Квалификация",
	"first-round":             "1-й раунд",
	"second-round":            "2-й раунд",
	"third-round":             "3-й раунд",
	"fourth-round":            "4-й раунд",
	"fifth-round":             "5-й раунд",
	"round-of-64":             "1/32",
	"round-of-32":             "1/16",
	"knockout-round-playoffs": "Стыки",
	"round-of-16":             "1/8",
	"quarterfinals":           "1/4",
	"semifinals":              "1/2",
	"final":                   "Финал",
}

type stageItem struct {
	ID   string `json:"id"`   // "round-of-16"
	Name string `json:"name"` // "1/8"
	Ties int    `json:"ties"`
}

// stagesOf lists the knockout stages a season has matches in, in the order a
// cup plays them, and the one being played: the stage of the first match
// not over, or the last stage once all are.
func stagesOf(games []season.Game) ([]stageItem, string) {
	stages := []stageItem{}
	for _, id := range season.Stages {
		if n := len(tiesOf(games, id)); n > 0 {
			stages = append(stages, stageItem{ID: id, Name: stageNames[id], Ties: n})
		}
	}
	current := ""
	for _, g := range games {
		if season.Knockout(g.Stage) {
			current = g.Stage
			if g.State != "post" {
				break
			}
		}
	}
	return stages, current
}

type tieLeg struct {
	MatchID   string    `json:"matchId"`
	KickoffAt time.Time `json:"kickoffAt"`
	State     string    `json:"state"` // pre, in or post
	Clock     string    `json:"clock"` // "104'" while in play, as ESPN shows it
}

type tieSide struct {
	Team     team   `json:"team"`
	Goals    []*int `json:"goals"`    // in each leg; null before it kicks off
	Total    *int   `json:"total"`    // over the legs played; null before the first
	Shootout *int   `json:"shootout"` // null without one
	Through  *bool  `json:"through"`  // null until the tie is decided
}

type tie struct {
	Home   tieSide  `json:"home"` // the home team of the first leg
	Away   tieSide  `json:"away"`
	Legs   []tieLeg `json:"legs"`
	Finish *string  `json:"finish"` // aet or pens, how the deciding match went; null otherwise
}

type tieList struct {
	LeagueID string `json:"leagueId"`
	Stage    string `json:"stage"`
	Items    []tie  `json:"items"`
}

// tiesOf pairs the matches of a stage into ties, in the order their first
// legs are played. A tie of two legs has both of its matches.
func tiesOf(games []season.Game, stage string) []tie {
	var ties []tie
	at := map[[2]string]int{} // a tie's index, by its teams' IDs
	for _, g := range games {
		if g.Stage != stage {
			continue
		}
		key := [2]string{g.Home.Team.ID, g.Away.Team.ID}
		slices.Sort(key[:])
		i, ok := at[key]
		if !ok {
			i = len(ties)
			at[key] = i
			ties = append(ties, tie{Home: tieSide{Team: toTeam(g.Home.Team)}, Away: tieSide{Team: toTeam(g.Away.Team)}})
		}
		t := &ties[i]
		leg := tieLeg{MatchID: g.ID, KickoffAt: g.Start, State: g.State}
		if g.State == "in" {
			leg.Clock = g.Clock
		}
		t.Legs = append(t.Legs, leg)

		// The second leg is played the other way round
		home, away := g.Home, g.Away
		if home.Team.ID != t.Home.Team.ID {
			home, away = away, home
		}
		t.Home.play(home, g)
		t.Away.play(away, g)
		if deciding(g) && g.State == "post" && g.Finish != "" {
			t.Finish = &g.Finish
		}
	}
	return ties
}

// deciding tells the match that decides a tie: its only one, or its second leg
func deciding(g season.Game) bool {
	return g.Leg == 0 || g.Leg == 2
}

// play adds a match of the tie to one of its sides
func (s *tieSide) play(m espn.Side, g season.Game) {
	var goals *int
	if n, err := strconv.Atoi(m.Score); err == nil && g.State != "pre" {
		goals = &n
		total := n
		if s.Total != nil {
			total += *s.Total
		}
		s.Total = &total
	}
	s.Goals = append(s.Goals, goals)
	if m.Shootout != nil {
		s.Shootout = m.Shootout
	}
	if deciding(g) && g.State == "post" {
		through := m.Winner
		if g.Leg == 2 {
			through = m.Through
		}
		s.Through = &through
	}
}

func tiesHandler(seasons Seasons, ids []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l, ok := league(w, r, seasons, ids)
		if !ok {
			return
		}
		stage := r.PathValue("stage")
		items := tiesOf(l.Games, stage)
		if len(items) == 0 {
			writeProblem(w, http.StatusNotFound, "Stage not found")
			return
		}
		writeJSON(w, http.StatusOK, tieList{LeagueID: l.ID, Stage: stage, Items: items})
	}
}
