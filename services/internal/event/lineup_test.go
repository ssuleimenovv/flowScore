package event

import (
	"math"
	"strings"
	"testing"
)

func TestArrange(t *testing.T) {
	names := []string{"ST", "GK", "RB", "LB", "CM", "CB"}
	spots := []Spot{{4, 0}, {0, 0}, {1, 2}, {1, -2}, {2.5, 0}, {1, 0}}
	players := func() []Player {
		out := make([]Player, len(names))
		for i, n := range names {
			out[i] = Player{Name: n}
		}
		return out
	}
	// A formation that does not add up is ignored: the lines are the players' own
	for _, formation := range []string{"", "4-4-2", "3-1-1"} {
		p := players()
		Arrange(formation, p, spots)
		var got []string
		for _, q := range p {
			got = append(got, q.Name)
		}
		if want := "GK LB CB RB CM ST"; strings.Join(got, " ") != want {
			t.Errorf("%q: order %q, want %q", formation, strings.Join(got, " "), want)
		}
		if formation != "" {
			continue
		}
		// Four lines: the keeper, the back three across the pitch, then a line
		// every 12.7
		lb, rb, st := p[1].Pos, p[3].Pos, p[5].Pos
		if lb.X != 20 || lb.Y != 16 || rb.Y != 84 || !near(p[4].Pos.X, 32.67) || !near(st.X, 45.33) || st.Y != 50 {
			t.Errorf("LB %+v, RB %+v, CM %+v, ST %+v", *lb, *rb, *p[4].Pos, *st)
		}
	}
	// The formation's own lines: one striker and the midfielder on the same line
	p := players()
	Arrange("3-2", p, spots)
	if p[4].Pos.X != p[5].Pos.X || p[4].Pos.Y >= p[5].Pos.Y {
		t.Errorf("CM %+v and ST %+v, want one line, CM on the left", *p[4].Pos, *p[5].Pos)
	}
}

func near(a, b float64) bool {
	return math.Abs(a-b) < 0.01
}
