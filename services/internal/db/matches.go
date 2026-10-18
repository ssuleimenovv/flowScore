package db

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ssuleimenovv/flowscore/services/internal/live"
)

// Matches keeps the finished matches in Postgres, each as its last snapshot.
type Matches struct {
	Pool *pgxpool.Pool
}

// Save keeps the match. A match saved again replaces what was there.
func (m *Matches) Save(ctx context.Context, s live.Snapshot) error {
	body, err := json.Marshal(s)
	if err != nil {
		return err
	}
	_, err = m.Pool.Exec(ctx, `
		INSERT INTO matches (id, source, competition, kickoff_at, snapshot)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET snapshot = EXCLUDED.snapshot, finished_at = now()`,
		s.Match.ID, s.Match.Source, s.Match.CompetitionID, s.Match.KickoffAt, body)
	return err
}

// Get is live.Archive: the saved match, or false. A database that does not
// answer in time counts as not having it, so a viewer gets a 404, not a hang.
func (m *Matches) Get(matchID string) (live.Snapshot, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var body []byte
	err := m.Pool.QueryRow(ctx, `SELECT snapshot FROM matches WHERE id = $1`, matchID).Scan(&body)
	if errors.Is(err, pgx.ErrNoRows) {
		return live.Snapshot{}, false
	}
	if err != nil {
		log.Printf("db: match %s: %v", matchID, err)
		return live.Snapshot{}, false
	}
	var s live.Snapshot
	if err := json.Unmarshal(body, &s); err != nil {
		log.Printf("db: match %s: %v", matchID, err)
		return live.Snapshot{}, false
	}
	return s, true
}
