package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ssuleimenovv/flowscore/services/internal/flow"
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

	state := flow.NewState(flow.DefaultParams())
	for e := range events {
		state.Apply(e)
		home, away := state.Flow()
		fmt.Printf("%s  %-15s %-4s  %-34s  flow %3.0f : %-3.0f\n",
			e.Clock(), e.Type, e.Side, e.Player, home, away)
	}

	if ctx.Err() != nil {
		fmt.Println("stopped")
	}
}
