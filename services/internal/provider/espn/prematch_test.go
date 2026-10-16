package espn

import (
	"os"
	"testing"
	"time"
)

func TestParsePreMatch(t *testing.T) {
	// Everton – Chelsea, 17 October 2026, a day before kick-off
	body, err := os.ReadFile("testdata/prematch-401878771.json")
	if err != nil {
		t.Fatal(err)
	}
	p, err := ParsePreMatch(body)
	if err != nil || p == nil {
		t.Fatalf("pre-match %v, %v", p, err)
	}

	// Everton's last five, oldest first: 2-2 at home to United, 0-0 at Spurs...
	if len(p.Home) != 5 || len(p.Away) != 5 {
		t.Fatalf("%d and %d matches", len(p.Home), len(p.Away))
	}
	first, spurs, cup := p.Home[0], p.Home[1], p.Home[2]
	if first.ID != "401879291" || first.Home.Code != "EVE" || first.Away.Name != "Manchester United" ||
		first.HomeGoals != 2 || first.AwayGoals != 2 || !first.Date.Equal(time.Date(2026, 9, 6, 13, 0, 0, 0, time.UTC)) {
		t.Errorf("first %+v", first)
	}
	// Away at Tottenham: the opponent is the home team
	if spurs.Home.Code != "TOT" || spurs.Away.Code != "EVE" {
		t.Errorf("at Spurs %+v", spurs)
	}
	if cup.Competition != "Carabao Cup" {
		t.Errorf("third %+v", cup)
	}
	if p.Away[0].Away.Code != "CHE" || p.Away[0].HomeGoals != 2 || p.Away[0].AwayGoals != 1 {
		t.Errorf("Chelsea's first %+v", p.Away[0])
	}

	// The meetings, newest first, without the season in the competition's name
	if len(p.Meetings) != 5 {
		t.Fatalf("%d meetings", len(p.Meetings))
	}
	if m := p.Meetings[0]; m.Home.Code != "EVE" || m.HomeGoals != 3 || m.AwayGoals != 0 ||
		m.Competition != "English Premier League" || !m.Date.Equal(time.Date(2026, 3, 21, 17, 30, 0, 0, time.UTC)) {
		t.Errorf("last meeting %+v", m)
	}

	if len(p.News) != 4 || p.News[0].URL == "" || p.News[0].Headline == "" || p.News[0].Published.IsZero() {
		t.Errorf("news %+v", p.News)
	}
	for i := 1; i < len(p.News); i++ {
		if p.News[i].Published.After(p.News[i-1].Published) {
			t.Errorf("story %d is newer than the one before it", i)
		}
	}

	if p, err := ParsePreMatch([]byte(`{}`)); p != nil || err != nil {
		t.Errorf("nothing before a match: %v, %v", p, err)
	}
}
