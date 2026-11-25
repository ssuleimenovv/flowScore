package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"testing"

	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
	"github.com/ssuleimenovv/flowscore/services/internal/season"
)

func TestTeam(t *testing.T) {
	// Arsenal also played Inter in the Champions League and won on penalties
	cup := season.League{ID: "uefa.champions", Games: []season.Game{
		game("9", 0, "post", "INT", "ARS", "1", "1", nil),
	}}
	cup.Games[0].Away.Winner = true
	cup.Games[0].Away.Colors = []string{"db0007", "ffffff"}

	mux := http.NewServeMux()
	RegisterTeams(mux, seasons{"eng.1": premierLeague(), "uefa.champions": cup}, nil, []string{"eng.1", "uefa.champions"})

	var p teamPage
	getLeague(t, mux, "/api/v1/teams/ARS", http.StatusOK, &p)
	if p.Team.ID != "ARS" || p.League == nil || p.League.ID != "eng.1" || p.League.Rank != 1 {
		t.Fatalf("team %+v, league %+v", p.Team, p.League)
	}
	if !slices.Equal(p.Colors, []string{"db0007", "ffffff"}) {
		t.Errorf("colours %v", p.Colors)
	}
	if len(p.Matches) != 4 {
		t.Fatalf("%d matches, want 4", len(p.Matches))
	}

	// By kick-off: the cup's match on day 0 comes first
	pens, won, away, next := p.Matches[0], p.Matches[1], p.Matches[2], p.Matches[3]
	if pens.League != "Лига чемпионов" || pens.Result != "win" {
		t.Errorf("won on penalties: %+v", pens)
	}
	if won.Result != "win" || won.Flow == nil || won.Flow.Home != 70 {
		t.Errorf("3:0 at home: %+v", won)
	}
	if away.Result != "win" || away.Home.ID != "LEE" {
		t.Errorf("0:1 away: %+v", away)
	}
	if next.State != "pre" || next.Result != "" || next.Score != nil || next.Flow != nil {
		t.Errorf("to come: %+v", next)
	}
}

func TestTeamNotFound(t *testing.T) {
	mux := http.NewServeMux()
	RegisterTeams(mux, seasons{"eng.1": premierLeague()}, nil, []string{"eng.1"})
	getLeague(t, mux, "/api/v1/teams/XXX", http.StatusNotFound, nil)

	// A season still being read may have the team
	mux = http.NewServeMux()
	RegisterTeams(mux, seasons{}, nil, []string{"eng.1"})
	getLeague(t, mux, "/api/v1/teams/ARS", http.StatusServiceUnavailable, nil)
}

// squads is a fake espn.Client: it counts how often ESPN was asked
type squads struct {
	players []espn.Player
	err     error
	asked   int
	league  string
}

func (s *squads) Squad(_ context.Context, league, _ string) ([]espn.Player, error) {
	s.asked++
	s.league = league
	return s.players, s.err
}

func TestSquad(t *testing.T) {
	espnSquad := &squads{players: []espn.Player{
		{ID: "9", Name: "Haaland", Jersey: "9", Position: "forward"},
		{ID: "3", Name: "Dias", Jersey: "3", Position: "defender"},
		{ID: "x", Name: "Youth", Position: ""},
		{ID: "31", Name: "Ederson", Jersey: "31", Position: "goalkeeper"},
		{ID: "5", Name: "Stones", Jersey: "5", Position: "defender"},
	}}
	mux := http.NewServeMux()
	RegisterTeams(mux, seasons{"eng.1": premierLeague()}, espnSquad, []string{"eng.1"})

	var list squadList
	getLeague(t, mux, "/api/v1/teams/ARS/squad", http.StatusOK, &list)
	var order []string
	for _, p := range list.Items {
		order = append(order, p.Name)
	}
	if got := fmt.Sprint(order); got != "[Ederson Dias Stones Haaland Youth]" {
		t.Errorf("order %s", got)
	}
	if espnSquad.league != "eng.1" {
		t.Errorf("asked ESPN under %q", espnSquad.league)
	}

	// Kept: ESPN is asked once, and when it fails the kept squad stands in
	getLeague(t, mux, "/api/v1/teams/ARS/squad", http.StatusOK, nil)
	if espnSquad.asked != 1 {
		t.Errorf("ESPN asked %d times, want 1", espnSquad.asked)
	}
}

func TestSquadFails(t *testing.T) {
	espnSquad := &squads{err: errors.New("503")}
	mux := http.NewServeMux()
	RegisterTeams(mux, seasons{"eng.1": premierLeague()}, espnSquad, []string{"eng.1"})
	getLeague(t, mux, "/api/v1/teams/ARS/squad", http.StatusBadGateway, nil)
	// A team none of the seasons has: ESPN is not even asked
	getLeague(t, mux, "/api/v1/teams/XXX/squad", http.StatusNotFound, nil)
	if espnSquad.asked != 1 {
		t.Errorf("ESPN asked %d times, want 1", espnSquad.asked)
	}
}
