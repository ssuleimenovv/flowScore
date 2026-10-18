package main

import (
	"context"
	"log"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/db"
	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
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
	archive      *db.Matches // nil without a database

	games map[string]*liveGame // by match ID; only run's goroutine touches it
}

type liveGame struct {
	league    string
	info      event.Match
	publisher *live.Publisher
	started   bool
	lineups   bool          // ESPN has put up the lineups
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
		l.games[g.ID] = game
		l.readSummary(game, body)
		log.Printf("espn: %s, kick-off %s", g.Name, g.Start.Format(time.RFC3339))
	}
	// Until kick-off the chances follow the season's ratings, which are read
	// a few minutes after a restart and again every few hours
	if g.State == "pre" {
		if m := g.Market; m != nil {
			game.publisher.UseMarket(predict.Outcome{Home: m.Home, Draw: m.Draw, Away: m.Away})
		}

		if !game.lineups && time.Until(g.Start) < lineupsAhead {
			if body, err := l.client.Summary(ctx, league, g.ID); err == nil {
				l.readSummary(game, body)
			}
		}

		game.publisher.Schedule(game.info)
	}
	if g.State != "pre" && !game.started {
		game.started = true
		go l.play(ctx, game)
	}
}

// ESPN puts the lineups up about an hour before kick-off. From this long
// before, the summary of a match still without them is read on every poll
const lineupsAhead = 90 * time.Minute

// readSummary gives the match what a summary read before kick-off has: the
// teams' form, meetings and news, and the lineups once ESPN has them.
func (l *liveESPN) readSummary(game *liveGame, body []byte) {
	if pm, err := espn.ParsePreMatch(body); err != nil {
		log.Printf("espn %s: %v", game.info.ID, err)
	} else if pm != nil {
		game.publisher.UsePreMatch(pm)
	}
	lineups, err := espn.ParseRosters(body)
	if err != nil {
		log.Printf("espn %s: %v", game.info.ID, err)
		return
	}
	if lineups != nil {
		game.publisher.UseLineups(lineups)
		game.lineups = true
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

	// The lineups come with every summary, with the substitutions so far
	source := espn.Live{Client: l.client, League: g.league, Every: l.every, Lineups: g.publisher.UseLineups}

	events, err := source.Stream(ctx, g.info.ID)
	if err != nil {
		log.Printf("espn %s: %v", g.info.ID, err)
		return
	}
	engine := flow.Engine{Params: l.params, Speed: 1, Tick: 5 * time.Second}
	log.Printf("espn: %s — %s kicked off", g.info.Home.Name, g.info.Away.Name)
	g.publisher.Run(engine.Run(ctx, g.info.ID, events))
	log.Printf("espn: %s — %s finished", g.info.Home.Name, g.info.Away.Name)
	// The match outlives the server: it opens after a restart too. Run also
	// returns when the server stops, so only a match that is over is kept
	if l.archive != nil {
		if snap, ok := l.store.Get(g.info.ID); ok && snap.Status == "finished" {
			if err := l.archive.Save(ctx, snap); err != nil {
				log.Printf("espn %s: archive: %v", g.info.ID, err)
			}
		}
	}
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
