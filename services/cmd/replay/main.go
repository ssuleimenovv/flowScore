package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ssuleimenovv/flowscore/services/internal/provider/statsbomb"
)

func main() {
	matchID := flag.Int("match", 3754314, "StatsBomb match ID")
	matches := flag.String("matches", "data/statsbomb/matches-2-27.json", "matches file")
	events := flag.String("events", "data/statsbomb/events-3754314.json", "events file")
	flag.Parse()

	evs, err := statsbomb.LoadMatch(*matches, *events, *matchID)
	if err != nil {
		log.Fatal(err)
	}

	counts := map[string]int{}
	for _, e := range evs {
		counts[string(e.Type)]++

		xg := ""
		if e.XG != nil {
			xg = fmt.Sprintf("xG %.2f", *e.XG)
		}
		fmt.Printf("%s  %-15s %-4s  %-25s %s\n", e.Clock(), e.Type, e.Side, e.Player, xg)
	}
	fmt.Println(counts)
}
