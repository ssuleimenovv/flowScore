package espn

import (
	"strconv"
	"strings"
)

// Market is what the bookmakers ESPN shows think of a match before
// kick-off: the chances of the three results, without their margin.
type Market struct {
	Home, Draw, Away float64
}

// moneyline is a price of ESPN's odds: "+215", "-145" or "EVEN"
type moneyline struct {
	Open, Close struct {
		Odds string `json:"odds"`
	}
}

// market turns the last prices of the three results into chances, nil
// when any of them is missing.
func market(home, draw, away moneyline) *Market {
	var p [3]float64
	for i, m := range []moneyline{home, draw, away} {
		odds := m.Close.Odds
		if odds == "" {
			odds = m.Open.Odds
		}
		v, ok := chance(odds)
		if !ok {
			return nil
		}
		p[i] = v
	}
	sum := p[0] + p[1] + p[2]
	return &Market{Home: p[0] / sum, Draw: p[1] / sum, Away: p[2] / sum}
}

// chance is the chance an American price stands for: +215 pays 215 on 100,
// so 100 / 315; -145 takes 145 to win 100, so 145 / 245.
func chance(odds string) (float64, bool) {
	if strings.EqualFold(odds, "EVEN") {
		return 0.5, true
	}
	v, err := strconv.Atoi(strings.TrimPrefix(odds, "+"))
	switch {
	case err != nil || (v > -100 && v < 100):
		return 0, false
	case v > 0:
		return 100 / float64(v+100), true
	default:
		return float64(-v) / float64(-v+100), true
	}
}
