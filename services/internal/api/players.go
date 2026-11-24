package api

import (
	"context"
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
	"github.com/ssuleimenovv/flowscore/services/internal/season"
)

// Athletes is where a player's club, number and line come from: espn.Client
type Athletes interface {
	Athlete(ctx context.Context, id string) (espn.Athlete, error)
}

// A player changes club in the transfer windows only
const athleteFor = 6 * time.Hour

// RegisterPlayers adds the player route of api/openapi.yaml: who the player
// is, from ESPN, and their season in our matches, in every competition the
// gateway follows.
func RegisterPlayers(mux *http.ServeMux, seasons Seasons, athletes Athletes, ids []string) {
	kept := newCache(athleteFor, func(ctx context.Context, _, id string) (espn.Athlete, error) {
		return athletes.Athlete(ctx, id)
	})
	mux.HandleFunc("GET /api/v1/players/{playerId}", playerHandler(seasons, kept, ids))
}

type playerInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Number   string `json:"number"`   // "" when ESPN gives none
	Position string `json:"position"` // goalkeeper, defender, midfielder or forward; "" for none
	Age      int    `json:"age"`
	Country  string `json:"country"`
	Team     *team  `json:"team"` // null for a player without a club
}

type playerMatch struct {
	ID        string      `json:"id"`
	LeagueID  string      `json:"leagueId"`
	League    string      `json:"league"`
	KickoffAt time.Time   `json:"kickoffAt"`
	Home      team        `json:"home"`
	Away      team        `json:"away"`
	Score     *live.Score `json:"score"`
	Side      event.Side  `json:"side"` // the player's team
	Starter   bool        `json:"starter"`
	Minutes   int         `json:"minutes"`
	Goals     int         `json:"goals"`
	Assists   int         `json:"assists"`
	Shots     int         `json:"shots"`
	OnTarget  int         `json:"onTarget"`
	XG        float64     `json:"xg"` // our model's
}

type playerPage struct {
	Player  playerInfo    `json:"player"`
	Matches []playerMatch `json:"matches"` // the ones the player got on the pitch in, newest first
}

// playerHandler answers 404 for a player ESPN does not know, and 502 when
// ESPN does not answer and nothing was kept
func playerHandler(seasons Seasons, kept *cache[espn.Athlete], ids []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("playerId")
		a, err := kept.get(r.Context(), "", id)
		switch {
		case errors.Is(err, espn.ErrNotFound):
			writeProblem(w, http.StatusNotFound, "Player not found")
			return
		case err != nil:
			writeProblem(w, http.StatusBadGateway, "ESPN did not give the player")
			return
		}
		p := playerPage{
			Player: playerInfo{
				ID: a.ID, Name: a.Name, Number: a.Jersey, Position: a.Position, Age: a.Age, Country: a.Country,
			},
			Matches: []playerMatch{},
		}
		if a.Team != nil {
			t := toTeam(*a.Team)
			p.Player.Team = &t
		}
		for _, lid := range ids {
			l, ok := seasons.League(lid)
			if !ok {
				continue
			}
			for _, g := range l.Games {
				if m, ok := playerMatchOf(g, l.ID, id); ok {
					p.Matches = append(p.Matches, m)
				}
			}
		}
		slices.SortStableFunc(p.Matches, func(a, b playerMatch) int { return b.KickoffAt.Compare(a.KickoffAt) })
		writeJSON(w, http.StatusOK, p)
	}
}

// playerMatchOf is a match as the player's page lists it; false for one the
// player did not get on the pitch in, or the season has not read
func playerMatchOf(g season.Game, leagueID, playerID string) (playerMatch, bool) {
	if g.Flow == nil {
		return playerMatch{}, false
	}
	i := slices.IndexFunc(g.Flow.Players, func(p season.PlayerMatch) bool { return p.ID == playerID })
	if i < 0 {
		return playerMatch{}, false
	}
	p := g.Flow.Players[i]
	return playerMatch{
		ID: g.ID, LeagueID: leagueID, League: nameOf(leagueID), KickoffAt: g.Start,
		Home: toTeam(g.Home.Team), Away: toTeam(g.Away.Team), Score: scoreOf(g),
		Side: p.Side, Starter: p.Starter, Minutes: p.Minutes,
		Goals: p.Goals, Assists: p.Assists, Shots: p.Shots, OnTarget: p.OnTarget, XG: p.XG,
	}, true
}
