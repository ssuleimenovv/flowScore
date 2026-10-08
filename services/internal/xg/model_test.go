package xg

import (
	"encoding/json"
	"math"
	"os"
	"testing"
)

const modelPath = "../../../ai/xg/model.json"

func load(t *testing.T) Model {
	t.Helper()
	m, err := Load(modelPath)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// The Go model must give what Python gives with the same model.json
// (testdata/golden.json is written by ai/xg/golden.py).
func TestMatchesPython(t *testing.T) {
	m := load(t)
	data, err := os.ReadFile("testdata/golden.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		X         float64   `json:"x"`
		Y         float64   `json:"y"`
		Header    bool      `json:"header"`
		Situation Situation `json:"situation"`
		XG        float64   `json:"xg"`
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}

	worst := 0.0
	for i, c := range cases {
		got := m.XG(Shot{X: c.X, Y: c.Y, Header: c.Header, Situation: c.Situation})
		diff := math.Abs(got - c.XG)
		if diff > 1e-12 {
			t.Errorf("case %d %+v: got %.6f, Python %.6f", i, c, got, c.XG)
		}
		worst = max(worst, diff)
	}
	t.Logf("%d cases, worst difference %.1e", len(cases), worst)
}

// What a fan knows about shots, the model must know too
func TestShotsMakeSense(t *testing.T) {
	m := load(t)
	near := m.XG(Shot{X: 112, Y: 40})
	far := m.XG(Shot{X: 95, Y: 40})
	wide := m.XG(Shot{X: 112, Y: 20})
	header := m.XG(Shot{X: 112, Y: 40, Header: true})
	penalty := m.XG(Shot{X: 108, Y: 40, Situation: Penalty})

	if !(near > far && near > wide && near > header) {
		t.Errorf("near %.3f, far %.3f, wide %.3f, header %.3f", near, far, wide, header)
	}
	if penalty < 0.7 || penalty > 0.85 {
		t.Errorf("penalty %.3f, want about 0.8, as scored in the season", penalty)
	}
}
