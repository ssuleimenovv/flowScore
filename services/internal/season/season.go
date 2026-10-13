// Package season keeps what the current season of each league has come to:
// its table, all its matches by round, and the Flow of every match played,
// which it gets by running the match through the Flow engine once it is over.
package season

import (
	"context"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
	"github.com/ssuleimenovv/flowscore/services/internal/xg"
)

// League is one league's season as far as it has been read.
type League struct {
	ID        string // "eng.1"
	Table     espn.Table
	Games     []Game // the whole season, by kick-off
	UpdatedAt time.Time
}

// Game is one match of the season.
type Game struct {
	espn.Game
	Round int        // 1 for the first round; worked out, ESPN does not say
	Flow  *MatchFlow // nil until the match is over and read
}

// Seasons reads the seasons of some leagues from ESPN and keeps them. Its
// methods are safe to call from many goroutines.
type Seasons struct {
	Client *espn.Client
	Params flow.Params
	XG     *xg.Model // nil leaves xG at 0

	mu      sync.RWMutex
	leagues map[string]League
	flows   map[string]MatchFlow // by match ID: a match over stays over
}

// League returns a league's season, and false before it has been read.
func (s *Seasons) League(id string) (League, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.leagues[id]
	return l, ok
}

// Run reads every league now and again every every, until ctx is done.
func (s *Seasons) Run(ctx context.Context, leagues []string, every time.Duration) {
	for {
		for _, id := range leagues {
			start := time.Now()
			if err := s.Refresh(ctx, id); err != nil {
				log.Printf("season %s: %v", id, err)
			} else {
				log.Printf("season %s: read in %s", id, time.Since(start).Round(time.Second))
			}
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(every):
		}
	}
}

// Refresh reads a league again: its table, every day of its calendar and
// the matches over since the last time. A match that cannot be read is
// left for the next time.
func (s *Seasons) Refresh(ctx context.Context, id string) error {
	table, err := s.Client.Standings(ctx, id)
	if err != nil {
		return err
	}
	today, err := s.Client.Board(ctx, id, time.Time{})
	if err != nil {
		return err
	}

	seen := map[string]bool{}
	var games []Game
	for _, day := range today.Days {
		b, err := s.Client.Board(ctx, id, day)
		if err != nil {
			return err
		}
		for _, g := range b.Games {
			// A day of ESPN's runs into the next one: a late match is on both
			if !seen[g.ID] {
				seen[g.ID] = true
				games = append(games, Game{Game: g})
			}
		}
	}

	for i := range games {
		g := &games[i]
		if g.State != "post" {
			continue
		}
		f, ok := s.flow(g.ID)
		if !ok {
			if f, err = s.read(ctx, id, g.ID); err != nil {
				log.Printf("season %s, match %s: %v", id, g.ID, err)
				continue
			}
			s.mu.Lock()
			s.flows[g.ID] = f
			s.mu.Unlock()
		}
		g.Flow = &f
	}
	slices.SortStableFunc(games, func(a, b Game) int { return a.Start.Compare(b.Start) })
	numberRounds(games)

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.leagues == nil {
		s.leagues = map[string]League{}
	}
	s.leagues[id] = League{ID: id, Table: table, Games: games, UpdatedAt: time.Now().UTC()}
	return nil
}

func (s *Seasons) flow(matchID string) (MatchFlow, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.flows == nil {
		s.flows = map[string]MatchFlow{}
	}
	f, ok := s.flows[matchID]
	return f, ok
}

// read reads a finished match and runs it through the engine.
func (s *Seasons) read(ctx context.Context, league, matchID string) (MatchFlow, error) {
	body, err := s.Client.Summary(ctx, league, matchID)
	if err != nil {
		return MatchFlow{}, err
	}
	_, events, err := espn.Played(body)
	if err != nil {
		return MatchFlow{}, err
	}
	return Measure(events, s.Params, s.XG), nil
}

// numberRounds gives each match its round: one more than the last round
// either team has played. A match put off to later gets a later round,
// which is where it is played.
func numberRounds(games []Game) {
	last := map[string]int{} // by team ID
	for i := range games {
		g := &games[i]
		g.Round = max(last[g.Home.Team.ID], last[g.Away.Team.ID]) + 1
		last[g.Home.Team.ID], last[g.Away.Team.ID] = g.Round, g.Round
	}
}
