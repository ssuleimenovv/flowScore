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

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

const (
	baseURL      = "https://site.api.espn.com/apis/site/v2/sports/soccer"
	standingsURL = "https://site.api.espn.com/apis/v2/sports/soccer" // the tables live elsewhere
)

// A summary is about 200 KB; a body far past that is not one
const maxBody = 10 << 20

// Client asks the API. BaseURL and StandingsURL are there for tests.
type Client struct {
	HTTP         *http.Client
	BaseURL      string
	StandingsURL string
}

func NewClient() *Client {
	return &Client{HTTP: &http.Client{Timeout: 15 * time.Second}, BaseURL: baseURL, StandingsURL: standingsURL}
}

// Game is one match of the scoreboard, with what picking it and listing it
// needs.
type Game struct {
	ID         string
	Name       string // "Leeds United at Arsenal"
	Start      time.Time
	State      string // "pre", "in" or "post"
	Clock      string // "72'" as ESPN shows it; while in play it is the minute
	Home, Away Side
}

// Side is one team of a game and its goals, "" before kick-off.
type Side struct {
	Team  event.Team
	Score string
}

// Board is a league's scoreboard of one day: its matches, and every day of
// the season that has any.
type Board struct {
	Games []Game
	Days  []time.Time
}

// Scoreboard returns today's matches of a league: "eng.1" is the Premier League.
func (c *Client) Scoreboard(ctx context.Context, league string) ([]Game, error) {
	b, err := c.Board(ctx, league, time.Time{})
	return b.Games, err
}

// Board returns the scoreboard of a league on a day, by its date in UTC.
// The zero day is today.
func (c *Client) Board(ctx context.Context, league string, day time.Time) (Board, error) {
	url := fmt.Sprintf("%s/%s/scoreboard", c.BaseURL, league)
	if !day.IsZero() {
		url += "?dates=" + day.UTC().Format("20060102")
	}
	body, err := c.get(ctx, url)
	if err != nil {
		return Board{}, err
	}
	var raw struct {
		Leagues []struct {
			Calendar []string `json:"calendar"` // "2026-08-21T07:00Z"
		} `json:"leagues"`
		Events []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Date   string `json:"date"` // "2026-10-10T11:30Z"
			Status struct {
				DisplayClock string `json:"displayClock"`
				Type         struct {
					State string `json:"state"`
				} `json:"type"`
			} `json:"status"`
			Competitions []struct {
				Competitors []struct {
					HomeAway string `json:"homeAway"`
					Score    string `json:"score"`
					Team     struct {
						ID           string `json:"id"`
						Abbreviation string `json:"abbreviation"`
						DisplayName  string `json:"displayName"`
					} `json:"team"`
				} `json:"competitors"`
			} `json:"competitions"`
		} `json:"events"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return Board{}, fmt.Errorf("scoreboard %s: %w", league, err)
	}

	var b Board
	if len(raw.Leagues) > 0 {
		for _, d := range raw.Leagues[0].Calendar {
			t, err := time.Parse(dateLayout, d)
			if err != nil {
				return Board{}, fmt.Errorf("scoreboard %s, calendar: %w", league, err)
			}
			b.Days = append(b.Days, t)
		}
	}
	for _, e := range raw.Events {
		start, err := time.Parse(dateLayout, e.Date)
		if err != nil {
			return Board{}, fmt.Errorf("scoreboard %s, match %s: %w", league, e.ID, err)
		}
		g := Game{ID: e.ID, Name: e.Name, Start: start, State: e.Status.Type.State, Clock: e.Status.DisplayClock}
		if len(e.Competitions) > 0 {
			for _, t := range e.Competitions[0].Competitors {
				side := Side{Team: event.Team{ID: t.Team.ID, Code: t.Team.Abbreviation, Name: t.Team.DisplayName}, Score: t.Score}
				if t.HomeAway == "home" {
					g.Home = side
				} else {
					g.Away = side
				}
			}
		}
		b.Games = append(b.Games, g)
	}
	return b, nil
}

// ESPN's dates have no seconds
const dateLayout = "2006-01-02T15:04Z07:00"

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
