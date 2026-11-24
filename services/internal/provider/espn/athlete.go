package espn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Athlete is a player as ESPN knows them: their club, number and line
type Athlete struct {
	ID       string
	Name     string
	Jersey   string // "" when ESPN gives none
	Position string // goalkeeper, defender, midfielder or forward; "" for none
	Age      int
	Country  string
	Team     *event.Team // their club; nil for a player without one
}

// Athlete returns a player by their ID: "253989" is Erling Haaland. ESPN
// finds a player under any league, so it is asked under "all".
// ErrNotFound for an ID it does not know.
func (c *Client) Athlete(ctx context.Context, id string) (Athlete, error) {
	body, err := c.get(ctx, fmt.Sprintf("%s/all/athletes/%s", c.AthleteURL, url.PathEscape(id)))
	if err != nil {
		return Athlete{}, err
	}
	return ParseAthlete(body)
}

// ParseAthlete reads a player's page. ESPN leaves the injuries out and has
// photos of few players, so neither is read, as in a squad.
func ParseAthlete(body []byte) (Athlete, error) {
	var raw struct {
		Athlete struct {
			ID          string `json:"id"`
			DisplayName string `json:"displayName"`
			Jersey      string `json:"jersey"`
			Age         int    `json:"age"`
			Citizenship string `json:"citizenship"`
			Position    struct {
				Name string `json:"name"` // "Forward"
			} `json:"position"`
			Team *struct {
				ID           string `json:"id"`
				Abbreviation string `json:"abbreviation"`
				DisplayName  string `json:"displayName"`
			} `json:"team"`
		} `json:"athlete"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return Athlete{}, fmt.Errorf("athlete: %w", err)
	}
	r := raw.Athlete
	if r.ID == "" {
		return Athlete{}, fmt.Errorf("athlete: %w", ErrNotFound)
	}
	a := Athlete{
		ID: r.ID, Name: r.DisplayName, Jersey: r.Jersey, Position: positions[r.Position.Name],
		Age: r.Age, Country: r.Citizenship,
	}
	if r.Team != nil {
		a.Team = &event.Team{ID: r.Team.ID, Code: r.Team.Abbreviation, Name: r.Team.DisplayName}
	}
	return a, nil
}
