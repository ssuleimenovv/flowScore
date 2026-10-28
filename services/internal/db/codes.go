package db

import (
	"context"
	"crypto/subtle"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// What a code is for
const (
	Verify = "verify" // confirm the email after signing up
	Reset  = "reset"  // set a new password
)

// MaxAttempts is how many wrong tries a code takes before it is gone
const MaxAttempts = 5

var (
	// ErrTooSoon: a code went out less than the gap ago
	ErrTooSoon = errors.New("code sent too recently")
	// ErrNoCode: no code to check, or it expired, or its tries ran out
	ErrNoCode = errors.New("no live code")
	// ErrWrongCode: the code does not match; one try is used up
	ErrWrongCode = errors.New("wrong code")
)

// Codes keeps the codes sent by email.
type Codes struct {
	Pool *pgxpool.Pool
}

// Put keeps a new code for the user, replacing the one before. When that
// one went out less than gap ago nothing changes and the answer is
// ErrTooSoon: someone pressing "send again" cannot flood a mailbox
func (c *Codes) Put(ctx context.Context, userID, purpose, hash string, ttl, gap time.Duration) error {
	tag, err := c.Pool.Exec(ctx, `
		INSERT INTO email_codes (user_id, purpose, code_hash, expires_at)
		VALUES ($1, $2, $3, now() + make_interval(secs => $4))
		ON CONFLICT (user_id, purpose) DO UPDATE
		SET code_hash = EXCLUDED.code_hash, expires_at = EXCLUDED.expires_at, attempts = 0, sent_at = now()
		WHERE email_codes.sent_at <= now() - make_interval(secs => $5)`,
		userID, purpose, hash, ttl.Seconds(), gap.Seconds())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrTooSoon
	}
	return nil
}

// Use checks a code. The right one works once and is gone; a wrong one uses
// up a try and tells how many are left. The row is locked while it is
// checked, so two tries at once cannot both get through
func (c *Codes) Use(ctx context.Context, userID, purpose, hash string) (left int, err error) {
	tx, err := c.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx) // does nothing after Commit

	var kept string
	var attempts int
	err = tx.QueryRow(ctx, `
		SELECT code_hash, attempts FROM email_codes
		WHERE user_id = $1 AND purpose = $2 AND expires_at > now() AND attempts < $3
		FOR UPDATE`, userID, purpose, MaxAttempts).Scan(&kept, &attempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, ErrNoCode
	}
	if err != nil {
		return 0, err
	}

	// Compared in constant time: how long it takes says nothing of how close a guess was
	if subtle.ConstantTimeCompare([]byte(kept), []byte(hash)) == 1 {
		if _, err := tx.Exec(ctx, `DELETE FROM email_codes WHERE user_id = $1 AND purpose = $2`, userID, purpose); err != nil {
			return 0, err
		}
		return 0, tx.Commit(ctx)
	}
	if _, err := tx.Exec(ctx, `UPDATE email_codes SET attempts = attempts + 1 WHERE user_id = $1 AND purpose = $2`, userID, purpose); err != nil {
		return 0, err
	}
	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}
	return MaxAttempts - attempts - 1, ErrWrongCode
}
