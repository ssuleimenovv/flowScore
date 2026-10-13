package espn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Table is a league's table as it stands
type Table struct {
	Season int // 2026 for 2026/27
	Rows   []Standing
}

// Standing is one team's row of the table
type Standing struct {
	Rank                     int
	Team                     event.Team
	Played, Won, Drawn, Lost int
	For, Against             int // goals
	Points                   int
	Zone                     string // what the place leads to: "Champions League", "Relegation"; "" for nothing
}

// Standings returns a league's table
func (c *Client) Standings(ctx context.Context, league string) (Table, error) {
	body, err := c.get(ctx, fmt.Sprintf("%s/%s/standings", c.StandingsURL, league))
	if err != nil {
		return Table{}, err
	}
	var raw struct {
		Season struct {
			Year int `json:"year"`
		} `json:"season"`
		Children []struct {
			Standings struct {
				Entries []struct {
					Team struct {
						ID           string `json:"id"`
						Abbreviation string `json:"abbreviation"`
						DisplayName  string `json:"displayName"`
					} `json:"team"`
					Note *struct {
						Description string `json:"description"`
					} `json:"note"`
					Stats []struct {
						Name  string   `json:"name"`
						Value *float64 `json:"value"` // null for the "overall" record
					} `json:"stats"`
				} `json:"entries"`
			} `json:"standings"`
		} `json:"children"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return Table{}, fmt.Errorf("standings %s: %w", league, err)
	}
	// A league of one table has one child; groups are for cups, not read yet
	if len(raw.Children) == 0 {
		return Table{}, fmt.Errorf("standing %s: %w", league, errors.New("no table"))
	}

	t := Table{Season: raw.Season.Year}
	for _, e := range raw.Children[0].Standings.Entries {
		stat := map[string]int{}
		for _, s := range e.Stats {
			if s.Value != nil {
				stat[s.Name] = int(*s.Value)
			}
		}
		row := Standing{
			Rank:    stat["rank"],
			Team:    event.Team{ID: e.Team.ID, Code: e.Team.Abbreviation, Name: e.Team.DisplayName},
			Played:  stat["gamesPlayed"],
			Won:     stat["wins"],
			Drawn:   stat["ties"],
			Lost:    stat["losses"],
			For:     stat["pointsFor"],
			Against: stat["pointsAgainst"],
			Points:  stat["points"],
		}
		if e.Note != nil {
			row.Zone = e.Note.Description
		}
		t.Rows = append(t.Rows, row)
	}
	return t, nil
}
