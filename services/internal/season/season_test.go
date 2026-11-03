package season

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
)

func TestMeasure(t *testing.T) {
	share := 0.5
	events := []event.Event{{Type: event.Possession, Period: 1, Elapsed: time.Minute, HomeShare: &share}}
	// Five home shots in a row, then nothing till the end
	for m := 2; m <= 6; m++ {
		events = append(events, event.Event{Type: event.ShotOnTarget, Side: event.Home, Period: 1, Elapsed: time.Duration(m) * time.Minute})
	}
	events = append(events, event.Event{Type: event.Possession, Period: 2, Elapsed: 90 * time.Minute, HomeShare: &share})

	f := Measure(events, flow.DefaultParams(), nil)
	if f.PeakSide != event.Home || f.PeakMinute != 6 || f.Peak < 70 {
		t.Errorf("peak %.1f of %s at %d', want the home team's at 6'", f.Peak, f.PeakSide, f.PeakMinute)
	}
	// The shots fade, so on the whole the home team is well below its peak but above the away team
	if f.Home >= f.Peak || f.Home <= f.Away {
		t.Errorf("average %.1f : %.1f, peak %.1f", f.Home, f.Away, f.Peak)
	}
}

func TestNumberRounds(t *testing.T) {
	game := func(home, away string) Game {
		return Game{Game: espn.Game{Home: espn.Side{Team: event.Team{ID: home}}, Away: espn.Side{Team: event.Team{ID: away}}}}
	}
	// Round 1: A–B, C–D; round 2: A–C; B–D was put off and played after round 2
	games := []Game{game("A", "B"), game("C", "D"), game("A", "C"), game("B", "D")}
	numberRounds(games)
	var got []int
	for _, g := range games {
		got = append(got, g.Round)
	}
	if want := []int{1, 1, 2, 2}; !equal(got, want) {
		t.Errorf("rounds %v, want %v", got, want)
	}
}

func TestKnockoutsHaveNoRound(t *testing.T) {
	game := func(stage, home, away string) Game {
		return Game{Game: espn.Game{Stage: stage, Home: espn.Side{Team: event.Team{ID: home}}, Away: espn.Side{Team: event.Team{ID: away}}}}
	}
	games := []Game{game("league-phase", "A", "B"), game("round-of-16", "A", "C"), game("league-phase", "C", "D")}
	numberRounds(games)
	if games[0].Round != 1 || games[1].Round != 0 || games[2].Round != 1 {
		t.Errorf("rounds %d, %d, %d; want 1, 0, 1", games[0].Round, games[1].Round, games[2].Round)
	}
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ESPN as the season needs it: a table, a calendar of one day, a finished
// match and one to come
const (
	table = `{"season": {"year": 2026}, "children": [{"standings": {"entries": [
		{"team": {"id": "359", "abbreviation": "ARS", "displayName": "Arsenal"}, "stats": [{"name": "rank", "value": 1}]}]}}]}`
	today = `{"leagues": [{"calendar": ["2026-10-10T07:00Z"]}], "events": []}`
	day   = `{"events": [
		{"id": "1", "date": "2026-10-10T14:00Z", "status": {"type": {"state": "post"}}, "competitions": [{"competitors": [
			{"homeAway": "home", "score": "1", "team": {"id": "359", "displayName": "Arsenal"}},
			{"homeAway": "away", "score": "0", "team": {"id": "357", "displayName": "Leeds United"}}]}]},
		{"id": "2", "date": "2026-10-10T16:30Z", "status": {"type": {"state": "pre"}}, "competitions": [{"competitors": [
			{"homeAway": "home", "team": {"id": "363", "displayName": "Chelsea"}},
			{"homeAway": "away", "team": {"id": "359", "displayName": "Arsenal"}}]}]}]}`
	summary = `{
		"header": {"id": "1", "league": {"id": "700"}, "competitions": [{"date": "2026-10-10T14:00Z",
			"status": {"type": {"state": "post", "name": "STATUS_FULL_TIME"}},
			"competitors": [{"homeAway": "home", "score": "1", "team": {"id": "359", "displayName": "Arsenal"}},
				{"homeAway": "away", "score": "0", "team": {"id": "357", "displayName": "Leeds United"}}]}]},
		"boxscore": {"teams": [{"homeAway": "home", "statistics": [{"name": "possessionPct", "displayValue": "60"}]},
			{"homeAway": "away", "statistics": [{"name": "possessionPct", "displayValue": "40"}]}]},
		"commentary": [{"play": {"id": "9", "type": {"type": "goal"}, "period": {"number": 1},
			"clock": {"value": 600, "displayValue": "10'"}, "team": {"displayName": "Arsenal"},
			"participants": [{"athlete": {"displayName": "Bukayo Saka"}}], "fieldPositionX": 88, "fieldPositionY": 50}}]
	}`
)

func TestRefresh(t *testing.T) {
	var summaries atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/eng.1/standings":
			io.WriteString(w, table)
		case "/eng.1/scoreboard":
			io.WriteString(w, today)
		case "/eng.1/scoreboard?dates=20261010":
			io.WriteString(w, day)
		case "/eng.1/summary?event=1":
			summaries.Add(1)
			io.WriteString(w, summary)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	s := &Seasons{Client: &espn.Client{HTTP: srv.Client(), BaseURL: srv.URL, StandingsURL: srv.URL}, Params: flow.DefaultParams()}
	if _, ok := s.League("eng.1"); ok {
		t.Fatal("a league before it was read")
	}
	for range 2 {
		if err := s.Refresh(context.Background(), "eng.1"); err != nil {
			t.Fatal(err)
		}
	}
	l, ok := s.League("eng.1")
	if !ok || l.Table.Season != 2026 || len(l.Games) != 2 {
		t.Fatalf("league %+v", l)
	}
	played, next := l.Games[0], l.Games[1]
	if played.Flow == nil || played.Flow.Home <= played.Flow.Away || played.Round != 1 {
		t.Errorf("played %+v, flow %+v", played, played.Flow)
	}
	// Arsenal's second match is the second round; it has no Flow yet
	if next.Flow != nil || next.Round != 2 {
		t.Errorf("to come %+v", next)
	}
	// A match over is read once, not on every refresh
	if n := summaries.Load(); n != 1 {
		t.Errorf("summary read %d times", n)
	}
}

// A cup lists stages, not days: its season is read a year at a time,
// without the last season's final that the year still has
func TestRefreshCup(t *testing.T) {
	final := `{"id": "0", "date": "2026-05-30T19:00Z", "status": {"type": {"state": "post"}}, "competitions": [{"competitors": [
		{"homeAway": "home", "score": "2", "team": {"id": "83", "displayName": "Barcelona"}},
		{"homeAway": "away", "score": "1", "team": {"id": "160", "displayName": "Paris Saint-Germain"}}]}]},`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/uefa.champions/standings":
			io.WriteString(w, table)
		case "/uefa.champions/scoreboard":
			io.WriteString(w, `{"leagues": [{"season": {"startDate": "2026-07-01T04:00Z"}, "calendar": [{"label": "League Phase"}]}], "events": []}`)
		case "/uefa.champions/scoreboard?dates=2026&limit=1000":
			io.WriteString(w, strings.Replace(day, `"events": [`, `"events": [`+final, 1))
		case "/uefa.champions/summary?event=1":
			io.WriteString(w, summary)
		default:
			if strings.HasSuffix(r.URL.RawQuery, "limit=1000") {
				io.WriteString(w, `{"events": []}`) // a year with nothing yet
				return
			}
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	s := &Seasons{Client: &espn.Client{HTTP: srv.Client(), BaseURL: srv.URL, StandingsURL: srv.URL}, Params: flow.DefaultParams()}
	if err := s.Refresh(context.Background(), "uefa.champions"); err != nil {
		t.Fatal(err)
	}
	l, _ := s.League("uefa.champions")
	if len(l.Games) != 2 || l.Games[0].ID != "1" || l.Games[0].Flow == nil || l.Games[0].Round != 1 {
		t.Errorf("games %+v", l.Games)
	}
}
