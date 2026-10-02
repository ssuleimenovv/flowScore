package live

import (
	"math"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
)

// Message is the WebSocket envelope from api/openapi.yaml (WsEnvelope).
type Message struct {
	Type    string    `json:"type"`
	MatchID string    `json:"matchId"`
	Seq     int64     `json:"seq"`
	SentAt  time.Time `json:"sentAt"`
	Data    any       `json:"data"`
}

type FlowValues struct {
	Home float64 `json:"home"`
	Away float64 `json:"away"`
}

type FlowPoint struct {
	Minute int     `json:"minute"`
	Home   float64 `json:"home"`
	Away   float64 `json:"away"`
}

// FlowUpdate is the data of a flow.update message.
type FlowUpdate struct {
	Current FlowValues   `json:"current"`
	Delta10 FlowValues   `json:"delta10"`
	Point   FlowPoint    `json:"point"`
	Clock   Clock        `json:"clock"`
	Factors []FlowFactor `json:"factors"`
}

// FlowFactor is one line of "why Flow is what it is" (FlowFactor in the contract):
// how much a group of events adds to a team's impulse right now.
type FlowFactor struct {
	Side    event.Side `json:"side"`
	Key     string     `json:"key"`
	Value   float64    `json:"value"`
	Count   int        `json:"count"`   // events of the group in the last 10 minutes
	Minutes int        `json:"minutes"` // since the first of them, rounded up
	Share   *float64   `json:"share,omitempty"`
}

func toFactors(fs []flow.Factor) []FlowFactor {
	out := make([]FlowFactor, 0, len(fs)) // [] rather than null when there are none
	for _, f := range fs {
		ff := FlowFactor{
			Side:  f.Side,
			Key:   string(f.Group),
			Value: math.Round(f.Value*10) / 10,
			Count: f.Count,
			Share: f.Share,
		}
		if f.Count > 0 {
			ff.Minutes = max(1, int(math.Ceil(f.Since.Minutes())))
		}
		out = append(out, ff)
	}
	return out
}

// Clock is the match time as the contract sends it, in the REST match and in
// every flow.update. The client keeps it ticking from ObservedAt.
type Clock struct {
	ElapsedSeconds int       `json:"elapsedSeconds"`
	Period         string    `json:"period"`
	ObservedAt     time.Time `json:"observedAt"`
}

func ClockOf(period int, elapsed time.Duration, observedAt time.Time) Clock {
	return Clock{
		ElapsedSeconds: int(elapsed.Seconds()),
		Period:         periodName(period),
		ObservedAt:     observedAt,
	}
}

func periodName(period int) string {
	switch period {
	case 2:
		return "second_half"
	case 3, 4:
		return "extra_time"
	case 5:
		return "penalties"
	default:
		return "first_half"
	}
}

// StatRow is one row of the match stats (StatRow in the contract).
type StatRow struct {
	Key  string  `json:"key"`
	Home float64 `json:"home"`
	Away float64 `json:"away"`
}

// MatchStats is the data of a match.stats message: every row, not a diff,
// so a client can simply replace what it has.
type MatchStats struct {
	Stats []StatRow `json:"stats"`
}

type PersonRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MatchEvent is the data of a match.event message (MatchEvent in the contract).
type MatchEvent struct {
	ID         string          `json:"id"`
	Type       event.Type      `json:"type"`
	Side       *event.Side     `json:"side"` // nil for the whistles
	Minute     int             `json:"minute"`
	AddedTime  *int            `json:"addedTime"`
	Player     *PersonRef      `json:"player"`
	XG         *float64        `json:"xG"`
	Position   *event.Position `json:"position,omitempty"`
	FlowImpact *float64        `json:"flowImpact"`
}

// periodEnd is the regular last minute of each period.
var periodEnd = map[int]int{1: 45, 2: 90, 3: 105, 4: 120}

// minuteOf turns match time into the minute on the clock, the way the live
// badge shows it: 72:14 is 72′. Past the regular end it is added time:
// 45:49 in the first half is 45+1.
func minuteOf(period int, elapsed time.Duration) (minute int, added *int) {
	minute = int(elapsed.Minutes())
	if end, ok := periodEnd[period]; ok && minute >= end {
		extra := minute - end + 1
		return end, &extra
	}
	return minute, nil
}
