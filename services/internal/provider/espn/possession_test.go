package espn

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func TestBall(t *testing.T) {
	var b ball
	share := func(v float64) *float64 { return &v }

	// Before ESPN counts the ball there is nothing to send
	if e := b.next("m1", 1, 10*time.Minute, nil); e != nil {
		t.Errorf("no share: %+v", e)
	}

	// The first reading is the share of every minute so far
	e := b.next("m1", 1, 10*time.Minute, share(0.6))
	if len(e) != 10 || e[0].ID != "m1:possession:1" || e[9].Elapsed != 10*time.Minute ||
		*e[9].HomeShare != 0.6 || e[9].Period != 1 {
		t.Fatalf("first reading: %d events, %+v", len(e), e)
	}

	// Half a minute later no minute has ended
	if e := b.next("m1", 1, 10*time.Minute+30*time.Second, share(0.61)); e != nil {
		t.Errorf("within a minute: %+v", e)
	}

	// Ten minutes at 40% take the whole match from 60% to 50%
	e = b.next("m1", 1, 20*time.Minute, share(0.5))
	if len(e) != 10 || e[0].ID != "m1:possession:11" || math.Abs(*e[0].HomeShare-0.4) > 1e-9 {
		t.Errorf("second reading: %d events, first %+v", len(e), e[0])
	}

	// A share that has not moved waits, for ESPN may not have counted yet...
	if e := b.next("m1", 1, 22*time.Minute, share(0.5)); e != nil {
		t.Errorf("not moved: %+v", e)
	}
	// ...and is taken once it has stood for five minutes
	if e := b.next("m1", 1, 25*time.Minute, share(0.5)); len(e) != 5 || *e[0].HomeShare != 0.5 {
		t.Errorf("stood still: %+v", e)
	}

	// Read in the second half, the minutes of the first keep their period
	e = b.next("m1", 2, 47*time.Minute, share(0.45))
	if len(e) != 22 || e[19].Elapsed != 45*time.Minute || e[19].Period != 1 || e[20].Period != 2 {
		t.Errorf("across the break: %d events", len(e))
	}
}

// boxscore puts ESPN's possession into a summary
func boxscore(body, home, away string) string {
	return strings.Replace(body, `"gameInfo"`, `"boxscore": {"teams": [
		{"homeAway": "home", "statistics": [{"name": "foulsCommitted", "displayValue": "3"}, {"name": "possessionPct", "displayValue": "`+home+`"}]},
		{"homeAway": "away", "statistics": [{"name": "possessionPct", "displayValue": "`+away+`"}]}
	]}, "gameInfo"`, 1)
}

func TestParseSummaryPossession(t *testing.T) {
	s := parse(t, boxscore(summary("in", "STATUS_SECOND_HALF", 2, 3600, "4", "1"), "33.3", "66.7"))
	if s.Possession == nil || math.Abs(*s.Possession-0.333) > 1e-9 {
		t.Errorf("possession = %v", s.Possession)
	}
	// Before kick-off ESPN has no share, or 0 for both
	if s := parse(t, summary("pre", "STATUS_SCHEDULED", 0, 0, "", "")); s.Possession != nil {
		t.Errorf("no boxscore: %v", *s.Possession)
	}
	if s := parse(t, boxscore(summary("pre", "STATUS_SCHEDULED", 0, 0, "", ""), "0", "0")); s.Possession != nil {
		t.Errorf("0 for both: %v", *s.Possession)
	}
}

// Joining a match under way, the minutes of possession so far go between its
// plays, in match order, before the clock
func TestTrackerSendsPossessionInMatchOrder(t *testing.T) {
	tr := tracker{matchID: "m1", sent: map[string]bool{}}
	shot := event.Event{ID: "m1:1", Type: event.ShotOnTarget, Period: 1, Elapsed: 90 * time.Second}

	events, _ := tr.next(parse(t, boxscore(summary("in", "STATUS_FIRST_HALF", 1, 180, "0", "0"), "60", "40")), []event.Event{shot})
	var got []string
	for _, e := range events {
		got = append(got, fmt.Sprintf("%s %s", e.Type, e.Clock()))
	}
	want := "[possession 01:00 shot_on_target 01:30 possession 02:00 possession 03:00 clock 03:00]"
	if fmt.Sprint(got) != want {
		t.Errorf("events %v, want %s", got, want)
	}
}
