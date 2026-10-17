package db

import (
	"context"
	"os"
	"testing"
)

// A real Postgres, as FLOWSCORE_TEST_DATABASE_URL names it; without one the
// test is skipped
func TestOpenMigratesOnce(t *testing.T) {
	url := os.Getenv("FLOWSCORE_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("FLOWSCORE_TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()

	// Twice: the second start finds every migration run and changes nothing
	for range 2 {
		pool, err := Open(ctx, url)
		if err != nil {
			t.Fatal(err)
		}
		pool.Close()
	}

	pool, err := Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	var runs int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM schema_migrations WHERE name = 'migrations/001_matches.sql'`).Scan(&runs); err != nil {
		t.Fatal(err)
	}
	if runs != 1 {
		t.Errorf("001 recorded %d times", runs)
	}
	if _, err := pool.Exec(ctx, `SELECT id, snapshot FROM matches LIMIT 1`); err != nil {
		t.Errorf("matches table: %v", err)
	}
}
