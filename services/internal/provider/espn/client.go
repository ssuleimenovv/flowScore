// Package espn reads live football from ESPN's public site API
// (docs/ARCHITECTURE.md): the scoreboard of a league and the summary of a
// match, with its commentary and key events. The API is unofficial and has
// no key; its format may change, so only this package knows it.
package espn

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://site.api.espn.com/apis/site/v2/sports/soccer"

// A summary is about 200 KB; a body far past that is not one
const maxBody = 10 << 20

// Client asks the API. BaseURL is there for tests.
type Client struct {
	HTTP    *http.Client
	BaseURL string
}

func NewClient() *Client {
	return &Client{HTTP: &http.Client{Timeout: 15 * time.Second}, BaseURL: baseURL}
}

// Game is one match of the scoreboard, with what picking it needs.
type Game struct {
	ID    string
	Name  string // "Leeds United at Arsenal"
	Start time.Time
	State string // "pre", "in" or "post"
}

// Scoreboard returns today's matches of a league: "eng.1" is the Premier League.
func (c *Client) Scoreboard(ctx context.Context, league string) ([]Game, error) {
	body, err := c.get(ctx, fmt.Sprintf("%s/%s/scoreboard", c.BaseURL, league))
	if err != nil {
		return nil, err
	}
	var board struct {
		Events []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Date   string `json:"date"` // "2026-10-10T11:30Z", without seconds
			Status struct {
				Type struct {
					State string `json:"state"`
				} `json:"type"`
			} `json:"status"`
		} `json:"events"`
	}
	if err := json.Unmarshal(body, &board); err != nil {
		return nil, fmt.Errorf("scoreboard %s: %w", league, err)
	}

	games := make([]Game, 0, len(board.Events))
	for _, e := range board.Events {
		start, err := time.Parse("2006-01-02T15:04Z07:00", e.Date)
		if err != nil {
			return nil, fmt.Errorf("scoreboard %s, match %s: %w", league, e.ID, err)
		}
		games = append(games, Game{ID: e.ID, Name: e.Name, Start: start, State: e.Status.Type.State})
	}
	return games, nil
}

// Summary returns the whole summary of a match as the API sends it.
func (c *Client) Summary(ctx context.Context, league, id string) ([]byte, error) {
	return c.get(ctx, fmt.Sprintf("%s/%s/summary?event=%s", c.BaseURL, league, id))
}

func (c *Client) get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: %s", url, res.Status)
	}
	return io.ReadAll(io.LimitReader(res.Body, maxBody))
}
