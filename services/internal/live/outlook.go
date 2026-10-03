package live

import "github.com/ssuleimenovv/flowscore/services/internal/predict"

// Outlook is what the outcome model knows about the match right now: enough
// for the simulator to predict it again with changes.
type Outlook struct {
	model     *predict.Model
	rating    float64
	situation predict.Situation
}

// Simulate returns the chances as the match stands and with the reds added.
func (o Outlook) Simulate(reds []predict.Red) (now, scenario Probabilities) {
	now = toPercents(o.model.Predict(o.rating, o.situation))
	scenario = toPercents(o.model.Simulate(o.rating, o.situation, reds))
	return now, scenario
}

// Minute is the match minute on the clock, as the live badge shows it.
func (o Outlook) Minute() int {
	return int(o.situation.At.Minutes())
}
