package espn

import (
	"os"
	"testing"
)

func TestParseSquad(t *testing.T) {
	// Manchester City, 2026/27: one player of each line of its roster
	body, err := os.ReadFile("testdata/squad-382.json")
	if err != nil {
		t.Fatal(err)
	}
	players, err := ParseSquad(body)
	if err != nil || len(players) != 4 {
		t.Fatalf("%d players, %v", len(players), err)
	}
	dias := players[1]
	if dias.Name != "Rúben Dias" || dias.Jersey != "3" || dias.Position != "defender" || dias.Age != 29 || dias.Country != "Portugal" {
		t.Errorf("Dias %+v", dias)
	}
	want := []string{"goalkeeper", "defender", "midfielder", "forward"}
	for i, p := range players {
		if p.Position != want[i] {
			t.Errorf("%s plays %q, want %q", p.Name, p.Position, want[i])
		}
	}
}
