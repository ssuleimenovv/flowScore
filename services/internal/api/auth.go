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
	"github.com/ssuleimenovv/flowscore/services/internal/mail"
)

// Accounts is what signing up and in need from the users table: db.Users
// does it, a test fakes it.
type Accounts interface {
	Create(ctx context.Context, email, name, hash string) (db.User, error)
	ByEmail(ctx context.Context, email string) (db.User, string, error)
	ByID(ctx context.Context, id string) (db.User, error)
	MarkVerified(ctx context.Context, id string) error
}

// CodeStore keeps the codes sent by email: db.Codes.
type CodeStore interface {
	Put(ctx context.Context, userID, purpose, hash string, ttl, gap time.Duration) error
	Use(ctx context.Context, userID, purpose, hash string) (int, error)
}

// Auth is signing up and in: the users, the email codes and the way they
// are sent, and the access tokens.
type Auth struct {
	Users  Accounts
	Codes  CodeStore
	Hasher auth.Codes
	Mail   mail.Sender
	Tokens auth.Tokens
}

// A code counts for 10 minutes; a new one goes out at most once a minute
const (
	codeTTL = 10 * time.Minute
	codeGap = time.Minute
)

// Register adds the auth routes of api/openapi.yaml.
func (a Auth) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/register", a.register)
	mux.HandleFunc("POST /api/v1/auth/login", a.login)
	mux.HandleFunc("POST /api/v1/auth/verify", a.verify)
	mux.HandleFunc("POST /api/v1/auth/resend", a.resend)
	mux.HandleFunc("GET /api/v1/me", signedIn(a.Tokens, a.me))
}

type user struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// session is what a confirmed sign-in gets: the token to send from now on,
// and who it is
type session struct {
	Token string `json:"token"`
	User  user   `json:"user"`
}

// pending is the answer while a code is on its way: the email it went to
type pending struct {
	Email string `json:"email"`
}

type credentials struct {
	Name     string `json:"name"` // signing up only
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"` // confirming only
}

func (a Auth) register(w http.ResponseWriter, r *http.Request) {
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
	u, err := a.Users.Create(r.Context(), in.Email, name, hash)
	if errors.Is(err, db.ErrEmailTaken) {
		writeProblem(w, http.StatusConflict, "Email already registered")
		return
	}
	if err != nil {
		serverError(w, "sign up", err)
		return
	}
	// No token yet: the account waits for the code from the email
	if err := a.sendCode(r.Context(), u, db.Verify); err != nil {
		mailError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, pending{Email: u.Email})
}

// dummyHash is checked against when no one has the email, so an unknown
// email takes as long to answer as a wrong password: the timing does not
// tell which emails are signed up
var dummyHash, _ = auth.HashPassword("not anyone's password 0")

func (a Auth) login(w http.ResponseWriter, r *http.Request) {
	in, ok := readCredentials(w, r)
	if !ok {
		return
	}
	u, hash, err := a.Users.ByEmail(r.Context(), in.Email)
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
	// The right password but an email not confirmed yet: a new code, unless
	// one went out a moment ago, and the client opens the code screen
	if !u.Verified {
		if err := a.sendCode(r.Context(), u, db.Verify); err != nil && !errors.Is(err, db.ErrTooSoon) {
			mailError(w, err)
			return
		}
		writeProblem(w, http.StatusForbidden, "Email not verified")
		return
	}
	a.writeSession(w, http.StatusOK, u)
}

// verify checks the code from the email. The right one confirms the email
// and signs in
func (a Auth) verify(w http.ResponseWriter, r *http.Request) {
	in, ok := readCredentials(w, r)
	if !ok {
		return
	}
	u, _, err := a.Users.ByEmail(r.Context(), in.Email)
	// An unknown email gets what a lost code gets
	if errors.Is(err, db.ErrNoUser) {
		writeProblem(w, http.StatusGone, "No code to check")
		return
	}
	if err != nil {
		serverError(w, "verify", err)
		return
	}
	code := strings.TrimSpace(in.Code)
	left, err := a.Codes.Use(r.Context(), u.ID, db.Verify, a.Hasher.Hash(u.ID, db.Verify, code))
	switch {
	case errors.Is(err, db.ErrWrongCode):
		writeWrongCode(w, left)
		return
	case errors.Is(err, db.ErrNoCode):
		writeProblem(w, http.StatusGone, "No code to check")
		return
	case err != nil:
		serverError(w, "verify", err)
		return
	}
	if err := a.Users.MarkVerified(r.Context(), u.ID); err != nil {
		serverError(w, "verify", err)
		return
	}
	a.writeSession(w, http.StatusOK, u)
}

// resend sends a new code. It answers an unknown or confirmed email the same,
// so it does not tell which emails are signed up
func (a Auth) resend(w http.ResponseWriter, r *http.Request) {
	in, ok := readCredentials(w, r)
	if !ok {
		return
	}
	u, _, err := a.Users.ByEmail(r.Context(), in.Email)
	if err != nil && !errors.Is(err, db.ErrNoUser) {
		serverError(w, "resend", err)
		return
	}
	if err == nil && !u.Verified {
		err := a.sendCode(r.Context(), u, db.Verify)
		if errors.Is(err, db.ErrTooSoon) {
			w.Header().Set("Retry-After", "60")
			writeProblem(w, http.StatusTooManyRequests, "Wait a minute before another code")
			return
		}
		if err != nil {
			mailError(w, err)
			return
		}
	}
	writeJSON(w, http.StatusAccepted, pending{Email: in.Email})
}

func (a Auth) me(w http.ResponseWriter, r *http.Request, userID string) {
	u, err := a.Users.ByID(r.Context(), userID)
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

// sendCode makes a code, keeps its hash and emails the code itself.
// db.ErrTooSoon means one went out less than a minute ago
func (a Auth) sendCode(ctx context.Context, u db.User, purpose string) error {
	code, err := auth.NewCode()
	if err != nil {
		return err
	}
	if err := a.Codes.Put(ctx, u.ID, purpose, a.Hasher.Hash(u.ID, purpose, code), codeTTL, codeGap); err != nil {
		return err
	}
	return a.Mail.Send(ctx, mail.CodeMail(u.Email, code, purpose == db.Reset))
}

func (a Auth) writeSession(w http.ResponseWriter, status int, u db.User) {
	token, err := a.Tokens.Issue(u.ID, time.Now())
	if err != nil {
		serverError(w, "token", err)
		return
	}
	writeJSON(w, status, session{Token: token, User: user{ID: u.ID, Email: u.Email, Name: u.Name}})
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

// wrongCode is the problem a wrong code gets: how many tries are left
type wrongCode struct {
	problem
	AttemptsLeft int `json:"attemptsLeft"`
}

func writeWrongCode(w http.ResponseWriter, left int) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	p := wrongCode{problem{Type: "about:blank", Title: "Wrong code", Status: http.StatusUnprocessableEntity}, left}
	if err := json.NewEncoder(w).Encode(p); err != nil {
		log.Printf("write problem: %v", err)
	}
}

// serverError logs what went wrong and tells the client only that it did
func serverError(w http.ResponseWriter, what string, err error) {
	log.Printf("auth: %s: %v", what, err)
	writeProblem(w, http.StatusInternalServerError, "Something went wrong")
}

// mailError: the code could not be sent; the user can ask for it again
func mailError(w http.ResponseWriter, err error) {
	log.Printf("auth: send code: %v", err)
	writeProblem(w, http.StatusBadGateway, "Could not send the email")
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
