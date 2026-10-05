package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
)

func TestListByKickoff(t *testing.T) {
	model, err := predict.Load("../../../ai/prediction/model.json")
	if err != nil {
		t.Fatal(err)
	}
	store := live.NewStore()
	now := time.Now().UTC()
	// Announced in the order late, early: the list still starts with the earlier kick-off
	for _, m := range []event.Match{
		{ID: "late", KickoffAt: now.Add(time.Hour), Home: event.Team{Name: "Arsenal"}, Away: event.Team{Name: "Chelsea"}},
		{ID: "early", KickoffAt: now, Home: event.Team{Name: "Liverpool"}, Away: event.Team{Name: "Everton"}},
	} {
		p := live.NewPublisher(live.NewHub(), store, m.ID, flow.DefaultParams())
		p.UseModel(&model)
		p.Schedule(m)
	}

	mux := http.NewServeMux()
	Register(mux, store)
	srv := httptest.NewServer(mux)
	defer srv.Close()

	var body struct {
		Items []struct {
			ID         string
			Status     string
			Points     []any
			Factors    []any
			Prediction *live.Prediction
		}
	}
	res := get(t, srv.URL+"/api/v1/matches", &body)
	if res.StatusCode != http.StatusOK || len(body.Items) != 2 {
		t.Fatalf("status %d, items %+v", res.StatusCode, body.Items)
	}
	first := body.Items[0]
	if first.ID != "early" || body.Items[1].ID != "late" {
		t.Errorf("order: %s, %s; want early, late", first.ID, body.Items[1].ID)
	}
	if first.Status != "scheduled" || first.Prediction == nil || first.Points == nil || first.Factors == nil {
		t.Errorf("a match before kick-off: %+v", first)
	}
}
