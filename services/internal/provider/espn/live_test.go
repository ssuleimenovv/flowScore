package espn

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// summary is a summary as ESPN sends it, cut down to what ParseSummary reads
func summary(state, status string, period int, clock float64, home, away string) string {
	return fmt.Sprintf(`{
		"header": {
			"id": "401879268",
			"league": {"id": "700", "shortName": "Premier League"},
			"competitions": [{
				"date": "2026-10-10T11:30Z",
				"status": {"clock": %v, "period": %d, "type": {"name": %q, "state": %q}},
				"competitors": [
					{"homeAway": "home", "score": %q, "team": {"id": "359", "abbreviation": "ARS", "displayName": "Arsenal"}},
					{"homeAway": "away", "score": %q, "team": {"id": "357", "abbreviation": "LEE", "displayName": "Leeds United"}}
				]
			}]
		},
		"gameInfo": {"venue": {"fullName": "Emirates Stadium"}}
	}`, clock, period, status, state, home, away)
}

func TestParseSummary(t *testing.T) {
	s, err := ParseSummary([]byte(summary("in", "STATUS_SECOND_HALF", 2, 2770, "2", "1")))
	if err != nil {
		t.Fatal(err)
	}
	m := s.Match
	if m.ID != "401879268" || m.Source != "espn" || m.Competition != "Premier League" || m.Venue != "Emirates Stadium" ||
		m.Home != (event.Team{ID: "359", Code: "ARS", Name: "Arsenal"}) || m.Away.Code != "LEE" ||
		!m.KickoffAt.Equal(time.Date(2026, 10, 10, 11, 30, 0, 0, time.UTC)) {
		t.Errorf("match = %+v", m)
	}
	if s.State != "in" || s.Period != 2 || s.Clock != 46*time.Minute+10*time.Second || s.Score != (Score{2, 1}) {
		t.Errorf("summary = %+v", s)
	}

	// Before kick-off there is no score
	s, err = ParseSummary([]byte(summary("pre", "STATUS_SCHEDULED", 0, 0, "", "")))
	if err != nil || s.Score != (Score{}) || s.State != "pre" {
		t.Errorf("before kick-off: %+v, %v", s, err)
	}
}

func parse(t *testing.T, body string) Summary {
	t.Helper()
	s, err := ParseSummary([]byte(body))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestTracker(t *testing.T) {
	tr := tracker{matchID: "m1"}

	if events, over := tr.next(parse(t, summary("pre", "STATUS_SCHEDULED", 0, 0, "", ""))); events != nil || over {
		t.Errorf("before kick-off: %v, %v", events, over)
	}

	// 1:0 at 23:00: the goal, then where the clock is
	events, _ := tr.next(parse(t, summary("in", "STATUS_FIRST_HALF", 1, 1380, "1", "0")))
	if len(events) != 2 || events[0].Type != event.Goal || events[0].Side != event.Home ||
		events[0].Elapsed != 23*time.Minute || events[1].Type != event.Clock || events[1].Stopped {
		t.Errorf("first half: %+v", events)
	}

	// The same score again sends only the clock
	events, _ = tr.next(parse(t, summary("in", "STATUS_FIRST_HALF", 1, 1410, "1", "0")))
	if len(events) != 1 || events[0].Type != event.Clock {
		t.Errorf("no new goal: %+v", events)
	}

	// The break: the clock stands at 45:00 of the second half
	events, _ = tr.next(parse(t, summary("in", "STATUS_HALFTIME", 1, 2820, "1", "0")))
	if c := events[0]; c.Type != event.Clock || !c.Stopped || c.Period != 2 || c.Elapsed != 45*time.Minute {
		t.Errorf("break: %+v", c)
	}

	// Two goals at once, one each, then the final whistle
	events, over := tr.next(parse(t, summary("post", "STATUS_FULL_TIME", 2, 5580, "2", "1")))
	if len(events) != 3 || events[0].ID != "m1:goal:home:2" || events[1].ID != "m1:goal:away:1" || !over {
		t.Errorf("full time: %+v, over %v", events, over)
	}
}

func TestLiveStreamsUntilTheFinalWhistle(t *testing.T) {
	polls := []string{
		summary("in", "STATUS_FIRST_HALF", 1, 600, "0", "0"),
		summary("in", "STATUS_FIRST_HALF", 1, 630, "1", "0"),
		summary("post", "STATUS_FULL_TIME", 2, 5580, "1", "0"),
	}
	var n atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		i := min(int(n.Add(1))-1, len(polls)-1)
		w.Write([]byte(polls[i]))
	}))
	defer srv.Close()

	l := Live{Client: &Client{HTTP: srv.Client(), BaseURL: srv.URL}, League: "eng.1", Every: time.Millisecond}
	events, err := l.Stream(context.Background(), "401879268")
	if err != nil {
		t.Fatal(err)
	}
	var types []event.Type
	for e := range events { // the channel closes at the final whistle
		types = append(types, e.Type)
	}
	want := []event.Type{event.Clock, event.Goal, event.Clock, event.Clock}
	if fmt.Sprint(types) != fmt.Sprint(want) {
		t.Errorf("events %v, want %v", types, want)
	}
}
