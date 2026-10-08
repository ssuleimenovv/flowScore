package event

import "time"

// Match is what is known about a fixture apart from its events
type Match struct {
	ID            string
	Source        string // "espn" for a real match, "replay" for a demo one
	CompetitionID string
	Competition   string
	Round         string
	Venue         string
	KickoffAt     time.Time
	Home, Away    Team
}

type Team struct {
	ID   string
	Code string // three letters, e.g. MCI
	Name string
}
