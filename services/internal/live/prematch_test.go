package live

import (
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
)

// flows is a fake season: the Flow of the matches it has read
type flows map[string][2]float64

func (f flows) PlayedFlow(id string) (home, away float64, ok bool) {
	v, ok := f[id]
	return v[0], v[1], ok
}

func TestPreMatchOf(t *testing.T) {
	everton := event.Team{ID: "368", Code: "EVE", Name: "Everton"}
	chelsea := event.Team{ID: "363", Code: "CHE", Name: "Chelsea"}
	spurs := event.Team{ID: "367", Code: "TOT", Name: "Tottenham Hotspur"}
	wolves := event.Team{ID: "380", Code: "WOL", Name: "Wolverhampton Wanderers"}
	day := time.Date(2026, 9, 12, 16, 30, 0, 0, time.UTC)

	store := NewStore()
	p := NewPublisher(NewHub(), store, "m1", flow.DefaultParams())
	p.UsePreMatch(&event.PreMatch{
		Home: []event.PastMatch{
			// A 0-0 at Spurs, and a cup win the season did not read
			{ID: "a", Date: day, Competition: "Premier League", Home: spurs, Away: everton},
			{ID: "b", Date: day, Competition: "Carabao Cup", Home: everton, Away: wolves, HomeGoals: 1},
		},
		Meetings: []event.PastMatch{
			{ID: "c", Home: everton, Away: chelsea, HomeGoals: 3},
			{ID: "d", Home: chelsea, Away: everton, HomeGoals: 2},
			{ID: "e", Home: chelsea, Away: everton},
		},
		News: []event.Article{{Headline: "Premier League Top 50", URL: "https://www.espn.com/x", Published: day}},
	})
	// The pre-match info outlives the reset of a new kick-off, as the lineups do
	p.Schedule(event.Match{ID: "m1", Home: everton, Away: chelsea})

	snap, _ := store.Get("m1")
	pm, ok := PreMatchOf(snap, flows{"a": {58, 44}})
	if !ok {
		t.Fatal("no pre-match info")
	}
	form := pm.Home
	if len(form.Matches) != 2 || form.Matches[0].Result != "draw" || form.Matches[1].Result != "win" {
		t.Fatalf("Everton's form %+v", form.Matches)
	}
	// Everton were the away team at Spurs: their Flow is the away one, of the one match read
	if form.Flow == nil || *form.Flow != 44 || form.FlowGames != 1 || form.Matches[1].Flow != nil {
		t.Errorf("Everton's Flow %v over %d", form.Flow, form.FlowGames)
	}
	if len(pm.Away.Matches) != 0 || pm.Away.Flow != nil {
		t.Errorf("Chelsea's form %+v", pm.Away)
	}
	// Everton won one, lost one, drew one, whoever was at home
	if m := pm.Meetings; m.HomeWins != 1 || m.AwayWins != 1 || m.Draws != 1 || len(m.Matches) != 3 {
		t.Errorf("meetings %+v", m)
	}
	if len(pm.News) != 1 || pm.News[0].URL != "https://www.espn.com/x" {
		t.Errorf("news %+v", pm.News)
	}

	if _, ok := PreMatchOf(Snapshot{}, nil); ok {
		t.Error("pre-match info of a demo match")
	}
}
