package live

import (
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
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
	Current FlowValues `json:"current"`
	Delta10 FlowValues `json:"delta10"`
	Point   FlowPoint  `json:"point"`
}

type PersonRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MatchEvent is the data of a match.event message (MatchEvent in the contract).
type MatchEvent struct {
	ID         string          `json:"id"`
	Type       event.Type      `json:"type"`
	Side       event.Side      `json:"side"`
	Minute     int             `json:"minute"`
	AddedTime  *int            `json:"addedTime"`
	Player     *PersonRef      `json:"player"`
	XG         *float64        `json:"xG"`
	Position   *event.Position `json:"position,omitempty"`
	FlowImpact *float64        `json:"flowImpact"`
}

// periodEnd is the regular last minute of each period.
var periodEnd = map[int]int{1: 45, 2: 90, 3: 105, 4: 120}

// minuteOf turns match time into the football minute: 45:49 in the first half
// is the 46th minute, shown as 45+1.
func minuteOf(period int, elapsed time.Duration) (minute int, added *int) {
	minute = int(elapsed.Minutes()) + 1
	if end, ok := periodEnd[period]; ok && minute > end {
		extra := minute - end
		return end, &extra
	}
	return minute, nil
}
