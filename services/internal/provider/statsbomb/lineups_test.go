package statsbomb

import (
	"math"
	"testing"
	"time"
)

func TestLoadLineups(t *testing.T) {
	// Manchester City – Arsenal, 2015/16, a demo match
	l, err := LoadLineups("../../../demo/statsbomb/lineups-3754314.json", "../../../demo/statsbomb/events-3754314.json", 36)
	if err != nil {
		t.Fatal(err)
	}
	city, arsenal := l.Home, l.Away
	if city.Formation != "4-4-2" || arsenal.Formation != "4-2-3-1" {
		t.Fatalf("formations %q and %q", city.Formation, arsenal.Formation)
	}
	if len(city.Players) != 18 || !city.Players[10].Starter || city.Players[11].Starter {
		t.Fatalf("%d players, the eleven first", len(city.Players))
	}

	// The keeper first, then each line from the left wing to the right one
	var names []string
	for _, p := range city.Players[:11] {
		names = append(names, p.Name)
	}
	want := []string{
		"Joe Hart",
		"Gaël Clichy", "Eliaquim Mangala", "Nicolás Otamendi", "Bacary Sagna",
		"Kevin De Bruyne", "Fernando Reges", "Fernandinho", "Jesús Navas",
		"Sergio Agüero", "Kelechi Iheanacho",
	}
	for i := range want {
		if names[i] != want[i] {
			t.Fatalf("starters %q, want %q", names, want)
		}
	}
	hart, clichy, aguero := city.Players[0], city.Players[1], city.Players[9]
	if hart.Pos.X != 5 || hart.Pos.Y != 50 || clichy.Pos.X != 16 || clichy.Pos.Y != 14 || math.Abs(aguero.Pos.X-45) > 0.01 || aguero.Pos.Y != 33 {
		t.Errorf("Hart at %+v, Clichy at %+v, Agüero at %+v", *hart.Pos, *clichy.Pos, *aguero.Pos)
	}

	// Iheanacho went off at 70:26 for Sterling, the first one on
	iheanacho, sterling := city.Players[10], city.Players[11]
	if iheanacho.Off == nil || iheanacho.Off.Period != 2 || iheanacho.Off.Elapsed != 70*time.Minute+26*time.Second {
		t.Errorf("Iheanacho off at %+v", iheanacho.Off)
	}
	if sterling.Name != "Raheem Sterling" || sterling.Number != "7" || sterling.On == nil || sterling.Pos != nil {
		t.Errorf("first on %+v", sterling)
	}
}
