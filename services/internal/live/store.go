package live

import (
	"slices"
	"sync"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Snapshot is the state of one match that the REST API returns.
// A client loads it first and then applies stream messages with a greater Seq.
type Snapshot struct {
	Match      event.Match
	Status     string // scheduled, live, halftime, finished
	Score      Score
	Halftime   *Score // score at the break, nil before it
	Period     int
	At         time.Duration // match time of the last update
	Flow       FlowValues
	Delta10    FlowValues
	Factors    []FlowFactor
	Points     []FlowPoint  // one per match minute, in order
	Events     []MatchEvent // oldest first
	Stats      []StatRow    // empty before kick-off
	Prediction *Prediction  // nil when no outcome model is loaded
	Outlook    *Outlook     // the same, for the simulator
	// the Latest analysis, nil until the agent answers. A new one replaces
	// it whole, so a copy may share it
	Explanation *Explanation
	Seq         int64 // last stream message already included
	UpdatedAt   time.Time
}

type Score struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

// Store keeps the latest snapshot of every match. The Publisher writes while
// HTTP handlers read at the same time, so access goes through an RWMutex.
type Store struct {
	mu      sync.RWMutex
	matches map[string]*Snapshot
}

func NewStore() *Store {
	return &Store{matches: map[string]*Snapshot{}}
}

// Schedule registers a match that has not started yet. seq works as in Start:
// a demo match is announced again after every round, and seq keeps growing.
func (s *Store) Schedule(m event.Match, seq int64) {
	s.reset(m, "scheduled", seq)
}

// Start resets the match to kick-off. seq is the last message already sent:
// the fresh snapshot counts as including it, so older messages are skipped.
func (s *Store) Start(m event.Match, seq int64) {
	s.reset(m, "live", seq)
}

func (s *Store) reset(m event.Match, status string, seq int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.matches[m.ID] = &Snapshot{Match: m, Status: status, Period: 1, Seq: seq, UpdatedAt: time.Now().UTC()}
}

// Get returns a copy of the snapshot. The caller may keep reading it
// after the lock is released while the Publisher goes on appending.
func (s *Store) Get(matchID string) (Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snap, ok := s.matches[matchID]
	if !ok {
		return Snapshot{}, false
	}
	return snap.clone(), true
}

// List returns a copy of every match, by kick-off time.
func (s *Store) List() []Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Snapshot, 0, len(s.matches))
	for _, snap := range s.matches {
		out = append(out, snap.clone())
	}
	slices.SortFunc(out, func(a, b Snapshot) int {
		return a.Match.KickoffAt.Compare(b.Match.KickoffAt)
	})
	return out
}

// clone copies the snapshot deep enough that the Publisher can go on
// appending to the original. Call it under the store's lock.
func (snap *Snapshot) clone() Snapshot {
	c := *snap
	c.Points = slices.Clone(snap.Points)
	c.Events = slices.Clone(snap.Events)
	c.Stats = slices.Clone(snap.Stats)
	c.Factors = slices.Clone(snap.Factors)

	if snap.Prediction != nil {
		p := *snap.Prediction // the Publisher updates it in place
		c.Prediction = &p
	}
	if snap.Outlook != nil {
		o := *snap.Outlook
		c.Outlook = &o
	}
	return c
}

func (s *Store) update(matchID string, fn func(*Snapshot)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if snap, ok := s.matches[matchID]; ok {
		fn(snap)
		snap.UpdatedAt = time.Now().UTC()
	}
}

// upsertPoint replaces the point of the same minute or appends a new one.
// Minutes can repeat: 45+1 of the first half and the start of the second
// half both count as minute 46.
func upsertPoint(points []FlowPoint, p FlowPoint) []FlowPoint {
	if i := slices.IndexFunc(points, func(q FlowPoint) bool { return q.Minute == p.Minute }); i >= 0 {
		points[i] = p
		return points
	}
	return append(points, p)
}
