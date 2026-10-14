package predict

import (
	"math"
	"testing"
)

func TestRatingFor(t *testing.T) {
	m, err := Load("../../../ai/prediction/model.json")
	if err != nil {
		t.Fatal(err)
	}
	// The rating the model's own chances came from comes back
	for _, r := range []float64{-1.3, 0, 0.5, 2} {
		if got := m.RatingFor(m.Predict(r, Situation{Period: 1})); math.Abs(got-r) > 1e-6 {
			t.Errorf("RatingFor at %.2f = %.4f", r, got)
		}
	}
	// City – Ipswich by the bookmakers: the model gives a clear favourite too
	o := m.Predict(m.RatingFor(Outcome{Home: 0.77, Draw: 0.14, Away: 0.09}), Situation{Period: 1})
	if math.Abs(o.Home-o.Away-0.68) > 1e-6 {
		t.Errorf("chances %+v, want the home win 68 points above the away win", o)
	}
}
