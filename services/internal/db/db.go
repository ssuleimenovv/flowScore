package db

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

// The schema, one file per change, run in name order: 001_, 002_, ...
// A file once run is never edited; a change is a new file.
//
//go:embed migrations/*.sql
var migrations embed.FS

// Open connects to Postgres and brings its schema up to date.
func Open(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: %w", err)
	}
	if err := migrate(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

// migrate runs every migration not run yet, each in its own transaction:
// the file and its row in schema_migrations go in together or not at all.
func migrate(ctx context.Context, pool *pgxpool.Pool) error {
	if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		name       text PRIMARY KEY,
		applied_at timestamptz NOT NULL DEFAULT now()
	)`); err != nil {
		return err
	}
	names, err := fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)
	for _, name := range names {
		if err := run(ctx, pool, name); err != nil {
			return fmt.Errorf("migration %s: %w", name, err)
		}
	}
	return nil
}

func run(ctx context.Context, pool *pgxpool.Pool, name string) error {
	sql, err := migrations.ReadFile(name)
	if err != nil {
		return err
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // does nothing after Commit

	// The row first: a second server starting at the same time waits on it
	// and then skips the file
	tag, err := tx.Exec(ctx, `INSERT INTO schema_migrations (name) VALUES ($1) ON CONFLICT DO NOTHING`, name)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return nil // run before
	}
	// No arguments: pgx sends it as a simple query, so a file may hold
	// several statements
	if _, err := tx.Exec(ctx, string(sql)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
