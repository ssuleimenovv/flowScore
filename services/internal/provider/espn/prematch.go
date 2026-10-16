package espn

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// The meetings' dates have seconds, the other dates of a summary do not
const meetingLayout = "2006-01-02T15:04:05Z07:00"

// A competition is named with its season: "2025-26 English Premier League"
var seasonPrefix = regexp.MustCompile(`^\d{4}(-\d{2})?\s+`)

// ParsePreMatch reads what a summary says before kick-off: each team's last
// five matches, their last meetings and the league's news. A summary has
// them before, during and after the match alike.
func ParsePreMatch(body []byte) (*event.PreMatch, error) {
	var raw struct {
		Header struct {
			Competitions []struct {
				Competitors []struct {
					HomeAway string  `json:"homeAway"`
					Team     rawTeam `json:"team"`
				} `json:"competitors"`
			} `json:"competitions"`
		} `json:"header"`
		LastFiveGames []struct {
			Team   rawTeam `json:"team"`
			Events []struct {
				ID          string  `json:"id"`
				Date        string  `json:"gameDate"` // "2026-09-06T13:00Z"
				HomeTeamID  string  `json:"homeTeamId"`
				HomeScore   string  `json:"homeTeamScore"`
				AwayScore   string  `json:"awayTeamScore"`
				Competition string  `json:"leagueAbbreviation"` // "Premier League", "Carabao Cup"
				Opponent    rawTeam `json:"opponent"`
			} `json:"events"`
		} `json:"lastFiveGames"`
		SeasonSeries []struct {
			Type   string `json:"type"` // "head-to-head"
			Events []struct {
				ID          string `json:"id"`
				Date        string `json:"date"` // "2026-03-21T17:30:00Z"
				Status      string `json:"status"`
				Competition string `json:"competitionName"`
				Competitors []struct {
					HomeAway string  `json:"homeAway"`
					Score    string  `json:"score"`
					Team     rawTeam `json:"team"`
				} `json:"competitors"`
			} `json:"events"`
		} `json:"seasonseries"`
		News struct {
			Articles []struct {
				Headline  string `json:"headline"`
				Published string `json:"published"` // "2026-10-16T11:13:11Z"
				Links     struct {
					Web struct {
						Href string `json:"href"`
					} `json:"web"`
				} `json:"links"`
			} `json:"articles"`
		} `json:"news"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("pre-match: %w", err)
	}

	var homeID string
	if len(raw.Header.Competitions) > 0 {
		for _, c := range raw.Header.Competitions[0].Competitors {
			if c.HomeAway == "home" {
				homeID = c.Team.ID
			}
		}
	}

	var out event.PreMatch
	for _, t := range raw.LastFiveGames {
		var games []event.PastMatch
		for _, e := range t.Events {
			date, err := time.Parse(dateLayout, e.Date)
			if err != nil {
				return nil, fmt.Errorf("pre-match, match %s: %w", e.ID, err)
			}
			g := event.PastMatch{ID: e.ID, Date: date, Competition: e.Competition, Home: t.Team.team(), Away: e.Opponent.team()}
			g.HomeGoals, g.AwayGoals = goals(e.HomeScore), goals(e.AwayScore)
			if e.HomeTeamID != t.Team.ID {
				g.Home, g.Away = g.Away, g.Home
			}
			games = append(games, g)
		}
		if t.Team.ID == homeID {
			out.Home = games
		} else {
			out.Away = games
		}
	}

	for _, s := range raw.SeasonSeries {
		if s.Type != "head-to-head" {
			continue
		}
		for _, e := range s.Events {
			date, err := time.Parse(meetingLayout, e.Date)
			if err != nil || e.Status != "post" {
				continue // a meeting still to come, or one ESPN dates oddly
			}
			g := event.PastMatch{ID: e.ID, Date: date, Competition: seasonPrefix.ReplaceAllString(e.Competition, "")}
			for _, c := range e.Competitors {
				if c.HomeAway == "home" {
					g.Home, g.HomeGoals = c.Team.team(), goals(c.Score)
				} else {
					g.Away, g.AwayGoals = c.Team.team(), goals(c.Score)
				}
			}
			out.Meetings = append(out.Meetings, g)
		}
	}

	for _, a := range raw.News.Articles {
		published, err := time.Parse(time.RFC3339, a.Published)
		if err != nil || a.Headline == "" {
			continue
		}
		out.News = append(out.News, event.Article{Headline: a.Headline, URL: a.Links.Web.Href, Published: published})
	}
	// ESPN lists its stories in an order of its own: the newest go first here
	slices.SortStableFunc(out.News, func(a, b event.Article) int { return b.Published.Compare(a.Published) })

	if out.Home == nil && out.Away == nil && out.Meetings == nil && out.News == nil {
		return nil, nil
	}
	return &out, nil
}

type rawTeam struct {
	ID           string `json:"id"`
	DisplayName  string `json:"displayName"`
	Abbreviation string `json:"abbreviation"`
}

func (t rawTeam) team() event.Team {
	return event.Team{ID: t.ID, Code: t.Abbreviation, Name: t.DisplayName}
}

// goals reads a score of a match over; ESPN writes it as text
func goals(score string) int {
	n, _ := strconv.Atoi(score)
	return n
}
