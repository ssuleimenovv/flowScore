package statsbomb

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata" // kick-off times are local; Windows has no time zone database for Go

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// teamCodes are the Premier League 2015/16 three-letter codes.
// StatsBomb gives team names only.
var teamCodes = map[string]string{
	"AFC Bournemouth":      "BOU",
	"Arsenal":              "ARS",
	"Aston Villa":          "AVL",
	"Chelsea":              "CHE",
	"Crystal Palace":       "CRY",
	"Everton":              "EVE",
	"Leicester City":       "LEI",
	"Liverpool":            "LIV",
	"Manchester City":      "MCI",
	"Manchester United":    "MUN",
	"Newcastle United":     "NEW",
	"Norwich City":         "NOR",
	"Southampton":          "SOU",
	"Stoke City":           "STK",
	"Sunderland":           "SUN",
	"Swansea City":         "SWA",
	"Tottenham Hotspur":    "TOT",
	"Watford":              "WAT",
	"West Bromwich Albion": "WBA",
	"West Ham United":      "WHU",
}

// LoadMatchInfo reads the fixture metadata: teams, competition, venue, kick-off.
func LoadMatchInfo(matchesPath string, matchID int) (event.Match, error) {
	m, err := findMatch(matchesPath, matchID)
	if err != nil {
		return event.Match{}, err
	}

	// StatsBomb stores local time; English matches are played on London time.
	london, err := time.LoadLocation("Europe/London")
	if err != nil {
		return event.Match{}, err
	}
	kickoff, err := time.ParseInLocation("2006-01-02 15:04:05.000", m.MatchDate+" "+m.KickOff, london)
	if err != nil {
		return event.Match{}, fmt.Errorf("kick-off of match %d: %w", matchID, err)
	}

	venue := ""
	if m.Stadium != nil {
		venue = m.Stadium.Name
	}

	return event.Match{
		ID:            strconv.Itoa(matchID),
		Source:        "replay", // a StatsBomb match is history, played again
		CompetitionID: strconv.Itoa(m.Competition.ID),
		Competition:   m.Competition.Name,
		Round:         strconv.Itoa(m.MatchWeek),
		Venue:         venue,
		KickoffAt:     kickoff.UTC(),
		Home:          team(m.HomeTeam.ID, m.HomeTeam.Name),
		Away:          team(m.AwayTeam.ID, m.AwayTeam.Name),
	}, nil
}

func team(id int, name string) event.Team {
	code, ok := teamCodes[name]
	if !ok {
		code = strings.ToUpper(string([]rune(name)[:min(3, len([]rune(name)))]))
	}
	return event.Team{ID: strconv.Itoa(id), Code: code, Name: name}
}
