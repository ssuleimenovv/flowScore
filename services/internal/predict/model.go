// Package predict is the outcome model of docs/PREDICTION.md: the chances of a
// home win, a draw and an away win from any moment of a match. The numbers come
// from ai/prediction/fit.py, which writes them to ai/prediction/model.json.
package predict

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
)

// Model is model.json: the goal rate weights, the draw fix and team ratings.
type Model struct {
	Name      string             `json:"model"`
	Intercept float64            `json:"intercept"`
	Weights   Weights            `json:"weights"`
	Draw      float64            `json:"draw"`
	LeadCap   int                `json:"leadCap"`
	MaxGoals  int                `json:"maxGoals"`
	Ratings   map[string]float64 `json:"ratings"` // by team name, after the last season
}

// Weights of the log goal rate, one per feature of docs/PREDICTION.md, section 1.
type Weights struct {
	Home   float64 `json:"home"`
	Rating float64 `json:"rating"`
	Red    float64 `json:"red"`
	Lead   float64 `json:"lead"`
}

func Load(path string) (Model, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Model{}, err
	}
	var m Model
	if err := json.Unmarshal(data, &m); err != nil {
		return Model{}, fmt.Errorf("%s: %w", path, err)
	}
	if m.MaxGoals <= 0 || m.Name == "" {
		return Model{}, fmt.Errorf("%s: not an outcome model", path)
	}
	return m, nil
}

// Situation is what a viewer knows at one moment of the match.
type Situation struct {
	Period               int
	At                   time.Duration // match time; the second half starts at 45:00
	HomeGoals, AwayGoals int
	HomeReds, AwayReds   int
	Finished             bool
}

// Outcome holds the chances of the three results; they add up to 1.
type Outcome struct {
	Home, Draw, Away float64
}

// The time a match is expected to last, stoppage included, as in dataset.py:
// a viewer does not know the real stoppage time, so neither does the model.
const (
	firstHalfEnd  = 47.0
	secondHalfEnd = 93.0
	secondHalf    = 48.0 // minutes the second half is expected to last
	lastMinute    = 0.5  // a second half that runs over still has a moment left
)

// Rating is the home team's rating minus the away team's. A team the model
// has not seen counts as average, 0.
func (m Model) Rating(home, away string) float64 {
	return m.Ratings[home] - m.Ratings[away]
}

// Predict returns the chances of each result from this moment on. rating is
// Rating of the two teams.
func (m Model) Predict(rating float64, s Situation) Outcome {
	left := remaining(s)
	lead := s.HomeGoals - s.AwayGoals
	home := m.rate(1, rating, s.AwayReds-s.HomeReds, lead) * left
	away := m.rate(0, -rating, s.HomeReds-s.AwayReds, -lead) * left

	// Two independent counts give too few draws: the draw is weighted by e^Draw
	o := m.outcome(home, away, lead)
	o.Draw *= math.Exp(m.Draw)
	return o.normalized()
}

// rate is one team's expected goals per minute.
func (m Model) rate(home, rating float64, red, lead int) float64 {
	lead = max(-m.LeadCap, min(m.LeadCap, lead))
	w := m.Weights
	return math.Exp(m.Intercept + w.Home*home + w.Rating*rating + w.Red*float64(red) + w.Lead*float64(lead))
}

// outcome adds the further goals of both teams to the current home lead,
// over every pair of counts up to MaxGoals.
func (m Model) outcome(homeMean, awayMean float64, lead int) Outcome {
	home, away := poisson(homeMean, m.MaxGoals), poisson(awayMean, m.MaxGoals)
	var o Outcome
	for i, ph := range home {
		for j, pa := range away {
			switch final := lead + i - j; {
			case final > 0:
				o.Home += ph * pa
			case final == 0:
				o.Draw += ph * pa
			default:
				o.Away += ph * pa
			}
		}
	}
	return o.normalized()
}

func (o Outcome) normalized() Outcome {
	sum := o.Home + o.Draw + o.Away
	return Outcome{Home: o.Home / sum, Draw: o.Draw / sum, Away: o.Away / sum}
}

// poisson returns P(0..n goals) for this mean, each term from the one before.
func poisson(mean float64, n int) []float64 {
	p := make([]float64, n+1)
	p[0] = math.Exp(-mean)
	for k := 1; k <= n; k++ {
		p[k] = p[k-1] * mean / float64(k)
	}
	return p
}

// remaining is how many minutes the model expects are left.
func remaining(s Situation) float64 {
	if s.Finished {
		return 0
	}
	minute := s.At.Minutes()
	if s.Period <= 1 {
		return max(firstHalfEnd-minute, 0) + secondHalf
	}
	return max(secondHalfEnd-minute, lastMinute)
}
