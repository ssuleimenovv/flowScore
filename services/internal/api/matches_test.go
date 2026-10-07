package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
)

func server(t *testing.T) *httptest.Server {
	t.Helper()
	store := live.NewStore()
	store.Start(event.Match{
		ID:   "m1",
		Home: event.Team{ID: "36", Code: "MCI", Name: "Manchester City"},
		Away: event.Team{ID: "1", Code: "ARS", Name: "Arsenal"},
	}, 0)

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
	if _, ok := flow["factors"].([]any); !ok {
		t.Errorf("factors = %v, want []", flow["factors"])
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

func TestInsight(t *testing.T) {
	model, err := predict.Load("../../../ai/prediction/model.json")
	if err != nil {
		t.Fatal(err)
	}
	store := live.NewStore()
	p := live.NewPublisher(live.NewHub(), store, "m1", flow.DefaultParams())
	p.UseModel(&model)
	p.Schedule(event.Match{ID: "m1", Home: event.Team{Name: "Manchester City"}, Away: event.Team{Name: "Arsenal"}})

	mux := http.NewServeMux()
	Register(mux, store)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	var body struct {
		Prediction  live.Prediction
		Explanation any
		Seq         int64
	}
	res := get(t, srv.URL+"/api/v1/matches/m1/insight", &body)
	pr := body.Prediction
	if res.StatusCode != http.StatusOK || pr.Model != model.Name || pr.Current != pr.PreMatch {
		t.Errorf("status %d, prediction %+v", res.StatusCode, pr)
	}
	if sum := pr.Current.Home + pr.Current.Draw + pr.Current.Away; sum != 100 {
		t.Errorf("chances add up to %d", sum)
	}
	// The agent has not answered yet: the field is there, as null
	if body.Explanation != nil {
		t.Errorf("explanation = %v, want null", body.Explanation)
	}
}

// Without a model there is no prediction, and the contract has no insight without one
func TestInsightWithoutModel(t *testing.T) {
	srv := server(t)

	var body map[string]any
	res := get(t, srv.URL+"/api/v1/matches/m1/insight", &body)
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}
