package espn

import (
	"math"
	"testing"
)

func TestChance(t *testing.T) {
	for odds, want := range map[string]float64{"+215": 100.0 / 315, "-145": 145.0 / 245, "EVEN": 0.5, "+100": 0.5} {
		if got, ok := chance(odds); !ok || math.Abs(got-want) > 1e-9 {
			t.Errorf("chance(%q) = %.4f, %v; want %.4f", odds, got, ok, want)
		}
	}
	for _, odds := range []string{"", "OFF", "+50"} {
		if _, ok := chance(odds); ok {
			t.Errorf("chance(%q) is a chance", odds)
		}
	}
}

func TestMarket(t *testing.T) {
	price := func(open, close string) moneyline {
		var m moneyline
		m.Open.Odds, m.Close.Odds = open, close
		return m
	}
	// Everton – Chelsea on 17 October 2026: the closing prices, the draw's open one
	m := market(price("+175", "+215"), price("+275", ""), price("+130", "+115"))
	if m == nil || math.Abs(m.Home+m.Draw+m.Away-1) > 1e-9 {
		t.Fatalf("market %+v, want chances adding up to 1", m)
	}
	if round := func(p float64) int { return int(math.Round(p * 100)) }; round(m.Home) != 30 || round(m.Draw) != 25 || round(m.Away) != 44 {
		t.Errorf("market %+v, want 30/25/44", m)
	}
	if market(price("+175", ""), moneyline{}, price("+130", "")) != nil {
		t.Error("a market without the draw")
	}
}
