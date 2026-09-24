package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ssuleimenovv/flowscore/services/internal/provider"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/statsbomb"
)

func main() {
	matchID := flag.String("match", "3754314", "StatsBomb match ID")
	speed := flag.Float64("speed", 60, "1 = real time, 60 = one match minute per second")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var p provider.Provider = &statsbomb.Replay{
		MatchesPath: "data/statsbomb/matches-2-27.json",
		EventsDir:   "data/statsbomb",
		Speed:       *speed,
	}

	events, err := p.Stream(ctx, *matchID)
	if err != nil {
		log.Fatal(err)
	}

	for e := range events {
		fmt.Printf("%s  %-15s %-4s  %s\n", e.Clock(), e.Type, e.Side, e.Player)
	}

	if ctx.Err() != nil {
		fmt.Println("stopped")
	}
}
