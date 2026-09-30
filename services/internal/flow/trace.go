package flow

import (
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Sample is both teams' Flow at one whole minute of match time.
type Sample struct {
	Period     int
	At         time.Duration
	Home, Away float64
}

// periodStart is the match time each period kicks off at.
var periodStart = map[int]time.Duration{1: 0, 2: 45 * time.Minute, 3: 90 * time.Minute, 4: 105 * time.Minute, 5: 120 * time.Minute}

// Trace replays a whole match and samples Flow at every whole minute of each
// period, up to its last event. At a sample all events up to that moment are
// applied. This is what calibration fits and the reference its Python port
// (ai/calibration/flow.py) is checked against, so the two must agree on the rule.
func Trace(events []event.Event, p Params) []Sample {
	s := NewState(p)
	var out []Sample

	for i := 0; i < len(events); {
		period := events[i].Period
		end := i
		for end < len(events) && events[end].Period == period {
			end++
		}
		last := events[end-1].Elapsed

		s.StartPeriod(period, periodStart[period])
		for t := periodStart[period] + time.Minute; t <= last; t += time.Minute {
			for i < end && events[i].Elapsed <= t {
				s.Apply(events[i])
				i++
			}
			s.Advance(t)
			home, away := s.Flow()
			out = append(out, Sample{Period: period, At: t, Home: home, Away: away})
		}
		for ; i < end; i++ { // events after the last whole minute
			s.Apply(events[i])
		}
	}
	return out
}
