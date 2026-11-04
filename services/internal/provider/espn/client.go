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
	Market     *Market // nil when ESPN shows no odds
	Stage      string  // a cup's stage: "league-phase", "round-of-16", "final"; a league's season
	Leg        int     // 1 or 2 in a tie of two matches; 0 for a single match
	Finish     string  // "aet" after extra time, "pens" after a shootout; "" otherwise
}

// Side is one team of a game and its goals, "" before kick-off.
type Side struct {
	Team      event.Team
	Score     string
	Winner    bool // won the match, on penalties too
	Shootout  *int // its goals in a shootout; nil without one
	Aggregate *int // its goals over both legs of a tie; nil for a single match
	Through   bool // went through a tie of two legs that is over
}

// Board is a league's scoreboard of one day: its matches, and every day of
// the season that has any.
type Board struct {
	Games       []Game
	Days        []time.Time // none for a cup
	SeasonStart time.Time   // zero when ESPN does not say
	SeasonYear  int         // 2026 for 2026/27; 0 when ESPN does not say
}

// Scoreboard returns today's matches of a league: "eng.1" is the Premier League.
func (c *Client) Scoreboard(ctx context.Context, league string) ([]Game, error) {
	b, err := c.Board(ctx, league, time.Time{})
	return b.Games, err
}

// Board returns the scoreboard of a league on a day, by its date in UTC.
// The zero day is today.
func (c *Client) Board(ctx context.Context, league string, day time.Time) (Board, error) {
	query := ""
	if !day.IsZero() {
		query = "?dates=" + day.UTC().Format("20060102")
	}
	return c.board(ctx, league, query)
}

// Year returns every match of a league in a calendar year. A cup lists no
// days to read one by one, so its season is read this way.
func (c *Client) Year(ctx context.Context, league string, year int) (Board, error) {
	return c.board(ctx, league, fmt.Sprintf("?dates=%d&limit=1000", year))
}

func (c *Client) board(ctx context.Context, league, query string) (Board, error) {
	body, err := c.get(ctx, fmt.Sprintf("%s/%s/scoreboard%s", c.BaseURL, league, query))
	if err != nil {
		return Board{}, err
	}
	var raw struct {
		Leagues []struct {
			// A league's match days, "2026-08-21T07:00Z"; a cup's stages
			// instead ({label: "League Phase", ...}), so it is read below
			Calendar json.RawMessage `json:"calendar"`
			Season   struct {
				Year      int    `json:"year"`
				StartDate string `json:"startDate"` // "2026-07-01T04:00Z"
			} `json:"season"`
		} `json:"leagues"`
		Events []struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Date   string `json:"date"` // "2026-10-10T11:30Z"
			Season struct {
				Slug string `json:"slug"`
			} `json:"season"`
			Status struct {
				DisplayClock string `json:"displayClock"`
				Type         struct {
					State string `json:"state"`
					Name  string `json:"name"` // "STATUS_FINAL_PEN"
				} `json:"type"`
			} `json:"status"`
			Competitions []struct {
				Leg struct {
					Value int `json:"value"`
				} `json:"leg"`
				Series struct {
					Competitors []struct {
						ID             string   `json:"id"`
						Winner         bool     `json:"winner"`
						AggregateScore *float64 `json:"aggregateScore"` // 7.0
					} `json:"competitors"`
				} `json:"series"`
				Competitors []struct {
					HomeAway      string   `json:"homeAway"`
					Score         string   `json:"score"`
					Winner        bool     `json:"winner"`
					ShootoutScore *float64 `json:"shootoutScore"`
					Team          struct {
						ID           string `json:"id"`
						Abbreviation string `json:"abbreviation"`
						DisplayName  string `json:"displayName"`
					} `json:"team"`
				} `json:"competitors"`
				Odds []struct {
					Moneyline struct {
						Home, Draw, Away moneyline
					} `json:"moneyline"`
				} `json:"odds"`
			} `json:"competitions"`
		} `json:"events"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return Board{}, fmt.Errorf("scoreboard %s: %w", league, err)
	}

	var b Board
	if len(raw.Leagues) > 0 {
		b.SeasonYear = raw.Leagues[0].Season.Year
		if s := raw.Leagues[0].Season.StartDate; s != "" {
			t, err := time.Parse(dateLayout, s)
			if err != nil {
				return Board{}, fmt.Errorf("scoreboard %s, season: %w", league, err)
			}
			b.SeasonStart = t
		}
		// A cup's calendar is not a list of days: it has none to give
		var days []string
		if json.Unmarshal(raw.Leagues[0].Calendar, &days) == nil {
			for _, d := range days {
				t, err := time.Parse(dateLayout, d)
				if err != nil {
					return Board{}, fmt.Errorf("scoreboard %s, calendar: %w", league, err)
				}
				b.Days = append(b.Days, t)
			}
		}
	}

	for _, e := range raw.Events {
		start, err := time.Parse(dateLayout, e.Date)
		if err != nil {
			return Board{}, fmt.Errorf("scoreboard %s, match %s: %w", league, e.ID, err)
		}
		g := Game{ID: e.ID, Name: e.Name, Start: start, State: e.Status.Type.State, Clock: e.Status.DisplayClock,
			Stage: e.Season.Slug, Finish: finish(e.Status.Type.Name)}
		if len(e.Competitions) > 0 {
			c := e.Competitions[0]
			g.Leg = c.Leg.Value
			if len(c.Odds) > 0 {
				ml := c.Odds[0].Moneyline
				g.Market = market(ml.Home, ml.Draw, ml.Away)
			}
			for _, t := range c.Competitors {
				side := Side{Team: event.Team{ID: t.Team.ID, Code: t.Team.Abbreviation, Name: t.Team.DisplayName}, Score: t.Score,
					Winner: t.Winner, Shootout: whole(t.ShootoutScore)}
				for _, s := range c.Series.Competitors {
					if s.ID == t.Team.ID {
						side.Aggregate, side.Through = whole(s.AggregateScore), s.Winner
					}
				}
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

// finish names how a match over was decided past the 90 minutes
func finish(status string) string {
	switch status {
	case "STATUS_FINAL_AET":
		return "aet"
	case "STATUS_FINAL_PEN":
		return "pens"
	}
	return ""
}

// whole turns ESPN's 7.0 into 7; nil stays nil
func whole(f *float64) *int {
	if f == nil {
		return nil
	}
	n := int(*f)
	return &n
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
