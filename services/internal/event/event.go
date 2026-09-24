package event

import (
	"fmt"
	"time"
)

type Type string

const (
	Goal          Type = "goal"
	ShotOnTarget  Type = "shot_on_target"
	ShotOffTarget Type = "shot_off_target"
	ShotBlocked   Type = "shot_blocked"
	Corner        Type = "corner"
	YellowCard    Type = "yellow_card"
	RedCard       Type = "red_card"
	Substitution  Type = "substitution"
	Foul          Type = "foul"
	Offside       Type = "offside"
	TakeOn        Type = "take_on"
)

type Side string

const (
	Home Side = "home"
	Away Side = "away"
)

// Position is on a 0–100 pitch, attacking left to right.
type Position struct {
	X, Y float64
}

type Event struct {
	ID      string
	MatchID string
	Type    Type
	Side    Side
	Period  int
	Elapsed time.Duration // match time: 46:10 in the second half is 46m10s
	Player  string
	Pos     *Position
	XG      *float64
}

func (e Event) Clock() string {
	total := int(e.Elapsed.Seconds())
	return fmt.Sprintf("%02d:%02d", total/60, total%60)
}

func (s Side) Opponent() Side {
	if s == Home {
		return Away
	}
	return Home
}
