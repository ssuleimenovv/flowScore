package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
	"github.com/ssuleimenovv/flowscore/services/internal/season"
)

// athletes is a fake espn.Client
type athletes struct {
	athlete espn.Athlete
	err     error
}

func (a athletes) Athlete(context.Context, string) (espn.Athlete, error) {
	return a.athlete, a.err
}

func TestPlayer(t *testing.T) {
	l := premierLeague()
	// He scored twice in the first round and came on in the second
	l.Games[0].Flow.Players = []season.PlayerMatch{{ID: "9", Side: event.Home, Starter: true, Minutes: 90, Goals: 2, XG: 1.3}}
	l.Games[1].Flow.Players = []season.PlayerMatch{{ID: "9", Side: event.Away, Minutes: 20}}
	espnPlayer := athletes{athlete: espn.Athlete{
		ID: "9", Name: "Haaland", Jersey: "9", Position: "forward",
		Team: &event.Team{ID: "ARS", Code: "ARS", Name: "ARS"},
	}}
	mux := http.NewServeMux()
	RegisterPlayers(mux, seasons{"eng.1": l}, espnPlayer, []string{"eng.1"})

	var p playerPage
	getLeague(t, mux, "/api/v1/players/9", http.StatusOK, &p)
	if p.Player.Name != "Haaland" || p.Player.Team == nil || p.Player.Team.ID != "ARS" {
		t.Errorf("player %+v", p.Player)
	}
	if len(p.Matches) != 2 {
		t.Fatalf("%d matches, want 2", len(p.Matches))
	}
	// The newest first
	if p.Matches[0].ID != "2" || p.Matches[0].Starter || p.Matches[0].Minutes != 20 {
		t.Errorf("came on: %+v", p.Matches[0])
	}
	if p.Matches[1].Goals != 2 || p.Matches[1].XG != 1.3 || p.Matches[1].Score == nil {
		t.Errorf("scored twice: %+v", p.Matches[1])
	}
}

func TestPlayerFails(t *testing.T) {
	notFound := athletes{err: fmt.Errorf("GET: %w", espn.ErrNotFound)}
	mux := http.NewServeMux()
	RegisterPlayers(mux, seasons{}, notFound, []string{"eng.1"})
	getLeague(t, mux, "/api/v1/players/1", http.StatusNotFound, nil)

	down := athletes{err: errors.New("503")}
	mux = http.NewServeMux()
	RegisterPlayers(mux, seasons{}, down, []string{"eng.1"})
	getLeague(t, mux, "/api/v1/players/1", http.StatusBadGateway, nil)
}
