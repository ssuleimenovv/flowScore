package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
)

// simulateServer serves a match before kick-off, with the outcome model.
func simulateServer(t *testing.T) *httptest.Server {
	t.Helper()
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
	t.Cleanup(srv.Close)
	return srv
}

func post(t *testing.T, url, body string, v any) *http.Response {
	t.Helper()
	res, err := http.Post(url, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(v); err != nil {
		t.Fatalf("decode %s: %v", url, err)
	}
	return res
}

type simulationBody struct {
	Minute            int
	Current, Scenario live.Probabilities
}

func TestSimulateRed(t *testing.T) {
	srv := simulateServer(t)

	var body simulationBody
	res := post(t, srv.URL+"/api/v1/matches/m1/simulate", `{"reds":[{"side":"home","minute":30}]}`, &body)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", res.StatusCode)
	}
	if body.Minute != 0 || body.Scenario.Home >= body.Current.Home {
		t.Errorf("a home red at 30' should lower the home chance: %+v", body)
	}

	// No changes: the scenario is the match as it stands
	post(t, srv.URL+"/api/v1/matches/m1/simulate", `{"reds":[]}`, &body)
	if body.Scenario != body.Current {
		t.Errorf("empty scenario %+v, current %+v", body.Scenario, body.Current)
	}
}

func TestSimulateRejectsNonsense(t *testing.T) {
	srv := simulateServer(t)

	for _, body := range []string{
		`{"reds":[{"side":"left","minute":30}]}`,
		`{"reds":[{"side":"home","minute":0}]}`,
		`not json`,
	} {
		var problem map[string]any
		if res := post(t, srv.URL+"/api/v1/matches/m1/simulate", body, &problem); res.StatusCode != http.StatusBadRequest {
			t.Errorf("%s: status %d, want 400", body, res.StatusCode)
		}
	}
}

func TestSimulateWithoutModel(t *testing.T) {
	srv := server(t)

	var problem map[string]any
	res := post(t, srv.URL+"/api/v1/matches/m1/simulate", `{"reds":[]}`, &problem)
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}
