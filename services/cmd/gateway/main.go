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
	"github.com/ssuleimenovv/flowscore/services/internal/insight"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/statsbomb"
	"github.com/ssuleimenovv/flowscore/services/internal/xg"
)

// The demo matches: Premier League 2015/16 games with goals for both sides
const demoMatches = "3754314,3754348,3754305,3754239,3754208,3754174"

// Every setting is a flag with a default for local development. A host such
// as Render sets the environment instead, so each default reads it first.
func main() {
	addr := flag.String("addr", ":"+env("PORT", "8080"), "listen address")
	matchList := flag.String("match", env("FLOWSCORE_MATCH", demoMatches), "StatsBomb match IDs to replay, comma-separated")
	speed := flag.Float64("speed", 40, "1 = real time, 60 = one match minute per second")
	wait := flag.Duration("wait", 3*time.Minute, "how long a match is announced before kick-off")
	rest := flag.Duration("rest", 2*time.Minute, "how long a finished match stays on the list")
	dataDir := flag.String("data", env("FLOWSCORE_DATA", "data/statsbomb"), "folder with the matches, events-<id> and lineups-<id> files")
	matchesFile := flag.String("matches", env("FLOWSCORE_MATCHES", "matches-2-27.json"), "matches file in the data folder")
	modelPath := flag.String("model", env("FLOWSCORE_MODEL", "../ai/prediction/model.json"), "outcome model from ai/prediction/fit.py")
	xgPath := flag.String("xg-model", env("FLOWSCORE_XG_MODEL", "../ai/xg/model.json"), "our xG model from ai/xg/fit.py")
	originList := flag.String("origins", env("FLOWSCORE_ORIGINS", "localhost:5173,*:5173"), "host patterns of the sites allowed to call the API, comma-separated")
	llmModel := flag.String("llm-model", env("FLOWSCORE_LLM_MODEL", "gemini-3.1-flash-lite"), "Gemini model that writes the match analysis")
	llmGap := flag.Duration("llm-gap", 7*time.Second, "pause after each request to the LLM, to stay under its rate limit")
	insights := flag.String("insights", env("FLOWSCORE_INSIGHTS", "demo/insights.json"), "file that keeps the analysis texts between starts")
	espnLeagues := flag.String("espn", env("FLOWSCORE_ESPN", "eng.1"), "ESPN leagues whose real matches to follow, comma-separated; empty for none")
	espnEvery := flag.Duration("espn-every", 30*time.Second, "how often to read ESPN")

	flag.Parse()

	// Ctrl+C locally, SIGTERM when a host stops or redeploys the container
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	replay := statsbomb.Replay{
		MatchesPath: filepath.Join(*dataDir, *matchesFile),
		EventsDir:   *dataDir,
		Speed:       *speed,
	}
	model, err := predict.Load(*modelPath)
	if err != nil {
		log.Fatalf("outcome model: %v", err)
	}
	xgModel, err := xg.Load(*xgPath)
	if err != nil {
		log.Fatalf("xG model: %v", err)
	}

	agent, err := newAgent(ctx, os.Getenv("GEMINI_API_KEY"), *llmModel, *llmGap, *insights)
	if err != nil {
		log.Fatalf("analysis agent: %v", err)
	}

	hub := live.NewHub()
	store := live.NewStore()
	params := flow.DefaultParams()

	// Every match gets its own Publisher with the same models
	newPublisher := func(matchID string) *live.Publisher {
		p := live.NewPublisher(hub, store, matchID, params)
		p.UseModel(&model)
		p.UseXG(&xgModel)
		if agent != nil {
			p.UseAgent(agent)
		}
		return p
	}

	ids := strings.Split(*matchList, ",")
	// A match takes 95 minutes of match time; the rounds of the matches are
	// spread evenly over one cycle, so some are always live and some are next
	round := *wait + time.Duration(float64(95*time.Minute) / *speed) + *rest
	for i, raw := range ids {
		id, err := strconv.Atoi(strings.TrimSpace(raw))
		if err != nil {
			log.Fatalf("match ID %q: %v", raw, err)
		}
		info, err := statsbomb.LoadMatchInfo(replay.MatchesPath, id)
		if err != nil {
			log.Fatal(err)
		}
		publisher := newPublisher(info.ID)

		show := schedule{wait: *wait, rest: *rest, first: 20*time.Second + time.Duration(i)*round/time.Duration(len(ids))}
		go show.run(ctx, hub, publisher, params, info, &replay)
	}

	if *espnLeagues != "" {
		follow := liveESPN{
			client:       espn.NewClient(),
			leagues:      strings.Split(*espnLeagues, ","),
			every:        *espnEvery,
			hub:          hub,
			store:        store,
			params:       params,
			newPublisher: newPublisher,
		}
		go follow.run(ctx)
	}

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
		log.Printf("listening on %s, %d matches on the schedule", *addr, len(ids))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Print("shutting down")

	hub.Close() // WebSocket handlers see their channel closed and return
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}

// newAgent starts the agent that writes the match analysis. Without an API
// key it only serves the texts kept in the file; with neither it returns nil,
// and the match screen explains Flow on its own.
func newAgent(ctx context.Context, apiKey, model string, gap time.Duration, file string) (*insight.Agent, error) {
	var writer insight.Writer
	if apiKey != "" {
		g, err := insight.NewGemini(ctx, apiKey, model, "")
		if err != nil {
			return nil, err
		}
		writer = g
	}
	agent := insight.NewAgent(writer, gap)
	if err := agent.UseFile(file); err != nil {
		return nil, err
	}

	switch {
	case writer != nil:
		log.Printf("AI analysis: %d texts kept in %s, new ones by %s", agent.Len(), file, model)
	case agent.Len() > 0:
		log.Printf("AI analysis: %d texts kept in %s; GEMINI_API_KEY is not set, so no new ones", agent.Len(), file)
	default:
		log.Print("GEMINI_API_KEY is not set and no texts are kept: no AI analysis")
		return nil, nil
	}
	go agent.Run(ctx)
	return agent, nil
}

// env returns the environment variable, or fallback when it is not set.
func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// schedule plays one demo match round and round, the way a day of football
// looks on the home screen: announced with a kick-off time, live, finished.
type schedule struct {
	first time.Duration // the wait before the first kick-off, to stagger the matches
	wait  time.Duration
	rest  time.Duration
}

func (s schedule) run(ctx context.Context, hub *live.Hub, publisher *live.Publisher, params flow.Params, info event.Match, replay *statsbomb.Replay) {
	wait := s.first
	for {
		info.KickoffAt = time.Now().Add(wait).UTC()
		publisher.Schedule(info)
		// Viewers of the last round reload, so they see the new kick-off
		hub.CloseMatch(info.ID)
		if !sleep(ctx, wait) {
			return
		}

		// Kick-off goes into the store before the viewers are dropped, so the
		// snapshot they reload is already the live match
		publisher.Start(info)
		hub.CloseMatch(info.ID)
		if err := replayOnce(ctx, publisher, params, info, replay); err != nil {
			log.Printf("replay %s: %v", info.ID, err)
			return
		}
		if !sleep(ctx, s.rest) {
			return
		}
		wait = s.wait
	}
}

// sleep waits for d and reports false if the server is shutting down instead.
func sleep(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
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
	log.Printf("match %s finished", matchID)
	return nil
}
