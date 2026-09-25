package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/statsbomb"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	matchID := flag.String("match", "3754314", "StatsBomb match ID to replay")
	speed := flag.Float64("speed", 60, "1 = real time, 60 = one match minute per second")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	hub := live.NewHub()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /ws/matches/{matchId}/stream", live.StreamHandler(hub, []string{"localhost:5173"}))
	srv := &http.Server{Addr: *addr, Handler: mux}

	go func() {
		log.Printf("listening on %s, replay starts when the first client connects", *addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	go replay(ctx, hub, *matchID, *speed)

	<-ctx.Done()
	log.Print("shutting down")

	hub.Close() // WebSocket handlers see their channel closed and return
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

func replay(ctx context.Context, hub *live.Hub, matchID string, speed float64) {
	select {
	case <-hub.FirstSubscriber():
	case <-ctx.Done():
		return
	}

	provider := &statsbomb.Replay{
		MatchesPath: "data/statsbomb/matches-2-27.json",
		EventsDir:   "data/statsbomb",
		Speed:       speed,
	}
	events, err := provider.Stream(ctx, matchID)
	if err != nil {
		log.Printf("replay: %v", err)
		return
	}

	params := flow.DefaultParams()
	engine := flow.Engine{Params: params, Speed: speed, Tick: 5 * time.Second}
	updates := engine.Run(ctx, matchID, events)

	log.Printf("replaying match %s at x%.0f", matchID, speed)
	live.NewPublisher(hub, matchID, params).Run(updates)
	log.Print("replay finished")
}
