package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
	"github.com/ssuleimenovv/flowscore/services/internal/season"
)

// seasons is a fake season.Seasons: what it has is read
type seasons map[string]season.League

func (s seasons) League(id string) (season.League, bool) {
	l, ok := s[id]
	return l, ok
}

func game(id string, round int, state, home, away, hs, as string, flow *season.MatchFlow) season.Game {
	return season.Game{
		Game: espn.Game{
			ID: id, State: state, Start: time.Date(2026, 10, round, 14, 0, 0, 0, time.UTC), Clock: "72'",
			Home: espn.Side{Team: event.Team{ID: home, Code: home, Name: home}, Score: hs},
			Away: espn.Side{Team: event.Team{ID: away, Code: away, Name: away}, Score: as},
		},
		Round: round,
		Flow:  flow,
	}
}

func premierLeague() season.League {
	return season.League{
		ID: "eng.1",
		Table: espn.Table{Season: 2026, Rows: []espn.Standing{
			{Rank: 1, Team: event.Team{ID: "ARS", Code: "ARS", Name: "ARS"}, Played: 2, Won: 2, For: 4, Points: 6, Zone: "Champions League"},
			{Rank: 2, Team: event.Team{ID: "CHE", Code: "CHE", Name: "CHE"}, Played: 1, Points: 0, Zone: "Relegation playoff"},
			{Rank: 3, Team: event.Team{ID: "LEE", Code: "LEE", Name: "LEE"}},
		}},
		Games: []season.Game{
			game("1", 1, "post", "ARS", "CHE", "3", "0", &season.MatchFlow{Home: 70, Away: 30, Peak: 99}),
			game("2", 2, "post", "LEE", "ARS", "0", "1", &season.MatchFlow{Home: 40, Away: 50, Peak: 80}),
			game("3", 2, "in", "CHE", "LEE", "1", "1", nil),
			game("4", 3, "pre", "ARS", "LEE", "", "", nil),
		},
		UpdatedAt: time.Date(2026, 10, 13, 12, 0, 0, 0, time.UTC),
	}
}

func getLeague(t *testing.T, mux *http.ServeMux, path string, want int, into any) {
	t.Helper()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	if rec.Code != want {
		t.Fatalf("GET %s: %d, want %d: %s", path, rec.Code, want, rec.Body)
	}
	if into != nil {
		if err := json.Unmarshal(rec.Body.Bytes(), into); err != nil {
			t.Fatal(err)
		}
	}
}

func leagueMux(s Seasons) *http.ServeMux {
	mux := http.NewServeMux()
	RegisterLeagues(mux, s, []string{"eng.1", "esp.1"})
	return mux
}

func TestLeagueList(t *testing.T) {
	var list struct {
		Items []struct {
			ID, Name string
			Ready    bool
		}
	}
	getLeague(t, leagueMux(seasons{"eng.1": premierLeague()}), "/api/v1/leagues", http.StatusOK, &list)
	// In the gateway's order; La Liga is not read yet
	if len(list.Items) != 2 || list.Items[0].Name != "Премьер-лига" || !list.Items[0].Ready || list.Items[1].Ready {
		t.Errorf("list %+v", list)
	}
}

func TestLeagueNotThere(t *testing.T) {
	mux := leagueMux(seasons{"eng.1": premierLeague()})
	getLeague(t, mux, "/api/v1/leagues/xxx.1", http.StatusNotFound, nil)
	// Followed, but its season is still being read
	getLeague(t, mux, "/api/v1/leagues/esp.1/standings", http.StatusServiceUnavailable, nil)
	getLeague(t, mux, "/api/v1/leagues/eng.1/rounds/4", http.StatusNotFound, nil)
	getLeague(t, mux, "/api/v1/leagues/eng.1/rounds/x", http.StatusNotFound, nil)
}

func TestLeagueHeader(t *testing.T) {
	var h struct {
		Season  string
		Round   int
		Rounds  int
		Records []struct {
			MatchID string
			Side    event.Side
			Flow    float64
		}
	}
	getLeague(t, leagueMux(seasons{"eng.1": premierLeague()}), "/api/v1/leagues/eng.1", http.StatusOK, &h)
	// Round 2 is under way; the best average Flow is Arsenal's 70 in match 1
	if h.Season != "2026/27" || h.Round != 2 || h.Rounds != 3 {
		t.Errorf("header %+v", h)
	}
	if len(h.Records) != 2 || h.Records[0].MatchID != "1" || h.Records[0].Side != event.Home || h.Records[0].Flow != 70 ||
		h.Records[1].Side != event.Away {
		t.Errorf("records %+v", h.Records)
	}
}

func TestStandings(t *testing.T) {
	var s struct {
		Items []struct {
			Rank   int
			Team   team
			Points int
			Zone   *string
			Flow   *float64
		}
	}
	getLeague(t, leagueMux(seasons{"eng.1": premierLeague()}), "/api/v1/leagues/eng.1/standings", http.StatusOK, &s)
	ars, che, lee := s.Items[0], s.Items[1], s.Items[2]
	// Arsenal: 70 at home, 50 away
	if ars.Team.Code != "ARS" || ars.Zone == nil || *ars.Zone != "champions" || ars.Flow == nil || *ars.Flow != 60 {
		t.Errorf("arsenal %+v", ars)
	}
	if che.Zone == nil || *che.Zone != "relegation" || che.Flow == nil || *che.Flow != 30 {
		t.Errorf("chelsea %+v", che)
	}
	// Leeds' match under way has no Flow yet; the one over has
	if lee.Zone != nil || lee.Flow == nil || *lee.Flow != 40 {
		t.Errorf("leeds %+v", lee)
	}
}

func TestRound(t *testing.T) {
	var r struct {
		Round int
		Items []struct {
			ID    string
			State string
			Clock string
			Score *struct{ Home, Away int }
			Flow  *struct{ Home, Away float64 }
		}
	}
	mux := leagueMux(seasons{"eng.1": premierLeague()})
	getLeague(t, mux, "/api/v1/leagues/eng.1/rounds/2", http.StatusOK, &r)
	if r.Round != 2 || len(r.Items) != 2 {
		t.Fatalf("round %+v", r)
	}
	over, live := r.Items[0], r.Items[1]
	if over.Clock != "" || over.Score == nil || over.Score.Away != 1 || over.Flow == nil || over.Flow.Away != 50 {
		t.Errorf("over %+v", over)
	}
	if live.Clock != "72'" || live.Score == nil || live.Flow != nil {
		t.Errorf("under way %+v", live)
	}

	getLeague(t, mux, "/api/v1/leagues/eng.1/rounds/3", http.StatusOK, &r)
	if len(r.Items) != 1 || r.Items[0].Score != nil {
		t.Errorf("to come %+v", r.Items)
	}
}
