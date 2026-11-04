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

// A cup's calendar lists its stages, not its days: the board still reads,
// with its matches and no days
func TestBoardOfACup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{
			"leagues": [{"season": {"year": 2026, "startDate": "2026-07-01T04:00Z"}, "calendar": [{"label": "UEFA Champions League", "entries": [
				{"label": "League Phase", "startDate": "2026-08-27T07:00Z", "endDate": "2027-02-15T07:59Z"}]}]}],
			"events": [{"id": "401915386", "name": "Bayern Munich at Atlético Madrid", "date": "2026-11-03T20:00Z",
				"status": {"type": {"state": "pre"}},
				"competitions": [{"competitors": [
					{"homeAway": "home", "team": {"id": "1068", "abbreviation": "ATM", "displayName": "Atlético Madrid"}},
					{"homeAway": "away", "team": {"id": "132", "abbreviation": "BAY", "displayName": "Bayern Munich"}}]}]}]
		}`))
	}))
	t.Cleanup(srv.Close)

	c := &Client{HTTP: srv.Client(), BaseURL: srv.URL}
	b, err := c.Board(context.Background(), "uefa.champions", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if len(b.Days) != 0 || len(b.Games) != 1 || b.Games[0].Home.Team.Name != "Atlético Madrid" {
		t.Errorf("board %+v", b)
	}
	if want := time.Date(2026, 7, 1, 4, 0, 0, 0, time.UTC); !b.SeasonStart.Equal(want) {
		t.Errorf("season starts %v, want %v", b.SeasonStart, want)
	}
	if b.SeasonYear != 2026 {
		t.Errorf("season %d, want 2026", b.SeasonYear)
	}
}

func TestYear(t *testing.T) {
	c, path := fake(t, http.StatusOK, `{"events": []}`)
	if _, err := c.Year(context.Background(), "uefa.champions", 2026); err != nil {
		t.Fatal(err)
	}
	if want := "/uefa.champions/scoreboard?dates=2026&limit=1000"; *path != want {
		t.Errorf("asked %s, want %s", *path, want)
	}
}

// A cup's knockouts as ESPN sends them: the second leg of a tie, and a
// single match settled on penalties
func TestBoardOfKnockouts(t *testing.T) {
	c, _ := fake(t, http.StatusOK, `{"events": [
		{"id": "1", "date": "2026-02-24T17:45Z", "season": {"slug": "knockout-round-playoffs"},
			"status": {"type": {"state": "post", "name": "STATUS_FULL_TIME"}},
			"competitions": [{"leg": {"value": 2},
				"series": {"competitors": [{"id": "1068", "winner": true, "aggregateScore": 7.0}, {"id": "570", "winner": false, "aggregateScore": 4.0}]},
				"competitors": [
					{"homeAway": "home", "score": "4", "winner": true, "team": {"id": "1068", "displayName": "Atlético Madrid"}},
					{"homeAway": "away", "score": "1", "team": {"id": "570", "displayName": "Club Brugge"}}]}]},
		{"id": "2", "date": "2026-02-10T20:00Z", "season": {"slug": "quarterfinals"},
			"status": {"type": {"state": "post", "name": "STATUS_FINAL_PEN"}},
			"competitions": [{"competitors": [
				{"homeAway": "home", "score": "1", "shootoutScore": 6, "team": {"id": "114", "displayName": "Napoli"}},
				{"homeAway": "away", "score": "1", "shootoutScore": 7, "winner": true, "team": {"id": "2572", "displayName": "Como"}}]}]}]}`)
	b, err := c.Board(context.Background(), "uefa.champions", time.Time{})
	if err != nil {
		t.Fatal(err)
	}
	if len(b.Games) != 2 {
		t.Fatalf("games %+v", b.Games)
	}
	tie, pens := b.Games[0], b.Games[1]
	if tie.Stage != "knockout-round-playoffs" || tie.Leg != 2 || tie.Finish != "" {
		t.Errorf("second leg %+v", tie)
	}
	if h := tie.Home; h.Aggregate == nil || *h.Aggregate != 7 || !h.Through || !h.Winner {
		t.Errorf("Atlético %+v", h)
	}
	if a := tie.Away; a.Aggregate == nil || *a.Aggregate != 4 || a.Through {
		t.Errorf("Brugge %+v", a)
	}
	if pens.Leg != 0 || pens.Finish != "pens" || pens.Home.Aggregate != nil {
		t.Errorf("shootout %+v", pens)
	}
	if h, a := pens.Home, pens.Away; h.Shootout == nil || *h.Shootout != 6 || a.Shootout == nil || *a.Shootout != 7 || !a.Winner || h.Winner {
		t.Errorf("shootout %+v : %+v", h, a)
	}
}
