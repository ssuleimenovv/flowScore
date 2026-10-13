package espn

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// summary is a summary as ESPN sends it, cut down to what ParseSummary and
// ParsePlays read
func summary(state, status string, period int, clock float64, home, away string, plays ...string) string {
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
		"gameInfo": {"venue": {"fullName": "Emirates Stadium"}},
		"commentary": [%s]
	}`, clock, period, status, state, home, away, strings.Join(plays, ","))
}

// play is one play of a summary's commentary
func play(id, typ string, period int, clock float64, team string) string {
	return fmt.Sprintf(`{"play": {"id": %q, "type": {"type": %q}, "period": {"number": %d},
		"clock": {"value": %v, "displayValue": ""}, "team": {"displayName": %q}}}`, id, typ, period, clock, team)
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
	tr := tracker{matchID: "m1", sent: map[string]bool{}}
	pen := event.Event{ID: "m1:1", Type: event.Goal, Side: event.Home, Period: 1, Elapsed: 6 * time.Minute}
	late := event.Event{ID: "m1:2", Type: event.ShotOnTarget, Side: event.Away, Period: 1, Elapsed: 47 * time.Minute}
	shot := event.Event{ID: "m1:3", Type: event.ShotBlocked, Side: event.Away, Period: 2, Elapsed: 50 * time.Minute}

	if events, over := tr.next(parse(t, summary("pre", "STATUS_SCHEDULED", 0, 0, "", "")), nil); events != nil || over {
		t.Errorf("before kick-off: %v, %v", events, over)
	}

	// The penalty, then where the clock is
	events, _ := tr.next(parse(t, summary("in", "STATUS_FIRST_HALF", 1, 420, "1", "0")), []event.Event{pen})
	if len(events) != 2 || events[0].ID != "m1:1" || events[1].Type != event.Clock || events[1].Stopped {
		t.Errorf("first half: %+v", events)
	}

	// A play sent already is not sent again
	events, _ = tr.next(parse(t, summary("in", "STATUS_FIRST_HALF", 1, 450, "1", "0")), []event.Event{pen})
	if len(events) != 1 || events[0].Type != event.Clock {
		t.Errorf("nothing new: %+v", events)
	}

	// The break: the clock stands at 45:00 of the second half
	events, _ = tr.next(parse(t, summary("in", "STATUS_HALFTIME", 1, 2700, "1", "0")), []event.Event{pen})
	if c := events[0]; c.Type != event.Clock || !c.Stopped || c.Period != 2 || c.Elapsed != 45*time.Minute {
		t.Errorf("break: %+v", c)
	}

	// A first-half play written down after the break comes too late, and a
	// finished match has no clock to send
	events, over := tr.next(parse(t, summary("post", "STATUS_FULL_TIME", 0, 0, "1", "0")), []event.Event{pen, late, shot})
	if len(events) != 1 || events[0].ID != "m1:3" || !over {
		t.Errorf("full time: %+v, over %v", events, over)
	}

}

func TestLiveStreamsUntilTheFinalWhistle(t *testing.T) {
	goal := play("1", "goal", 1, 630, "Arsenal")
	polls := []string{
		summary("in", "STATUS_FIRST_HALF", 1, 600, "0", "0"),
		summary("in", "STATUS_FIRST_HALF", 1, 630, "1", "0", goal),
		summary("post", "STATUS_FULL_TIME", 0, 0, "1", "0", goal),
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
	want := []event.Type{event.Clock, event.Goal, event.Clock}
	if fmt.Sprint(types) != fmt.Sprint(want) {
		t.Errorf("events %v, want %v", types, want)
	}
}
