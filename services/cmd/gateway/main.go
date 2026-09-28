package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/api"
	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/statsbomb"
)

const matchesPath = "data/statsbomb/matches-2-27.json"

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	matchID := flag.String("match", "3754314", "StatsBomb match ID to replay")
	speed := flag.Float64("speed", 60, "1 = real time, 60 = one match minute per second")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	id, err := strconv.Atoi(*matchID)
	if err != nil {
		log.Fatalf("match ID: %v", err)
	}
	info, err := statsbomb.LoadMatchInfo(matchesPath, id)
	if err != nil {
		log.Fatal(err)
	}

	hub := live.NewHub()
	store := live.NewStore()
	store.Schedule(info) // the page can load the match before anyone starts the replay

	mux := http.NewServeMux()
	live.Register(mux, hub, []string{"localhost:5173"})
	api.Register(mux, store)
	srv := &http.Server{Addr: *addr, Handler: mux}

	go func() {
		log.Printf("listening on %s, replay starts when the first client connects", *addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	go loopReplay(ctx, hub, store, info, *speed)

	<-ctx.Done()
	log.Print("shutting down")

	hub.Close() // WebSocket handlers see their channel closed and return
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

// loopReplay plays the match again and again while anyone watches it.
func loopReplay(ctx context.Context, hub *live.Hub, store *live.Store, info event.Match, speed float64) {
	matchID := info.ID
	params := flow.DefaultParams()
	publisher := live.NewPublisher(hub, store, matchID, params)

	for round := 0; ctx.Err() == nil; round++ {
		for hub.Count(matchID) == 0 {
			select {
			case <-hub.Joined():
			case <-ctx.Done():
				return
			}
		}

		// Kick-off goes into the store before the viewers of the last replay are
		// dropped, so the snapshot they reload is already the new match
		publisher.Start(info)
		if round > 0 {
			hub.CloseMatch(matchID)
		}
		if err := replayOnce(ctx, publisher, params, info, speed); err != nil {
			log.Printf("replay: %v", err)
			return
		}
	}
}

func replayOnce(ctx context.Context, publisher *live.Publisher, params flow.Params, info event.Match, speed float64) error {
	matchID := info.ID
	provider := &statsbomb.Replay{
		MatchesPath: matchesPath,
		EventsDir:   "data/statsbomb",
		Speed:       speed,
	}
	events, err := provider.Stream(ctx, matchID)
	if err != nil {
		return err
	}

	engine := flow.Engine{Params: params, Speed: speed, Tick: 5 * time.Second}
	log.Printf("replaying match %s at x%.0f", matchID, speed)
	publisher.Run(engine.Run(ctx, matchID, events))
	log.Print("replay finished")
	return nil
}
