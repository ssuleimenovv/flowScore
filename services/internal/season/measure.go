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

	if model == nil {
		return f
	}
	for _, e := range events {
		shot := e.Type == event.ShotOnTarget || e.Type == event.ShotOffTarget ||
			e.Type == event.ShotBlocked || (e.Type == event.Goal && !e.OwnGoal)
		if !shot || e.Pos == nil {
			continue
		}
		v := model.XG(xg.Shot{
			X:         e.Pos.X / 100 * pitchLength,
			Y:         e.Pos.Y / 100 * pitchWidth,
			Header:    e.Header,
			Situation: xg.Situation(e.Situation),
		})
		if e.Side == event.Home {
			f.XGHome += v
		} else {
			f.XGAway += v
		}
	}
	return f
}
