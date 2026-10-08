package main

import (
	"context"
	"log"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
)

// liveESPN follows the real matches of some leagues on ESPN. A match goes on
// the list when the scoreboard has it, is played from its summary once it
// kicks off and comes off the list when the scoreboard drops it.
type liveESPN struct {
	client       *espn.Client
	leagues      []string
	every        time.Duration // how often the scoreboard and each summary are read
	hub          *live.Hub
	store        *live.Store
	params       flow.Params
	newPublisher func(matchID string) *live.Publisher

	games map[string]*liveGame // by match ID; only run's goroutine touches it
}

type liveGame struct {
	league    string
	info      event.Match
	publisher *live.Publisher
	started   bool
	done      chan struct{} // closed when the match has been played out
}

func (l *liveESPN) run(ctx context.Context) {
	l.games = map[string]*liveGame{}
	for {
		for _, league := range l.leagues {
			l.poll(ctx, league)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(l.every):
		}
	}
}

// poll reads one league's scoreboard and follows every match on it. A failed
// read changes nothing, so a hiccup of the API does not empty the list.
func (l *liveESPN) poll(ctx context.Context, league string) {
	games, err := l.client.Scoreboard(ctx, league)
	if err != nil {
		log.Printf("espn %s: %v", league, err)
		return
	}
	seen := map[string]bool{}
	for _, g := range games {
		seen[g.ID] = true
		l.follow(ctx, league, g)
	}
	// A match the scoreboard dropped is from an earlier day
	for id, g := range l.games {
		if g.league == league && !seen[id] && !playing(g) {
			l.store.Remove(id)
			l.hub.CloseMatch(id)
			delete(l.games, id)
		}
	}
}

func (l *liveESPN) follow(ctx context.Context, league string, g espn.Game) {
	game, ok := l.games[g.ID]
	if !ok {
		if g.State == "post" {
			return // over before the gateway saw it: nothing to play
		}
		body, err := l.client.Summary(ctx, league, g.ID)
		if err != nil {
			log.Printf("espn %s: %v", g.Name, err)
			return
		}
		s, err := espn.ParseSummary(body)
		if err != nil {
			log.Printf("espn %s: %v", g.Name, err)
			return
		}
		game = &liveGame{league: league, info: s.Match, publisher: l.newPublisher(g.ID), done: make(chan struct{})}
		game.publisher.Schedule(game.info)
		l.games[g.ID] = game
		log.Printf("espn: %s, kick-off %s", g.Name, g.Start.Format(time.RFC3339))
	}
	if g.State != "pre" && !game.started {
		game.started = true
		go l.play(ctx, game)
	}
}

// play runs the match through the Flow engine in real time until the final
// whistle, the way schedule.run plays a demo match.
func (l *liveESPN) play(ctx context.Context, g *liveGame) {
	defer close(g.done)
	// Kick-off goes into the store before the viewers are dropped, so the
	// snapshot they reload is already the live match
	g.publisher.Start(g.info)
	l.hub.CloseMatch(g.info.ID)

	source := espn.Live{Client: l.client, League: g.league, Every: l.every}
	events, err := source.Stream(ctx, g.info.ID)
	if err != nil {
		log.Printf("espn %s: %v", g.info.ID, err)
		return
	}
	engine := flow.Engine{Params: l.params, Speed: 1, Tick: 5 * time.Second}
	log.Printf("espn: %s — %s kicked off", g.info.Home.Name, g.info.Away.Name)
	g.publisher.Run(engine.Run(ctx, g.info.ID, events))
	log.Printf("espn: %s — %s finished", g.info.Home.Name, g.info.Away.Name)
}

// playing reports whether the match kicked off and is not over yet.
func playing(g *liveGame) bool {
	if !g.started {
		return false
	}
	select {
	case <-g.done:
		return false
	default:
		return true
	}
}
