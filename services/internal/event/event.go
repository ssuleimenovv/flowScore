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
	Halftime      Type = "halftime" // the whistles have no side and no flow weight
	Fulltime      Type = "fulltime"
	// Internal events feed the match stats and are never shown in the timelines
	Possession Type = "possession" // one per match minute, with HomeShare
	KeyPass    Type = "key_pass"   // a pass that set up a shot
	Tackle     Type = "tackle"
)

// Internal reports whether the event only feeds the stats.
func (t Type) Internal() bool {
	return t == Possession || t == KeyPass || t == Tackle
}

type Side string

const (
	Home Side = "home"
	Away Side = "away"
)

// Position is on a 0–100 pitch, attacking left to right.
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Event struct {
	ID        string
	MatchID   string
	Type      Type
	Side      Side
	Period    int
	Elapsed   time.Duration // match time: 46:10 in the second half is 46m10s
	Player    string        // the name fans know: "Sergio Aguero", not the full legal name
	PlayerID  string
	Pos       *Position
	XG        *float64
	HomeShare *float64 // homeshare is set on possession events: the home team's share of the ball
	// during the minute, from 0 to 1
}

func (e Event) Clock() string {
	return FormatClock(e.Elapsed)
}

// FormatClock renders match time as mm:ss
func FormatClock(d time.Duration) string {
	total := int(d.Seconds())
	return fmt.Sprintf("%02d:%02d", total/60, total%60)
}

func (s Side) Opponent() Side {
	if s == Home {
		return Away
	}
	return Home
}
