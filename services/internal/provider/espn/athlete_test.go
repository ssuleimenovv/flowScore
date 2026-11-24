package espn

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestParseAthlete(t *testing.T) {
	// Erling Haaland, 2026/27
	body, err := os.ReadFile("testdata/athlete-253989.json")
	if err != nil {
		t.Fatal(err)
	}
	a, err := ParseAthlete(body)
	if err != nil {
		t.Fatal(err)
	}
	if a.Name != "Erling Haaland" || a.Jersey != "9" || a.Position != "forward" || a.Age != 26 || a.Country != "Norway" {
		t.Errorf("athlete %+v", a)
	}
	if a.Team == nil || a.Team.ID != "382" || a.Team.Code != "MNC" {
		t.Errorf("team %+v", a.Team)
	}
}

func TestAthleteNotFound(t *testing.T) {
	// ESPN answers 400 for a player ID it does not know
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()
	c := &Client{HTTP: srv.Client(), AthleteURL: srv.URL}
	if _, err := c.Athlete(context.Background(), "999999999"); !errors.Is(err, ErrNotFound) {
		t.Errorf("err %v, want ErrNotFound", err)
	}
}
