package espn

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// A scoreboard as ESPN sends it, cut down to what the client reads
const board = `{"events": [
	{"id": "401879268", "name": "Leeds United at Arsenal", "date": "2026-10-10T11:30Z",
	 "status": {"type": {"state": "in", "name": "STATUS_FIRST_HALF"}}},
	{"id": "401878776", "name": "Brentford at Aston Villa", "date": "2026-10-10T14:00Z",
	 "status": {"type": {"state": "pre", "name": "STATUS_SCHEDULED"}}}
]}`

func fake(t *testing.T, status int, body string) (*Client, *string) {
	t.Helper()
	path := new(string)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*path = r.URL.String()
		w.WriteHeader(status)
		io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	return &Client{HTTP: srv.Client(), BaseURL: srv.URL}, path
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
