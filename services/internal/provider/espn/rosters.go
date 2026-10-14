package espn

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// spots are ESPN's positions, Opta's names: "CD-L" is the left one of the
// centre-backs. The names follow the formation rather than the role: in a
// 4-2-3-1 the two holding midfielders are "LM" and "RM", in a 3-4-2-1 those
// are the wing-backs, and the two behind the lone "F" are "CF-L" and "CF-R".
// So "F" stands a little ahead of the other forwards, and each line of the
// formation still gets the right players.
var spots = map[string]event.Spot{
	"G":  {Line: 0},
	"LB": {Line: 1, Flank: -2}, "CD-L": {Line: 1, Flank: -1}, "CD": {Line: 1}, "CD-R": {Line: 1, Flank: 1}, "RB": {Line: 1, Flank: 2},
	"LWB": {Line: 1.5, Flank: -2}, "RWB": {Line: 1.5, Flank: 2},
	"DM-L": {Line: 2, Flank: -1}, "DM": {Line: 2}, "DM-R": {Line: 2, Flank: 1},
	"LM": {Line: 2.5, Flank: -2}, "CM-L": {Line: 2.5, Flank: -1}, "CM": {Line: 2.5}, "CM-R": {Line: 2.5, Flank: 1}, "RM": {Line: 2.5, Flank: 2},
	"AM-L": {Line: 3, Flank: -1}, "AM": {Line: 3}, "AM-R": {Line: 3, Flank: 1},
	"LW": {Line: 3.5, Flank: -2}, "RW": {Line: 3.5, Flank: 2},
	"LF": {Line: 4, Flank: -2}, "CF-L": {Line: 4, Flank: -1}, "CF-R": {Line: 4, Flank: 1}, "RF": {Line: 4, Flank: 2},
	"F": {Line: 4.2}, "CF": {Line: 4.2},
}

// A position ESPN has not shown us yet goes in the middle of the pitch
var unknownSpot = event.Spot{Line: 2.5}

// ParseRosters reads both teams' players from a summary: nil before ESPN
// has them, about an hour before kick-off. During the match it is updated
// with the substitutions.
func ParseRosters(body []byte) (*event.Lineups, error) {
	var raw struct {
		Rosters []struct {
			HomeAway  string `json:"homeAway"`
			Formation string `json:"formation"` // "4-2-3-1"
			Roster    []struct {
				Starter   bool   `json:"starter"`
				Jersey    string `json:"jersey"`
				SubbedIn  bool   `json:"subbedIn"`
				SubbedOut bool   `json:"subbedOut"`
				Athlete   struct {
					ID          string `json:"id"`
					DisplayName string `json:"displayName"`
				} `json:"athlete"`
				Position struct {
					Abbreviation string `json:"abbreviation"` // "CD-L"
				} `json:"position"`
				// The player's moments: goals, cards and substitutions
				Plays []struct {
					Clock struct {
						DisplayValue string `json:"displayValue"` // "83'", "90'+2'"
					} `json:"clock"`
					Substitution bool `json:"substitution"`
				} `json:"plays"`
			} `json:"roster"`
		} `json:"rosters"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("rosters: %w", err)
	}

	var out event.Lineups
	found := false
	for _, t := range raw.Rosters {
		var starters, bench []event.Player
		var starting []event.Spot
		for _, r := range t.Roster {
			p := event.Player{ID: r.Athlete.ID, Name: r.Athlete.DisplayName, Number: r.Jersey, Starter: r.Starter}
			// Came on at the first substitution, went off at the last: a
			// player can do both
			var subs []*event.Moment
			for _, play := range r.Plays {
				if m := moment(play.Clock.DisplayValue); play.Substitution && m != nil {
					subs = append(subs, m)
				}
			}
			if r.SubbedIn && len(subs) > 0 {
				p.On = subs[0]
			}
			if r.SubbedOut && len(subs) > 0 {
				p.Off = subs[len(subs)-1]
			}
			if r.Starter {
				spot, ok := spots[r.Position.Abbreviation]
				if !ok {
					spot = unknownSpot
				}
				starters = append(starters, p)
				starting = append(starting, spot)
			} else {
				bench = append(bench, p)
			}
		}
		if len(starters) == 0 {
			continue
		}
		found = true
		event.Arrange(t.Formation, starters, starting)
		event.SortBench(bench)
		l := event.Lineup{Formation: t.Formation, Players: append(starters, bench...)}
		if t.HomeAway == "home" {
			out.Home = l
		} else {
			out.Away = l
		}
	}
	if !found {
		return nil, nil
	}
	return &out, nil
}

// moment reads a minute as ESPN shows it, "83'" or "90'+2'", as the match
// time our clock shows the same minute at: 82:00 and 91:00. The period is
// the one the minute is in. nil if it does not parse.
func moment(clock string) *event.Moment {
	base, extra, _ := strings.Cut(clock, "+")
	minute, err := strconv.Atoi(strings.TrimSuffix(base, "'"))
	if err != nil || minute < 1 {
		return nil
	}
	stoppage := 0
	if extra != "" {
		if stoppage, err = strconv.Atoi(strings.TrimSuffix(extra, "'")); err != nil {
			return nil
		}
	}
	period := 4
	for p, end := range []int{45, 90, 105} {
		if minute <= end {
			period = p + 1
			break
		}
	}
	return &event.Moment{Period: period, Elapsed: time.Duration(minute+stoppage-1) * time.Minute}
}
