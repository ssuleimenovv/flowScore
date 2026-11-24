package espn

import (
	"context"
	"encoding/json"
	"fmt"
)

// Player is one player of a team's squad for the season
type Player struct {
	ID       string
	Name     string // "Rúben Dias"
	Jersey   string // "" when ESPN gives none
	Position string // goalkeeper, defender, midfielder or forward; "" for none
	Age      int
	Country  string // "Portugal"
}

// positions are the four of a squad: ESPN names only the line there
var positions = map[string]string{
	"Goalkeeper": "goalkeeper",
	"Defender":   "defender",
	"Midfielder": "midfielder",
	"Forward":    "forward",
}

// Squad returns a team's players of the season: "382" is Manchester City.
// ESPN finds the team under any league it plays in.
func (c *Client) Squad(ctx context.Context, league, team string) ([]Player, error) {
	body, err := c.get(ctx, fmt.Sprintf("%s/%s/teams/%s/roster", c.BaseURL, league, team))
	if err != nil {
		return nil, err
	}
	return ParseSquad(body)
}

// ParseSquad reads the players of a team's roster. ESPN leaves the injuries
// empty there and has photos of few players, so neither is read.
func ParseSquad(body []byte) ([]Player, error) {
	var raw struct {
		Athletes []struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
			Jersey      string `json:"jersey"`
			Age         int    `json:"age"`
			Citizenship string `json:"citizenship"`
			Position    struct {
				Name string `json:"name"` // "Defender"
			} `json:"position"`
		} `json:"athletes"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("squad: %w", err)
	}
	players := make([]Player, 0, len(raw.Athletes))
	for _, a := range raw.Athletes {
		players = append(players, Player{
			ID: a.ID, Name: a.DisplayName, Jersey: a.Jersey,
			Position: positions[a.Position.Name], Age: a.Age, Country: a.Citizenship,
		})
	}
	return players, nil
}
