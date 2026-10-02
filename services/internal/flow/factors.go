package flow

import (
	"math"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Group is one kind of event in the explanation of Flow (docs/FLOW.md, section 6).
type Group string

const (
	Shots         Group = "shots"
	Goals         Group = "goals"
	KeyPasses     Group = "key_passes"
	Possession    Group = "possession"
	Corners       Group = "corners"
	Cards         Group = "cards"     // yellow cards the opponent got
	RedCards      Group = "red_cards" // the team's own
	Substitutions Group = "substitutions"
	Other         Group = "other" // events weighted after calibration
)

// groups is the order factors are listed in, so updates are stable.
var groups = []Group{Shots, Goals, KeyPasses, Possession, Corners, Cards, RedCards, Substitutions, Other}

func groupOf(t event.Type) Group {
	switch t {
	case event.ShotOnTarget, event.ShotOffTarget, event.ShotBlocked:
		return Shots
	case event.Goal:
		return Goals
	case event.KeyPass:
		return KeyPasses
	case event.Possession:
		return Possession
	case event.Corner:
		return Corners
	case event.YellowCard:
		return Cards
	case event.RedCard:
		return RedCards
	case event.Substitution:
		return Substitutions
	default:
		return Other
	}
}

// Factor is what one group of events adds to a team's impulse right now.
type Factor struct {
	Side  event.Side
	Group Group
	Value float64 // the decayed sum of the group's weights; a team's factors add up to its S
	// For the label "4 удара за последние 7 минут": the group's events in the
	// last 10 minutes of the period and how long ago the first of them was
	Count int
	Since time.Duration
	Share *float64 // possession only: the team's share of the ball over those minutes
}

// recentWindow is how far back the labels count events.
const recentWindow = 10 * time.Minute

// minFactor leaves out groups that no longer move Flow in a visible way.
const minFactor = 0.5

// mark is one weighted event of the last minutes, kept for the labels.
type mark struct {
	side  event.Side
	group Group
	at    time.Duration
	share float64
}

// Factors explains the current impulse of both teams group by group.
func (s *State) Factors() []Factor {
	var out []Factor
	for _, side := range []event.Side{event.Home, event.Away} {
		for _, g := range groups {
			v := s.parts[side][g]
			if math.Abs(v) < minFactor {
				continue
			}
			f := Factor{Side: side, Group: g, Value: v}

			oldest, shares := s.at, 0.0
			for _, m := range s.recent {
				if m.side == side && m.group == g && m.at >= s.at-recentWindow {
					f.Count++
					oldest = min(oldest, m.at)
					shares += m.share
				}
			}
			f.Since = s.at - oldest
			if g == Possession && f.Count > 0 {
				share := shares / float64(f.Count)
				f.Share = &share
			}
			out = append(out, f)
		}
	}
	return out
}

// remember keeps the event for the labels and forgets what left the window.
func (s *State) remember(m mark) {
	s.recent = append(s.recent, m)
	keep := s.recent[:0]
	for _, r := range s.recent {
		if r.at >= s.at-recentWindow {
			keep = append(keep, r)
		}
	}
	s.recent = keep
}
