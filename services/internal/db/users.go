package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User is one who signed up.
type User struct {
	ID        string
	Email     string
	Name      string
	Verified  bool // the email is confirmed by a code
	CreatedAt time.Time
}

var (
	// ErrEmailTaken: a sign-up with an email someone already signed up with
	ErrEmailTaken = errors.New("email taken")
	// ErrNoUser: no one signed up with that email
	ErrNoUser = errors.New("no such user")
)

// Users keeps the users in Postgres.
type Users struct {
	Pool *pgxpool.Pool
}

// Create signs a new user up. hash is the password's bcrypt hash: the
// password itself never reaches the database.
func (u *Users) Create(ctx context.Context, email, name, hash string) (User, error) {
	user := User{Email: email, Name: name}
	err := u.Pool.QueryRow(ctx, `
		INSERT INTO users (email, name, password_hash) VALUES ($1, $2, $3)
		RETURNING id::text, created_at`, email, name, hash).Scan(&user.ID, &user.CreatedAt)
	// The UNIQUE on email answers for us, with no race between a check and
	// the insert
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
		return User{}, ErrEmailTaken
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// ByEmail finds a user and their password hash, for signing in.
func (u *Users) ByEmail(ctx context.Context, email string) (User, string, error) {
	var user User
	var hash string
	err := u.Pool.QueryRow(ctx, `
		SELECT id::text, email, name, created_at, email_verified_at IS NOT NULL, password_hash
		FROM users WHERE email = $1`,
		email).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.Verified, &hash)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, "", ErrNoUser
	}
	if err != nil {
		return User{}, "", err
	}
	return user, hash, nil
}

// ByID finds a user by the ID an access token names.
func (u *Users) ByID(ctx context.Context, id string) (User, error) {
	user := User{ID: id}
	err := u.Pool.QueryRow(ctx, `
		SELECT email, name, created_at, email_verified_at IS NOT NULL FROM users WHERE id = $1`,
		id).Scan(&user.Email, &user.Name, &user.CreatedAt, &user.Verified)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNoUser
	}
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// MarkVerified records that the user's email is confirmed.
func (u *Users) MarkVerified(ctx context.Context, id string) error {
	_, err := u.Pool.Exec(ctx, `
		UPDATE users SET email_verified_at = now() WHERE id = $1 AND email_verified_at IS NULL`, id)
	return err
}
