package statsbomb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

type rawLineup struct {
	TeamID int `json:"team_id"`
	Lineup []struct {
		ID        int     `json:"player_id"`
		Name      string  `json:"player_name"`
		Nickname  *string `json:"player_nickname"`
		Number    int     `json:"jersey_number"`
		Positions []struct {
			ID          int     `json:"position_id"`
			From        string  `json:"from"` // "78:20"
			To          *string `json:"to"`
			FromPeriod  int     `json:"from_period"`
			ToPeriod    *int    `json:"to_period"`
			StartReason string  `json:"start_reason"` // "Starting XI", "Substitution - On (Tactical)"
			EndReason   string  `json:"end_reason"`   // "Final Whistle", "Substitution - Off (Injury)"
		} `json:"positions"`
	} `json:"lineup"`
}

// loadNames returns the name to show for every player of the match. StatsBomb
// names players in full ("Sergio Leonel Agüero del Castillo"); the nickname is
// the one fans know ("Sergio Agüero"), when there is one.
func loadNames(path string) (map[int]string, error) {
	var teams []rawLineup
	if err := readJSON(path, &teams); err != nil {
		return nil, fmt.Errorf("read lineups: %w", err)
	}

	names := map[int]string{}
	for _, team := range teams {
		for _, p := range team.Lineup {
			names[p.ID] = shownName(p.Name, p.Nickname)
		}
	}
	return names, nil
}

func shownName(name string, nickname *string) string {
	if nickname != nil && *nickname != "" {
		return *nickname
	}
	return name
}

// spots are StatsBomb's positions on its grid of 25 (its Open Data spec):
// 1 is the keeper, 2–8 the defenders and wing-backs, 9–11 the defensive
// midfielders, 12–16 the midfielders, 17–21 the attacking midfielders and
// wingers, 22–25 the forwards. Their flanks run from left to right.
var spots = map[int]event.Spot{
	1: {Line: 0, Flank: 0},
	2: {Line: 1, Flank: 2}, 3: {Line: 1, Flank: 1}, 4: {Line: 1, Flank: 0}, 5: {Line: 1, Flank: -1}, 6: {Line: 1, Flank: -2},
	7: {Line: 1.5, Flank: 2}, 8: {Line: 1.5, Flank: -2},
	9: {Line: 2, Flank: 1}, 10: {Line: 2, Flank: 0}, 11: {Line: 2, Flank: -1},
	12: {Line: 2.5, Flank: 2}, 13: {Line: 2.5, Flank: 1}, 14: {Line: 2.5, Flank: 0}, 15: {Line: 2.5, Flank: -1}, 16: {Line: 2.5, Flank: -2},
	17: {Line: 3.5, Flank: 2}, 18: {Line: 3, Flank: 1}, 19: {Line: 3, Flank: 0}, 20: {Line: 3, Flank: -1}, 21: {Line: 3.5, Flank: -2},
	22: {Line: 4, Flank: 1}, 23: {Line: 4, Flank: 0}, 24: {Line: 4, Flank: -1}, 25: {Line: 3.8, Flank: 0},
}

// LoadLineups reads both teams' players from the lineups file and their
// formations from the "Starting XI" events of the match.
func LoadLineups(lineupsPath, eventsPath string, homeTeamID int) (*event.Lineups, error) {
	var teams []rawLineup
	if err := readJSON(lineupsPath, &teams); err != nil {
		return nil, fmt.Errorf("read lineups: %w", err)
	}
	var raw []struct {
		Type    ref `json:"type"`
		Team    ref `json:"team"`
		Tactics *struct {
			Formation int `json:"formation"` // 4231
		} `json:"tactics"`
	}
	if err := readJSON(eventsPath, &raw); err != nil {
		return nil, fmt.Errorf("read events: %w", err)
	}
	formations := map[int]string{}
	for _, r := range raw {
		if r.Type.Name == "Starting XI" && r.Tactics != nil {
			formations[r.Team.ID] = formation(r.Tactics.Formation)
		}
	}

	var out event.Lineups
	for _, t := range teams {
		l := lineup(t, formations[t.TeamID])
		if t.TeamID == homeTeamID {
			out.Home = l
		} else {
			out.Away = l
		}
	}
	return &out, nil
}

// formation writes 4231 as "4-2-3-1"
func formation(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Join(strings.Split(strconv.Itoa(n), ""), "-")
}

func lineup(t rawLineup, formation string) event.Lineup {
	var starters, bench []event.Player
	var starting []event.Spot
	for _, p := range t.Lineup {
		player := event.Player{ID: strconv.Itoa(p.ID), Name: shownName(p.Name, p.Nickname)}
		if p.Number > 0 {
			player.Number = strconv.Itoa(p.Number)
		}
		for _, pos := range p.Positions {
			if strings.HasPrefix(pos.StartReason, "Substitution - On") && player.On == nil {
				player.On = moment(pos.FromPeriod, pos.From)
			}
			if strings.HasPrefix(pos.EndReason, "Substitution - Off") && pos.To != nil && pos.ToPeriod != nil {
				player.Off = moment(*pos.ToPeriod, *pos.To)
			}
		}
		if len(p.Positions) > 0 && p.Positions[0].StartReason == "Starting XI" {
			player.Starter = true
			starters = append(starters, player)
			starting = append(starting, spots[p.Positions[0].ID])
		} else {
			bench = append(bench, player)
		}
	}
	event.Arrange(formation, starters, starting)

	event.SortBench(bench)
	return event.Lineup{Formation: formation, Players: append(starters, bench...)}
}

// moment reads StatsBomb's "78:20" of a period; nil if it does not parse
func moment(period int, clock string) *event.Moment {
	m, s, ok := strings.Cut(clock, ":")
	minutes, err1 := strconv.Atoi(m)
	seconds, err2 := strconv.Atoi(s)
	if !ok || err1 != nil || err2 != nil {
		return nil
	}
	return &event.Moment{Period: period, Elapsed: time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second}
}
