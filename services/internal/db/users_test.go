package db

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestUsers(t *testing.T) {
	url := os.Getenv("FLOWSCORE_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("FLOWSCORE_TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	t.Cleanup(func() { pool.Exec(ctx, `DELETE FROM users WHERE email = 'test@flowscore.dev'`) })

	u := Users{Pool: pool}
	made, err := u.Create(ctx, "test@flowscore.dev", "Тест", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if made.ID == "" || made.CreatedAt.IsZero() {
		t.Errorf("created %+v", made)
	}
	// One email, one user
	if _, err := u.Create(ctx, "test@flowscore.dev", "Другой", "hash"); !errors.Is(err, ErrEmailTaken) {
		t.Errorf("second sign-up: %v, want ErrEmailTaken", err)
	}

	got, hash, err := u.ByEmail(ctx, "test@flowscore.dev")
	if err != nil || got.ID != made.ID || got.Name != "Тест" || hash != "hash" {
		t.Errorf("by email: %+v, %q, %v", got, hash, err)
	}
	if _, _, err := u.ByEmail(ctx, "nobody@flowscore.dev"); !errors.Is(err, ErrNoUser) {
		t.Errorf("unknown email: %v, want ErrNoUser", err)
	}
	if got, err := u.ByID(ctx, made.ID); err != nil || got.Email != "test@flowscore.dev" {
		t.Errorf("by ID: %+v, %v", got, err)
	}
	// A well-formed ID no one has
	if _, err := u.ByID(ctx, "00000000-0000-0000-0000-000000000000"); !errors.Is(err, ErrNoUser) {
		t.Errorf("unknown ID: %v, want ErrNoUser", err)
	}
}
