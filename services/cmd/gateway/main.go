package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/api"
	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/statsbomb"
)

// Every setting is a flag with a default for local development. A host such
// as Render sets the environment instead, so each default reads it first.
func main() {
	addr := flag.String("addr", ":"+env("PORT", "8080"), "listen address")
	matchID := flag.String("match", env("FLOWSCORE_MATCH", "3754314"), "StatsBomb match ID to replay")
	speed := flag.Float64("speed", 60, "1 = real time, 60 = one match minute per second")
	dataDir := flag.String("data", env("FLOWSCORE_DATA", "data/statsbomb"), "folder with the matches, events-<id> and lineups-<id> files")
	matchesFile := flag.String("matches", env("FLOWSCORE_MATCHES", "matches-2-27.json"), "matches file in the data folder")
	modelPath := flag.String("model", env("FLOWSCORE_MODEL", "../ai/prediction/model.json"), "outcome model from ai/prediction/fit.py")
	originList := flag.String("origins", env("FLOWSCORE_ORIGINS", "localhost:5173,*:5173"), "host patterns of the sites allowed to call the API, comma-separated")
	flag.Parse()

	// Ctrl+C locally, SIGTERM when a host stops or redeploys the container
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	id, err := strconv.Atoi(*matchID)
	if err != nil {
		log.Fatalf("match ID: %v", err)
	}
	replay := statsbomb.Replay{
		MatchesPath: filepath.Join(*dataDir, *matchesFile),
		EventsDir:   *dataDir,
		Speed:       *speed,
	}
	info, err := statsbomb.LoadMatchInfo(replay.MatchesPath, id)
	if err != nil {
		log.Fatal(err)
	}

	model, err := predict.Load(*modelPath)
	if err != nil {
		log.Fatalf("outcome model: %v", err)
	}

	hub := live.NewHub()
	store := live.NewStore()
	params := flow.DefaultParams()
	publisher := live.NewPublisher(hub, store, info.ID, params)
	publisher.UseModel(&model)
	publisher.Schedule(info) // the page can load the match before anyone starts the replay

	origins := strings.Split(*originList, ",")
	mux := http.NewServeMux()
	// The host polls it to know the service is up
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	live.Register(mux, hub, origins)
	api.Register(mux, store)
	srv := &http.Server{
		Addr:              *addr,
		Handler:           api.CORS(mux, origins),
		ReadHeaderTimeout: 10 * time.Second, // a client that never finishes its headers is dropped
	}

	go func() {
		log.Printf("listening on %s, replay starts when the first client connects", *addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	go loopReplay(ctx, hub, publisher, params, info, &replay)

	<-ctx.Done()
	log.Print("shutting down")

	hub.Close() // WebSocket handlers see their channel closed and return
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

// env returns the environment variable, or fallback when it is not set.
func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// loopReplay plays the match again and again while anyone watches it.
func loopReplay(ctx context.Context, hub *live.Hub, publisher *live.Publisher, params flow.Params, info event.Match, replay *statsbomb.Replay) {
	matchID := info.ID
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
		if err := replayOnce(ctx, publisher, params, info, replay); err != nil {
			log.Printf("replay: %v", err)
			return
		}
	}
}

func replayOnce(ctx context.Context, publisher *live.Publisher, params flow.Params, info event.Match, replay *statsbomb.Replay) error {
	matchID := info.ID
	events, err := replay.Stream(ctx, matchID)
	if err != nil {
		return err
	}

	engine := flow.Engine{Params: params, Speed: replay.Speed, Tick: 5 * time.Second}
	log.Printf("replaying match %s at x%.0f", matchID, replay.Speed)
	publisher.Run(engine.Run(ctx, matchID, events))
	log.Print("replay finished")
	return nil
}
