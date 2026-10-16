package espn

import (
	"math"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func TestParseRosters(t *testing.T) {
	// Coventry City – Newcastle United, 2026/27, the rosters of its summary
	body, err := os.ReadFile("testdata/rosters-401879267.json")
	if err != nil {
		t.Fatal(err)
	}
	l, err := ParseRosters(body)
	if err != nil || l == nil {
		t.Fatalf("lineups %v, %v", l, err)
	}
	coventry, newcastle := l.Home, l.Away
	if coventry.Formation != "4-2-3-1" || len(coventry.Players) != 20 {
		t.Fatalf("Coventry %s, %d players", coventry.Formation, len(coventry.Players))
	}

	// ESPN's "LM" and "RM" of a 4-2-3-1 are its holding pair, not wingers
	var names []string
	for _, p := range coventry.Players[:11] {
		names = append(names, p.Name)
	}
	want := "Carl Rushworth, Jay Da Silva, Ethan Pinnock, Bobby Thomas, Milan van Ewijk, " +
		"Matt Grimes, Caleb Yirenkyi, Ephron Mason-Clark, Jack Rudoni, Brandon Thomas-Asante, Sidiki Cherif"
	if got := strings.Join(names, ", "); got != want {
		t.Fatalf("starters %s\nwant %s", got, want)
	}
	if grimes := coventry.Players[5]; math.Abs(grimes.Pos.X-25.67) > 0.01 || grimes.Pos.Y != 33 {
		t.Errorf("Grimes at %+v", *grimes.Pos)
	}

	// Rudoni went off at 83' for Torp; the bench starts with Sakamoto, on at 61'
	rudoni, sakamoto := coventry.Players[8], coventry.Players[11]
	if rudoni.Off == nil || *rudoni.Off != (event.Moment{Period: 2, Elapsed: 82 * time.Minute}) || rudoni.Number != "5" {
		t.Errorf("Rudoni %+v off at %+v", rudoni, rudoni.Off)
	}
	if sakamoto.Name != "Tatsuhiro Sakamoto" || sakamoto.On == nil || sakamoto.On.Elapsed != 60*time.Minute || sakamoto.Pos != nil {
		t.Errorf("first on %+v", sakamoto)
	}
	// The ones who did not come on, by shirt number
	if last := coventry.Players[19]; last.Name != "Gustavo Hamer" || last.On != nil {
		t.Errorf("last on the bench %+v", last)
	}

	// Touré came on at 59' and went off again at 90'
	for _, p := range newcastle.Players {
		if p.Name == "Bazoumana Touré" && (p.On == nil || p.Off == nil || p.On.Elapsed != 58*time.Minute || p.Off.Elapsed != 89*time.Minute) {
			t.Errorf("Touré on %+v, off %+v", p.On, p.Off)
		}
	}

	// Before ESPN has them
	if l, err := ParseRosters([]byte(`{"rosters": []}`)); l != nil || err != nil {
		t.Errorf("no rosters: %v, %v", l, err)
	}
}

func TestMoment(t *testing.T) {
	for clock, want := range map[string]event.Moment{
		"1'":      {Period: 1, Elapsed: 0},
		"45'+2'":  {Period: 1, Elapsed: 46 * time.Minute},
		"46'":     {Period: 2, Elapsed: 45 * time.Minute},
		"90'+2'":  {Period: 2, Elapsed: 91 * time.Minute},
		"105'+1'": {Period: 3, Elapsed: 105 * time.Minute},
		"118'":    {Period: 4, Elapsed: 117 * time.Minute},
	} {
		if got := moment(clock); got == nil || *got != want {
			t.Errorf("moment(%q) = %+v, want %+v", clock, got, want)
		}
	}
	for _, clock := range []string{"", "HT", "0'", "90'+x'"} {
		if got := moment(clock); got != nil {
			t.Errorf("moment(%q) = %+v", clock, got)
		}
	}
}
