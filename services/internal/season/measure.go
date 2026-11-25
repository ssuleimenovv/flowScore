package season

import (
	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/xg"
)

// MatchFlow is what a played match comes to in Flow and in our xG.
type MatchFlow struct {
	Home, Away     float64 // each team's Flow, averaged over the minutes of the match
	Peak           float64 // the highest Flow either team reached
	PeakSide       event.Side
	PeakMinute     int
	XGHome, XGAway float64 // our model's, over the shots of each team
	// XGKnown is false when the xG above says nothing: without a model, or
	// with a shot the source gave no place on the pitch. ESPN's commentary of
	// a cup match has the goals only, without their places.
	XGKnown bool
	Players []PlayerMatch // who got on the pitch and what they did; nil without lineups
}

// event.Position is 0–100; the xG model measures in StatsBomb's yards
const (
	pitchLength = 120.0
	pitchWidth  = 80.0
)

// Measure runs a whole match, its events in match order, through the Flow
// engine minute by minute, as calibration does (flow.Trace).
func Measure(events []event.Event, p flow.Params, model *xg.Model) MatchFlow {
	var f MatchFlow
	samples := flow.Trace(events, p)
	for _, s := range samples {
		f.Home += s.Home
		f.Away += s.Away
		minute := int(s.At.Minutes())
		if s.Home > f.Peak {
			f.Peak, f.PeakSide, f.PeakMinute = s.Home, event.Home, minute
		}
		if s.Away > f.Peak {
			f.Peak, f.PeakSide, f.PeakMinute = s.Away, event.Away, minute
		}
	}
	if n := float64(len(samples)); n > 0 {
		f.Home /= n
		f.Away /= n
	}

	f.XGKnown = model != nil && shotsPlaced(events)
	if !f.XGKnown {
		return f
	}
	for _, e := range events {
		v, ok := shotXG(e, model)
		if !ok {
			continue
		}
		if e.Side == event.Home {
			f.XGHome += v
		} else {
			f.XGAway += v
		}
	}

	return f
}

// isShot is a shot of the team that look it: an own goal in none
func isShot(e event.Event) bool {
	return e.Type == event.ShotOnTarget || e.Type == event.ShotOffTarget ||
		e.Type == event.ShotBlocked || (e.Type == event.Goal && !e.OwnGoal)
}

// shotsPlaced reports whether every shot of the match has its place on the
// pitch: a match without a shot has nothing to measure wrong either
func shotsPlaced(events []event.Event) bool {
	for _, e := range events {
		if isShot(e) && e.Pos == nil {
			return false
		}
	}
	return true
}

// shotXG is our model's xG of a shot; false for anything else, for a shot
// with no place on the pitch, and without a model
func shotXG(e event.Event, model *xg.Model) (float64, bool) {
	if model == nil || !isShot(e) || e.Pos == nil {
		return 0, false
	}
	return model.XG(xg.Shot{
		X:         e.Pos.X / 100 * pitchLength,
		Y:         e.Pos.Y / 100 * pitchWidth,
		Header:    e.Header,
		Situation: xg.Situation(e.Situation),
	}), true
}
