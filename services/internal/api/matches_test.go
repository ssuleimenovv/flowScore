package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
)

func server(t *testing.T) *httptest.Server {
	t.Helper()
	store := live.NewStore()
	store.Start(event.Match{
		ID:   "m1",
		Home: event.Team{ID: "36", Code: "MCI", Name: "Manchester City"},
		Away: event.Team{ID: "1", Code: "ARS", Name: "Arsenal"},
	})

	mux := http.NewServeMux()
	Register(mux, store)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func get(t *testing.T, url string, v any) *http.Response {
	t.Helper()
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		t.Fatalf("decode %s: %v", url, err)
	}
	return res
}

func TestGetMatch(t *testing.T) {
	srv := server(t)

	var body map[string]any
	res := get(t, srv.URL+"/api/v1/matches/m1", &body)

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	home := body["home"].(map[string]any)
	if home["code"] != "MCI" || body["status"] != "live" {
		t.Errorf("body = %v", body)
	}
	if body["halftimeScore"] != nil {
		t.Errorf("halftimeScore = %v, want null", body["halftimeScore"])
	}
}

func TestEmptyListsAreNotNull(t *testing.T) {
	srv := server(t)

	var flow map[string]any
	get(t, srv.URL+"/api/v1/matches/m1/flow", &flow)
	var events map[string]any
	get(t, srv.URL+"/api/v1/matches/m1/events", &events)

	if _, ok := flow["points"].([]any); !ok {
		t.Errorf("points = %v, want []", flow["points"])
	}
	if _, ok := events["items"].([]any); !ok {
		t.Errorf("items = %v, want []", events["items"])
	}
}

func TestUnknownMatchIsProblem(t *testing.T) {
	srv := server(t)

	var body map[string]any
	res := get(t, srv.URL+"/api/v1/matches/nope/flow", &body)

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "application/problem+json" {
		t.Errorf("content type = %q", ct)
	}
	if body["status"] != float64(404) {
		t.Errorf("body = %v", body)
	}
}
