package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/auth"
	"github.com/ssuleimenovv/flowscore/services/internal/db"
)

// Accounts is what signing up and in need from the database: db.Users
// does it, a test fakes it.
type Accounts interface {
	Create(ctx context.Context, email, name, hash string) (db.User, error)
	ByEmail(ctx context.Context, email string) (db.User, string, error)
	ByID(ctx context.Context, id string) (db.User, error)
}

// RegisterAuth adds signing up and signing in.
func RegisterAuth(mux *http.ServeMux, users Accounts, tokens auth.Tokens) {
	mux.HandleFunc("POST /api/v1/auth/register", registerHandler(users, tokens))
	mux.HandleFunc("POST /api/v1/auth/login", loginHandler(users, tokens))
	mux.HandleFunc("GET /api/v1/me", signedIn(tokens, meHandler(users)))
}

type user struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// session is the answer to both: the token to send from now on, and who it is
type session struct {
	Token string `json:"token"`
	User  user   `json:"user"`
}

type credentials struct {
	Name     string `json:"name"` // signing up only
	Email    string `json:"email"`
	Password string `json:"password"`
}

func registerHandler(users Accounts, tokens auth.Tokens) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, ok := readCredentials(w, r)
		if !ok {
			return
		}
		name := strings.TrimSpace(in.Name)
		if name == "" || !strings.Contains(in.Email, "@") {
			writeProblem(w, http.StatusUnprocessableEntity, "Name and email are required")
			return
		}
		if err := auth.CheckPassword(in.Password); err != nil {
			writeProblem(w, http.StatusUnprocessableEntity, "Weak password")
			return
		}
		hash, err := auth.HashPassword(in.Password)
		if err != nil {
			serverError(w, "hash", err)
			return
		}
		u, err := users.Create(r.Context(), in.Email, name, hash)
		if errors.Is(err, db.ErrEmailTaken) {
			writeProblem(w, http.StatusConflict, "Email already registered")
			return
		}
		if err != nil {
			serverError(w, "sign up", err)
			return
		}
		writeSession(w, http.StatusCreated, u, tokens)
	}
}

// dummyHash is checked against when no one has the email, so an unknown
// email takes as long to answer as a wrong password: the timing does not
// tell which emails are signed up
var dummyHash, _ = auth.HashPassword("not anyone's password 0")

func loginHandler(users Accounts, tokens auth.Tokens) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in, ok := readCredentials(w, r)
		if !ok {
			return
		}
		u, hash, err := users.ByEmail(r.Context(), in.Email)
		if errors.Is(err, db.ErrNoUser) {
			auth.SamePassword(dummyHash, in.Password)
		} else if err != nil {
			serverError(w, "sign in", err)
			return
		}
		// One answer for both, so it does not tell which one was wrong
		if err != nil || !auth.SamePassword(hash, in.Password) {
			writeProblem(w, http.StatusUnauthorized, "Wrong email or password")
			return
		}
		writeSession(w, http.StatusOK, u, tokens)
	}
}

func readCredentials(w http.ResponseWriter, r *http.Request) (credentials, bool) {
	var in credentials
	body := http.MaxBytesReader(w, r.Body, 4<<10)
	if err := json.NewDecoder(body).Decode(&in); err != nil {
		writeProblem(w, http.StatusBadRequest, "Invalid request")
		return credentials{}, false
	}
	in.Email = auth.NormalizeEmail(in.Email)
	return in, true
}

func writeSession(w http.ResponseWriter, status int, u db.User, tokens auth.Tokens) {
	token, err := tokens.Issue(u.ID, time.Now())
	if err != nil {
		serverError(w, "token", err)
		return
	}
	writeJSON(w, status, session{Token: token, User: user{ID: u.ID, Email: u.Email, Name: u.Name}})
}

// serverError logs what went wrong and tells the client only that it did
func serverError(w http.ResponseWriter, what string, err error) {
	log.Printf("auth: %s: %v", what, err)
	writeProblem(w, http.StatusInternalServerError, "Something went wrong")
}

// signedIn lets through a request with a valid access token, telling next
// whose it is; any other request gets 401, and the client signs in again.
func signedIn(tokens auth.Tokens, next func(w http.ResponseWriter, r *http.Request, userID string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			writeProblem(w, http.StatusUnauthorized, "Sign in first")
			return
		}
		id, err := tokens.UserID(token)
		if err != nil {
			writeProblem(w, http.StatusUnauthorized, "Session expired")
			return
		}
		next(w, r, id)
	}
}

func meHandler(users Accounts) func(http.ResponseWriter, *http.Request, string) {
	return func(w http.ResponseWriter, r *http.Request, userID string) {
		u, err := users.ByID(r.Context(), userID)
		// A token outlives an account deleted after it was issued
		if errors.Is(err, db.ErrNoUser) {
			writeProblem(w, http.StatusUnauthorized, "Sign in first")
			return
		}
		if err != nil {
			serverError(w, "me", err)
			return
		}
		writeJSON(w, http.StatusOK, user{ID: u.ID, Email: u.Email, Name: u.Name})
	}
}
