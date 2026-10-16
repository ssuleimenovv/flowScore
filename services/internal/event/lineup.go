package event

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Lineups are the players of both teams in one match.
type Lineups struct {
	Home, Away Lineup
}

// Lineup is one team's: the starting eleven where they stand, then the bench.
type Lineup struct {
	Formation string   // "4-2-3-1"; "" when the source does not say
	Players   []Player // the eleven first, from the keeper forward, then the bench
}

type Player struct {
	ID      string
	Name    string
	Number  string // the shirt number, "" if unknown
	Starter bool
	Pos     *Position // where a starter stands, attacking left to right; nil on the bench
	On, Off *Moment   // when the player came on or went off; nil if not
}

// Moment is a point of the match, as the events have it.
type Moment struct {
	Period  int
	Elapsed time.Duration
}

// Spot is where a source puts a starter: Line counts from the keeper (0) to
// the strikers (4) and Flank from the left wing (-2) to the right one (2).
// Arrange turns the spots into places on the pitch.
type Spot struct {
	Line, Flank float64
}

// Arrange places the starters, the keeper first, on the lines of the
// formation and puts them in that order, each line from left to right. A
// formation that does not add up to the outfield players is ignored, and the
// lines are then the players' own.
func Arrange(formation string, starters []Player, spots []Spot) {
	order := make([]int, len(starters))
	for i := range order {
		order[i] = i
	}
	slices.SortStableFunc(order, func(a, b int) int {
		return cmp.Or(cmp.Compare(spots[a].Line, spots[b].Line), cmp.Compare(spots[a].Flank, spots[b].Flank))
	})
	if len(order) == 0 {
		return
	}

	// The keeper on his own line, then the formation's lines
	lines := [][]int{order[:1]}
	if sizes := lineSizes(formation); sum(sizes) == len(order)-1 {
		rest := order[1:]
		for _, n := range sizes {
			lines = append(lines, rest[:n])
			rest = rest[n:]
		}
	} else {
		for _, i := range order[1:] {
			last := lines[len(lines)-1]
			if len(lines) > 1 && spots[last[0]].Line == spots[i].Line {
				lines[len(lines)-1] = append(last, i)
			} else {
				lines = append(lines, []int{i})
			}
		}
	}

	placed := make([]Player, 0, len(starters))
	for n, line := range lines {
		slices.SortStableFunc(line, func(a, b int) int { return cmp.Compare(spots[a].Flank, spots[b].Flank) })
		for k, i := range line {
			p := starters[i]
			p.Pos = &Position{X: lineX(n, len(lines)), Y: flankY(k, len(line))}
			placed = append(placed, p)
		}
	}
	copy(starters, placed)
}

// SortBench puts the players who came on first, by when, and the rest by
// shirt number.
func SortBench(bench []Player) {
	slices.SortStableFunc(bench, func(a, b Player) int {
		switch {
		case a.On != nil && b.On != nil:
			return cmp.Or(cmp.Compare(a.On.Period, b.On.Period), cmp.Compare(a.On.Elapsed, b.On.Elapsed))
		case a.On != nil:
			return -1
		case b.On != nil:
			return 1
		}
		an, _ := strconv.Atoi(a.Number)
		bn, _ := strconv.Atoi(b.Number)
		return cmp.Compare(an, bn)
	})
}

// Each team stays in its own half, so the two elevens never run into each
// other: the keeper stands at 5, the back line at 16 and the front line at
// 45, short of the halfway line, the lines between them evenly apart.
func lineX(n, lines int) float64 {
	if n == 0 {
		return 5
	}
	if lines <= 2 {
		return 30
	}
	return 16 + float64(n-1)*29/float64(lines-2)
}

// The players of a line are spread around the middle, never wider than 14–86
func flankY(k, n int) float64 {
	if n == 1 {
		return 50
	}
	gap := min(34, 72/float64(n-1))
	return 50 + (float64(k)-float64(n-1)/2)*gap
}

// lineSizes reads "4-2-3-1" as 4, 2, 3, 1
func lineSizes(formation string) []int {
	var sizes []int
	for part := range strings.SplitSeq(formation, "-") {
		n, err := strconv.Atoi(part)
		if err != nil || n <= 0 {
			return nil
		}
		sizes = append(sizes, n)
	}
	return sizes
}

func sum(ns []int) int {
	s := 0
	for _, n := range ns {
		s += n
	}
	return s
}
