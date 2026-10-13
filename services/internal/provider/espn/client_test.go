package espn

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// A scoreboard as ESPN sends it, cut down to what the client reads
const board = `{
	"leagues": [{"calendar": ["2026-08-21T07:00Z", "2026-10-10T07:00Z"]}],
	"events": [
		{"id": "401879268", "name": "Leeds United at Arsenal", "date": "2026-10-10T11:30Z",
		 "status": {"displayClock": "72'", "type": {"state": "in", "name": "STATUS_SECOND_HALF"}},
		 "competitions": [{"competitors": [
			{"homeAway": "home", "score": "2", "team": {"id": "359", "abbreviation": "ARS", "displayName": "Arsenal"}},
			{"homeAway": "away", "score": "1", "team": {"id": "357", "abbreviation": "LEE", "displayName": "Leeds United"}}
		 ]}]},
		{"id": "401878776", "name": "Brentford at Aston Villa", "date": "2026-10-10T14:00Z",
		 "status": {"displayClock": "0'", "type": {"state": "pre", "name": "STATUS_SCHEDULED"}}}
	]
}`

func fake(t *testing.T, status int, body string) (*Client, *string) {
	t.Helper()
	path := new(string)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*path = r.URL.String()
		w.WriteHeader(status)
		io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	return &Client{HTTP: srv.Client(), BaseURL: srv.URL, StandingsURL: srv.URL}, path
}

func TestScoreboard(t *testing.T) {
	c, path := fake(t, http.StatusOK, board)
	games, err := c.Scoreboard(context.Background(), "eng.1")
	if err != nil {
		t.Fatal(err)
	}
	if *path != "/eng.1/scoreboard" {
		t.Errorf("asked %s", *path)
	}
	want := Game{
		ID:    "401879268",
		Name:  "Leeds United at Arsenal",
		Start: time.Date(2026, 10, 10, 11, 30, 0, 0, time.UTC),
		State: "in",
	}
	if len(games) != 2 || !games[0].Start.Equal(want.Start) || games[0].ID != want.ID ||
		games[0].State != want.State || games[1].State != "pre" {
		t.Errorf("games = %+v", games)
	}
}

func TestBoardOfADay(t *testing.T) {
	c, path := fake(t, http.StatusOK, board)
	b, err := c.Board(context.Background(), "eng.1", time.Date(2026, 10, 10, 7, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if *path != "/eng.1/scoreboard?dates=20261010" {
		t.Errorf("asked %s", *path)
	}
	// The calendar is every day of the season with a match
	if len(b.Days) != 2 || !b.Days[0].Equal(time.Date(2026, 8, 21, 7, 0, 0, 0, time.UTC)) || len(b.Games) != 2 {
		t.Errorf("board = %+v", b)
	}
}

// A table as ESPN sends it, two rows of it, cut down to what the client reads
const table = `{
	"season": {"year": 2026},
	"children": [{"standings": {"entries": [
		{"team": {"id": "382", "abbreviation": "MNC", "displayName": "Manchester City"},
		 "note": {"description": "Champions League"},
		 "stats": [{"name": "gamesPlayed", "value": 6}, {"name": "wins", "value": 6}, {"name": "ties", "value": 0},
			{"name": "losses", "value": 0}, {"name": "pointsFor", "value": 16}, {"name": "pointsAgainst", "value": 5},
			{"name": "points", "value": 18}, {"name": "rank", "value": 1}, {"name": "overall", "value": null}]},
		{"team": {"id": "363", "abbreviation": "CHE", "displayName": "Chelsea"},
		 "stats": [{"name": "gamesPlayed", "value": 6}, {"name": "points", "value": 10}, {"name": "rank", "value": 6}]}
	]}}]
}`

func TestStandings(t *testing.T) {
	c, path := fake(t, http.StatusOK, table)
	tb, err := c.Standings(context.Background(), "eng.1")
	if err != nil {
		t.Fatal(err)
	}
	if *path != "/eng.1/standings" {
		t.Errorf("asked %s", *path)
	}
	city := Standing{Rank: 1, Team: event.Team{ID: "382", Code: "MNC", Name: "Manchester City"},
		Played: 6, Won: 6, For: 16, Against: 5, Points: 18, Zone: "Champions League"}
	if tb.Season != 2026 || len(tb.Rows) != 2 || tb.Rows[0] != city {
		t.Errorf("table = %+v", tb)
	}
	// A place that leads nowhere has no note
	if tb.Rows[1].Zone != "" || tb.Rows[1].Rank != 6 {
		t.Errorf("chelsea = %+v", tb.Rows[1])
	}
}

func TestSummaryAsksForTheMatch(t *testing.T) {
	c, path := fake(t, http.StatusOK, `{"header": {}}`)
	body, err := c.Summary(context.Background(), "eng.1", "401879268")
	if err != nil || string(body) != `{"header": {}}` {
		t.Fatalf("body %q, err %v", body, err)
	}
	if *path != "/eng.1/summary?event=401879268" {
		t.Errorf("asked %s", *path)
	}
}

func TestErrorStatus(t *testing.T) {
	c, _ := fake(t, http.StatusNotFound, `{"code": 404}`)
	if _, err := c.Scoreboard(context.Background(), "xxx.1"); err == nil {
		t.Error("a 404 came back as a scoreboard")
	}
}
