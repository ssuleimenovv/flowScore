package espn

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Summary is what one summary says about the match as a whole: who plays,
// where the clock is and the score. The plays in it are read elsewhere.
type Summary struct {
	Match  event.Match
	State  string        // "pre", "in" or "post"
	Status string        // "STATUS_FIRST_HALF", "STATUS_HALFTIME", "STATUS_FULL_TIME"...
	Period int           // 1 or 2; 0 before kick-off
	Clock  time.Duration // match time: 46:10 in the second half is 46m10s
	Score  Score
	// Possession is the home team's share of the ball over the match so fat
	// 0-1; nil until ESPN counts it
	Possession *float64
}

type Score struct {
	Home, Away int
}

// ParseSummary reads the header of a summary.
func ParseSummary(body []byte) (Summary, error) {
	var raw struct {
		Header struct {
			ID     string `json:"id"`
			League struct {
				ID        string `json:"id"`
				ShortName string `json:"shortName"` // "Premier League"
			} `json:"league"`
			Competitions []struct {
				Date   string `json:"date"`
				Status struct {
					Clock  float64 `json:"clock"` // seconds
					Period int     `json:"period"`
					Type   struct {
						Name  string `json:"name"`
						State string `json:"state"`
					} `json:"type"`
				} `json:"status"`
				Competitors []struct {
					HomeAway string `json:"homeAway"`
					Score    string `json:"score"` // "2"; missing before kick-off
					Team     struct {
						ID           string `json:"id"`
						Abbreviation string `json:"abbreviation"`
						DisplayName  string `json:"displayName"`
					} `json:"team"`
				} `json:"competitors"`
			} `json:"competitions"`
		} `json:"header"`

		Boxscore struct {
			Teams []struct {
				HomeAway   string `json:"homeAway"`
				Statistics []struct {
					Name         string `json:"name"`
					DisplayValue string `json:"displayValue"` // "33.3"
				} `json:"statistics"`
			} `json:"teams"`
		} `json:"boxscore"`

		GameInfo struct {
			Venue struct {
				FullName string `json:"fullName"`
			} `json:"venue"`
		} `json:"gameInfo"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return Summary{}, err
	}
	h := raw.Header
	if len(h.Competitions) == 0 {
		return Summary{}, errors.New("summary without a competition")
	}
	c := h.Competitions[0]
	start, err := time.Parse("2006-01-02T15:04Z07:00", c.Date)
	if err != nil {
		return Summary{}, fmt.Errorf("match %s: %w", h.ID, err)
	}

	s := Summary{
		Match: event.Match{
			ID:            h.ID,
			Source:        "espn",
			CompetitionID: h.League.ID,
			Competition:   h.League.ShortName,
			Venue:         raw.GameInfo.Venue.FullName,
			KickoffAt:     start,
		},
		State:  c.Status.Type.State,
		Status: c.Status.Type.Name,
		Period: c.Status.Period,
		Clock:  time.Duration(c.Status.Clock * float64(time.Second)),
	}
	for _, t := range c.Competitors {
		team := event.Team{ID: t.Team.ID, Code: t.Team.Abbreviation, Name: t.Team.DisplayName}
		goals, _ := strconv.Atoi(t.Score) // "" before kick-off is 0
		if t.HomeAway == "home" {
			s.Match.Home, s.Score.Home = team, goals
		} else {
			s.Match.Away, s.Score.Away = team, goals
		}
	}

	// Before kick-off there is no share, or one of 0 for both teams
	var home, away float64
	for _, t := range raw.Boxscore.Teams {
		for _, st := range t.Statistics {
			if st.Name != "possessionPct" {
				continue
			}
			pct, _ := strconv.ParseFloat(st.DisplayValue, 64)
			if t.HomeAway == "home" {
				home = pct
			} else {
				away = pct
			}
		}
	}
	if home+away > 0 {
		share := home / (home + away)
		s.Possession = &share
	}
	return s, nil
}
