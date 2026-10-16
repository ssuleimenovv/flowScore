package event

import "time"

// PreMatch is what is known of a match before it starts: how both teams
// have played lately, how their meetings went, and the news of the league.
type PreMatch struct {
	Home, Away []PastMatch // each team's last matches, oldest first
	Meetings   []PastMatch // the two teams' last meetings, newest first
	News       []Article   // newest first
}

// PastMatch is a match over, of any competition.
type PastMatch struct {
	ID                   string
	Date                 time.Time
	Competition          string // "Premier League", "Carabao Cup"
	Home, Away           Team
	HomeGoals, AwayGoals int
}

type Article struct {
	Headline  string
	URL       string
	Published time.Time
}
