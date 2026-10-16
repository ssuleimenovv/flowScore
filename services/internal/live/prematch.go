package live

import (
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Flows is where the Flow of a match over comes from: the season's run of
// every played match. ok is false for a match it has not read.
type Flows interface {
	PlayedFlow(matchID string) (home, away float64, ok bool)
}

// MatchPreMatch is "Перед матчем" (MatchPreMatch in the contract).
type MatchPreMatch struct {
	MatchID  string   `json:"matchId"`
	Home     TeamForm `json:"home"`
	Away     TeamForm `json:"away"`
	Meetings Meetings `json:"meetings"`
	News     []News   `json:"news"`
}

// TeamForm is a team's last matches and its Flow in the ones the season read.
type TeamForm struct {
	Matches   []FormMatch `json:"matches"`   // oldest first
	Flow      *float64    `json:"flow"`      // averaged over flowGames; null if none was read
	FlowGames int         `json:"flowGames"` // a cup match has no Flow of ours
}

type FormMatch struct {
	PastMatch
	Result string `json:"result"` // win, draw or loss, for the team
}

type PastMatch struct {
	ID          string      `json:"id"`
	Date        time.Time   `json:"date"`
	Competition string      `json:"competition"`
	Home        team        `json:"home"`
	Away        team        `json:"away"`
	Score       Score       `json:"score"`
	Flow        *FlowValues `json:"flow"` // averaged over the match; null if the season did not read it
}

// Meetings are the two teams' last meetings and how they went, for this
// match's home team and away team.
type Meetings struct {
	HomeWins int         `json:"homeWins"`
	Draws    int         `json:"draws"`
	AwayWins int         `json:"awayWins"`
	Matches  []PastMatch `json:"matches"` // newest first
}

type News struct {
	Headline    string    `json:"headline"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"publishedAt"`
}

type team struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// PreMatchOf is the snapshot's pre-match info, with the Flow of the matches
// the season read; false if the source has none, as for a demo match.
func PreMatchOf(s Snapshot, flows Flows) (MatchPreMatch, bool) {
	pm := s.PreMatch
	if pm == nil {
		return MatchPreMatch{}, false
	}
	past := func(m event.PastMatch) PastMatch {
		p := PastMatch{
			ID: m.ID, Date: m.Date, Competition: m.Competition,
			Home:  team{ID: m.Home.ID, Code: m.Home.Code, Name: m.Home.Name},
			Away:  team{ID: m.Away.ID, Code: m.Away.Code, Name: m.Away.Name},
			Score: Score{Home: m.HomeGoals, Away: m.AwayGoals},
		}
		if flows != nil {
			if home, away, ok := flows.PlayedFlow(m.ID); ok {
				p.Flow = &FlowValues{Home: home, Away: away}
			}
		}
		return p
	}
	form := func(teamID string, matches []event.PastMatch) TeamForm {
		f := TeamForm{Matches: []FormMatch{}}
		sum := 0.0
		for _, m := range matches {
			p := past(m)
			away := m.Away.ID == teamID
			mine, theirs := m.HomeGoals, m.AwayGoals
			if away {
				mine, theirs = theirs, mine
			}
			f.Matches = append(f.Matches, FormMatch{PastMatch: p, Result: result(mine, theirs)})
			if p.Flow != nil {
				if away {
					sum += p.Flow.Away
				} else {
					sum += p.Flow.Home
				}
				f.FlowGames++
			}
		}
		if f.FlowGames > 0 {
			avg := sum / float64(f.FlowGames)
			f.Flow = &avg
		}
		return f
	}

	out := MatchPreMatch{
		MatchID:  s.Match.ID,
		Home:     form(s.Match.Home.ID, pm.Home),
		Away:     form(s.Match.Away.ID, pm.Away),
		Meetings: Meetings{Matches: []PastMatch{}},
		News:     []News{},
	}
	for _, m := range pm.Meetings {
		out.Meetings.Matches = append(out.Meetings.Matches, past(m))
		goals := map[string]int{m.Home.ID: m.HomeGoals, m.Away.ID: m.AwayGoals}
		switch home, away := goals[s.Match.Home.ID], goals[s.Match.Away.ID]; {
		case home > away:
			out.Meetings.HomeWins++
		case home < away:
			out.Meetings.AwayWins++
		default:
			out.Meetings.Draws++
		}
	}
	for _, a := range pm.News {
		out.News = append(out.News, News{Headline: a.Headline, URL: a.URL, PublishedAt: a.Published})
	}
	return out, true
}

func result(mine, theirs int) string {
	switch {
	case mine > theirs:
		return "win"
	case mine < theirs:
		return "loss"
	default:
		return "draw"
	}
}
