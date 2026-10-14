package live

import (
	"testing"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
)

// ratings is a fake season.Seasons that knows one pair of teams
type ratings map[[2]string]float64

func (r ratings) Rating(home, away event.Team) (float64, bool) {
	v, ok := r[[2]string{home.ID, away.ID}]
	return v, ok
}

func TestPublisherUsesSeasonRatings(t *testing.T) {
	model, err := predict.Load("../../../ai/prediction/model.json")
	if err != nil {
		t.Fatal(err)
	}
	preMatch := func(r Ratings, home, away event.Team) Probabilities {
		store := NewStore()
		p := NewPublisher(NewHub(), store, "m1", flow.DefaultParams())
		p.UseModel(&model)
		if r != nil {
			p.UseRatings(r)
		}
		p.Start(event.Match{ID: "m1", Home: home, Away: away})
		snap, _ := store.Get("m1")
		return snap.Prediction.PreMatch
	}

	// In 2015/16 City were far above Sunderland; this season says the other way round
	city, sunderland := event.Team{ID: "382", Name: "Manchester City"}, event.Team{ID: "366", Name: "Sunderland"}
	season := ratings{{"382", "366"}: -1.0}
	old, now := preMatch(nil, city, sunderland), preMatch(season, city, sunderland)
	if old.Home <= old.Away || now.Home >= now.Away {
		t.Errorf("City – Sunderland: %+v by the model's ratings, %+v by this season's", old, now)
	}
	// Teams the season does not know keep the model's ratings
	if got := preMatch(season, sunderland, city); got != preMatch(nil, sunderland, city) {
		t.Errorf("Sunderland – City by the model's ratings = %+v", got)
	}
}
