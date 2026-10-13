package espn

import (
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// commentary is a few plays of Coventry City – Newcastle United, 12 October
// 2026, as ESPN sent them, cut down to what ParsePlays reads
const commentary = `{"commentary": [
	{"sequence": 0, "text": "Lineups are announced and players are warming up."},
	{"play": {"id": "53551559", "type": {"type": "kickoff"}, "period": {"number": 1}, "clock": {"value": 0, "displayValue": ""}}},
	{"play": {"id": "53551972", "type": {"type": "penalty---scored"}, "period": {"number": 1},
		"clock": {"value": 340, "displayValue": "6'"}, "team": {"displayName": "Coventry City"},
		"text": "Goal! Coventry City 1, Newcastle United 0. Jack Rudoni (Coventry City) converts the penalty with a left footed shot to the bottom right corner.",
		"participants": [{"athlete": {"displayName": "Jack Rudoni"}}], "fieldPositionX": 88.5, "fieldPositionY": 50.0}},
	{"play": {"id": "53558553", "type": {"type": "shot-off-target"}, "period": {"number": 2},
		"clock": {"value": 2988, "displayValue": "50'"}, "team": {"displayName": "Newcastle United"},
		"text": "Attempt missed. Joe Willock (Newcastle United) header from the left side of the six yard box is close, but misses the top left corner. Assisted by Jacob Murphy with a cross.",
		"participants": [{"athlete": {"displayName": "Joe Willock"}}, {"athlete": {"displayName": "Jacob Murphy"}}],
		"fieldPositionX": 97.8, "fieldPositionY": 57.7}},
	{"play": {"id": "53560148", "type": {"type": "yellow-card"}, "period": {"number": 2},
		"clock": {"value": 3902, "displayValue": "66'"}, "team": {"displayName": "Coventry City"},
		"participants": [{"athlete": {"displayName": "Milan van Ewijk"}}]}},
	{"play": {"id": "53563292", "type": {"type": "shot-blocked"}, "period": {"number": 2},
		"clock": {"value": 5400, "displayValue": "90'+2'"}, "team": {"displayName": "Newcastle United"},
		"participants": [{"athlete": {"displayName": "Amar Dedic"}}], "fieldPositionX": 75.8, "fieldPositionY": 38.8}}
]}`

func TestParsePlays(t *testing.T) {
	m := event.Match{ID: "401879267", Home: event.Team{Name: "Coventry City"}, Away: event.Team{Name: "Newcastle United"}}
	plays, err := ParsePlays([]byte(commentary), m)
	if err != nil {
		t.Fatal(err)
	}
	// Words without a play and the kick-off are not events of ours
	if len(plays) != 4 {
		t.Fatalf("%d plays: %+v", len(plays), plays)
	}

	pen := plays[0]
	if pen.ID != "401879267:53551972" || pen.Type != event.Goal || pen.Side != event.Home ||
		pen.Player != "Jack Rudoni" || pen.Situation != "penalty" || pen.Elapsed != 340*time.Second ||
		pen.Pos == nil || *pen.Pos != (event.Position{X: 88.5, Y: 50}) {
		t.Errorf("penalty: %+v", pen)
	}
	if h := plays[1]; h.Type != event.ShotOffTarget || h.Side != event.Away || !h.Header ||
		h.Situation != "open_play" || h.Player != "Joe Willock" || h.Period != 2 {
		t.Errorf("header: %+v", h)
	}
	if c := plays[2]; c.Type != event.YellowCard || c.Side != event.Home || c.Pos != nil || c.Situation != "" {
		t.Errorf("card: %+v", c)
	}
	// ESPN's clock stands at 90:00; the two added minutes are only in the text
	if b := plays[3]; b.Type != event.ShotBlocked || b.Elapsed != 92*time.Minute {
		t.Errorf("stoppage time: %+v", b)
	}
}
