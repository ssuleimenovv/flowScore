package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
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

	engine := flow.Engine{Params: flow.DefaultParams(), Speed: *speed, Tick: 5 * time.Second}
	updates := engine.Run(ctx, *matchID, events)

	lastMinute := -1
	for u := range updates {
		minute := int(u.At.Minutes())
		if u.Cause == nil && minute == lastMinute {
			continue // print ticks once per match minute
		}
		lastMinute = minute

		label := ""
		if u.Cause != nil {
			label = fmt.Sprintf("%s %s · %s", u.Cause.Type, u.Cause.Side, u.Cause.Player)
		}
		fmt.Printf("%s  %20s %3.0f │ %-3.0f %-20s  %s\n",
			event.FormatClock(u.At), bar(u.Home), u.Home, u.Away, bar(u.Away), label)
	}

	if ctx.Err() != nil {
		fmt.Println("stopped")
	}
}

// bar draws Flow as up to 20 blocks.
func bar(flow float64) string {
	return strings.Repeat("█", int(flow/5))
}
