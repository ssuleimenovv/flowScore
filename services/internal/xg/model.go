// Package xg is our own xG model of docs/FLOW.md, section 5: the chance that
// a shot goes in, from where it was taken and how. The numbers come from
// ai/xg/fit.py, which writes them to ai/xg/model.json.
package xg

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

// Situation is what led to the shot, in the words a live feed also has.
type Situation string

const (
	OpenPlay Situation = "open_play"
	SetPiece Situation = "set_piece" // play that started with a corner or a free kick
	FreeKick Situation = "free_kick" // a shot straight from a free kick
	Penalty  Situation = "penalty"
)

// Shot is what the model knows about a shot. X and Y are StatsBomb's yards:
// the pitch is 120 × 80 and the goal is on the line x = 120.
type Shot struct {
	X, Y      float64
	Header    bool
	Situation Situation
}

// Model is model.json: a logistic regression on the features of fit.py.
type Model struct {
	Name      string  `json:"name"`
	Intercept float64 `json:"intercept"`
	Weights   Weights `json:"weights"`
}

// Weights of the log-odds of a goal, one per feature.
type Weights struct {
	Distance float64 `json:"distance"`
	Angle    float64 `json:"angle"`
	Header   float64 `json:"header"`
	SetPiece float64 `json:"set_piece"`
	FreeKick float64 `json:"free_kick"`
	Penalty  float64 `json:"penalty"`
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
	if m.Name == "" || m.Weights.Distance == 0 {
		return Model{}, fmt.Errorf("%s: not an xG model", path)
	}
	return m, nil
}

// The goal: its middle and half its width, in yards
const (
	goalX    = 120.0
	goalY    = 40.0
	halfGoal = 4.0
)

// XG is the chance that the shot goes in, between 0 and 1.
func (m Model) XG(s Shot) float64 {
	dx := goalX - s.X
	dy := s.Y - goalY
	distance := math.Hypot(dx, dy)
	// The angle between the lines to the two posts, as in fit.py
	angle := math.Atan2(2*halfGoal*dx, dx*dx+dy*dy-halfGoal*halfGoal)

	w := m.Weights
	z := m.Intercept + w.Distance*distance + w.Angle*angle
	if s.Header {
		z += w.Header
	}
	switch s.Situation {
	case SetPiece:
		z += w.SetPiece
	case FreeKick:
		z += w.FreeKick
	case Penalty:
		z += w.Penalty
	}
	return 1 / (1 + math.Exp(-z))
}
