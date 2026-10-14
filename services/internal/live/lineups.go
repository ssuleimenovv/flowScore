package live

import (
	"cmp"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// MatchLineups are both teams' players (MatchLineups in the contract).
type MatchLineups struct {
	MatchID string     `json:"matchId"`
	Home    TeamLineup `json:"home"`
	Away    TeamLineup `json:"away"`
}

type TeamLineup struct {
	Formation *string        `json:"formation"` // "4-2-3-1"; null when the source does not say
	Starters  []LineupPlayer `json:"starters"`  // from the keeper forward, each line from left to right
	Bench     []LineupPlayer `json:"bench"`     // the ones who came on first
}

type LineupPlayer struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Number   *string         `json:"number"`
	Position *event.Position `json:"position"` // a starter's place, attacking left to right; null on the bench
	On       *Minute         `json:"on"`       // when the player came on; null if not yet
	Off      *Minute         `json:"off"`      // when the player went off; null if not yet
}

// Minute is a minute on the clock: 90+3 is 90 with 3 added.
type Minute struct {
	Minute    int  `json:"minute"`
	AddedTime *int `json:"addedTime"`
}

// LineupsOf is the snapshot's lineups as far as the match has come: a demo
// match knows its substitutions from the start, and a viewer at the 20th
// minute must not see the one of the 70th. False before the source has them.
func LineupsOf(s Snapshot) (MatchLineups, bool) {
	if s.Lineups == nil {
		return MatchLineups{}, false
	}
	shown := func(m *event.Moment) *Minute {
		if m == nil || !happened(s, *m) {
			return nil
		}
		minute, added := minuteOf(m.Period, m.Elapsed)
		return &Minute{Minute: minute, AddedTime: added}
	}
	team := func(l event.Lineup) TeamLineup {
		t := TeamLineup{Starters: []LineupPlayer{}, Bench: []LineupPlayer{}}
		if l.Formation != "" {
			t.Formation = &l.Formation
		}
		for _, p := range l.Players {
			lp := LineupPlayer{ID: p.ID, Name: p.Name, Position: p.Pos, On: shown(p.On), Off: shown(p.Off)}
			if p.Number != "" {
				lp.Number = &p.Number
			}
			if p.Starter {
				t.Starters = append(t.Starters, lp)
			} else {
				t.Bench = append(t.Bench, lp)
			}
		}
		return t
	}
	return MatchLineups{MatchID: s.Match.ID, Home: team(s.Lineups.Home), Away: team(s.Lineups.Away)}, true
}

// happened reports whether the moment is behind the match: all of a
// finished one, none of one yet to start.
func happened(s Snapshot, m event.Moment) bool {
	switch s.Status {
	case "finished":
		return true
	case "scheduled":
		return false
	}
	return cmp.Or(cmp.Compare(m.Period, s.Period), cmp.Compare(m.Elapsed, s.At)) <= 0
}
