package predict

import (
	"encoding/json"
	"math"
	"os"
	"testing"
	"time"
)

const modelPath = "../../../ai/prediction/model.json"

func load(t *testing.T) Model {
	t.Helper()
	m, err := Load(modelPath)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// The Go model must predict what Python predicts with the same model.json
// (testdata/golden.json is written by ai/prediction/golden.py).
func TestMatchesPython(t *testing.T) {
	m := load(t)
	data, err := os.ReadFile("testdata/golden.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Period    int     `json:"period"`
		Seconds   int     `json:"seconds"`
		HomeGoals int     `json:"homeGoals"`
		AwayGoals int     `json:"awayGoals"`
		HomeReds  int     `json:"homeReds"`
		AwayReds  int     `json:"awayReds"`
		Rating    float64 `json:"rating"`
		Finished  bool    `json:"finished"`
		Home      float64 `json:"home"`
		Draw      float64 `json:"draw"`
		Away      float64 `json:"away"`
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}

	worst := 0.0
	for i, c := range cases {
		got := m.Predict(c.Rating, Situation{
			Period:    c.Period,
			At:        time.Duration(c.Seconds) * time.Second,
			HomeGoals: c.HomeGoals,
			AwayGoals: c.AwayGoals,
			HomeReds:  c.HomeReds,
			AwayReds:  c.AwayReds,
			Finished:  c.Finished,
		})
		diff := max(math.Abs(got.Home-c.Home), math.Abs(got.Draw-c.Draw), math.Abs(got.Away-c.Away))
		if diff > 1e-9 {
			t.Errorf("case %d: got %+v, Python %.6f %.6f %.6f", i, got, c.Home, c.Draw, c.Away)
		}
		worst = max(worst, diff)
	}
	t.Logf("%d cases, largest difference %.1e", len(cases), worst)
}

func TestFinalWhistleIsCertain(t *testing.T) {
	m := load(t)
	got := m.Predict(-1, Situation{Period: 2, At: 94 * time.Minute, HomeGoals: 1, Finished: true})
	if got != (Outcome{Home: 1}) {
		t.Errorf("home won 1:0, got %+v", got)
	}
}

func TestStrongerTeamIsFavourite(t *testing.T) {
	m := load(t)
	start := Situation{Period: 1}
	even := m.Predict(0, start)
	if even.Home <= even.Away {
		t.Errorf("equal teams: the home side should be ahead, got %+v", even)
	}
	away := m.Predict(m.Rating("Aston Villa", "Arsenal"), start)
	if away.Away <= away.Home {
		t.Errorf("Arsenal away at Villa should be the favourite, got %+v", away)
	}
}

func TestLeadGrowsSaferWithTime(t *testing.T) {
	m := load(t)
	lead := Situation{Period: 2, HomeGoals: 1}
	var last float64
	for _, minute := range []time.Duration{50, 70, 85, 90} {
		lead.At = minute * time.Minute
		p := m.Predict(0, lead).Home
		if p <= last {
			t.Errorf("1:0 at %v: home win %.3f, not above %.3f earlier", lead.At, p, last)
		}
		last = p
	}
}

func TestUnknownTeamIsAverage(t *testing.T) {
	m := load(t)
	if r := m.Rating("Nobody FC", "Nobody FC"); r != 0 {
		t.Errorf("rating of unknown teams: got %v, want 0", r)
	}
}
