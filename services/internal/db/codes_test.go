package db

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"
)

func TestCodes(t *testing.T) {
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
	t.Cleanup(func() { pool.Exec(ctx, `DELETE FROM users WHERE email = 'codes@flowscore.dev'`) })

	users := Users{Pool: pool}
	u, err := users.Create(ctx, "codes@flowscore.dev", "Тест", "hash")
	if err != nil {
		t.Fatal(err)
	}
	if u.Verified {
		t.Error("a new user is verified before any code")
	}
	codes := Codes{Pool: pool}

	// A code, and "send again" right away is turned down
	if err := codes.Put(ctx, u.ID, Verify, "right", 10*time.Minute, time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := codes.Put(ctx, u.ID, Verify, "other", 10*time.Minute, time.Minute); !errors.Is(err, ErrTooSoon) {
		t.Errorf("second code at once: %v, want ErrTooSoon", err)
	}
	// A wrong try uses one up; the right code works once
	if left, err := codes.Use(ctx, u.ID, Verify, "wrong"); !errors.Is(err, ErrWrongCode) || left != MaxAttempts-1 {
		t.Errorf("wrong code: %d left, %v", left, err)
	}
	if _, err := codes.Use(ctx, u.ID, Verify, "right"); err != nil {
		t.Errorf("right code: %v", err)
	}
	if _, err := codes.Use(ctx, u.ID, Verify, "right"); !errors.Is(err, ErrNoCode) {
		t.Errorf("right code again: %v, want ErrNoCode", err)
	}

	// An expired code is no code
	if err := codes.Put(ctx, u.ID, Reset, "right", -time.Second, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := codes.Use(ctx, u.ID, Reset, "right"); !errors.Is(err, ErrNoCode) {
		t.Errorf("expired code: %v, want ErrNoCode", err)
	}

	// Out of tries, even the right code is turned down
	if err := codes.Put(ctx, u.ID, Reset, "right", 10*time.Minute, 0); err != nil {
		t.Fatal(err)
	}
	for range MaxAttempts {
		codes.Use(ctx, u.ID, Reset, "wrong")
	}
	if _, err := codes.Use(ctx, u.ID, Reset, "right"); !errors.Is(err, ErrNoCode) {
		t.Errorf("right code after %d wrong: %v, want ErrNoCode", MaxAttempts, err)
	}

	if err := users.MarkVerified(ctx, u.ID); err != nil {
		t.Fatal(err)
	}
	if got, _ := users.ByID(ctx, u.ID); !got.Verified {
		t.Error("not verified after MarkVerified")
	}
}
