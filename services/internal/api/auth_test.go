package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/auth"
	"github.com/ssuleimenovv/flowscore/services/internal/db"
	"github.com/ssuleimenovv/flowscore/services/internal/mail"
)

// fakeUsers is the users table as a map, by email
type fakeUsers map[string]*fakeUser

type fakeUser struct {
	user db.User
	hash string
}

func (f fakeUsers) Create(_ context.Context, email, name, hash string) (db.User, error) {
	if _, ok := f[email]; ok {
		return db.User{}, db.ErrEmailTaken
	}
	u := db.User{ID: "u-" + email, Email: email, Name: name}
	f[email] = &fakeUser{u, hash}
	return u, nil
}

func (f fakeUsers) ByEmail(_ context.Context, email string) (db.User, string, error) {
	row, ok := f[email]
	if !ok {
		return db.User{}, "", db.ErrNoUser
	}
	return row.user, row.hash, nil
}

func (f fakeUsers) ByID(_ context.Context, id string) (db.User, error) {
	for _, row := range f {
		if row.user.ID == id {
			return row.user, nil
		}
	}
	return db.User{}, db.ErrNoUser
}

func (f fakeUsers) MarkVerified(_ context.Context, id string) error {
	for _, row := range f {
		if row.user.ID == id {
			row.user.Verified = true
		}
	}
	return nil
}

// fakeCodes is the codes table: a code once put stays "too soon" to replace
type fakeCodes map[string]*fakeCode

type fakeCode struct {
	hash  string
	tries int
}

func (f fakeCodes) Put(_ context.Context, userID, purpose, hash string, _, gap time.Duration) error {
	key := userID + ":" + purpose
	if _, ok := f[key]; ok && gap > 0 {
		return db.ErrTooSoon
	}
	f[key] = &fakeCode{hash: hash}
	return nil
}

func (f fakeCodes) Use(_ context.Context, userID, purpose, hash string) (int, error) {
	key := userID + ":" + purpose
	c, ok := f[key]
	if !ok || c.tries >= db.MaxAttempts {
		return 0, db.ErrNoCode
	}
	if c.hash == hash {
		delete(f, key)
		return 0, nil
	}
	c.tries++
	return db.MaxAttempts - c.tries, db.ErrWrongCode
}

// outbox keeps the emails instead of sending them
type outbox struct {
	sent []mail.Message
}

func (o *outbox) Send(_ context.Context, m mail.Message) error {
	o.sent = append(o.sent, m)
	return nil
}

// code is the code of the last email, from its subject "Код FlowScore: 482193"
func (o *outbox) code() string {
	return strings.TrimPrefix(o.sent[len(o.sent)-1].Subject, "Код FlowScore: ")
}

type authServer struct {
	url    string
	users  fakeUsers
	codes  fakeCodes
	box    *outbox
	tokens auth.Tokens
}

func newAuthServer(t *testing.T) authServer {
	t.Helper()
	s := authServer{
		users:  fakeUsers{},
		codes:  fakeCodes{},
		box:    &outbox{},
		tokens: auth.Tokens{Secret: []byte("a secret of at least thirty-two bytes"), TTL: 15 * time.Minute},
	}
	mux := http.NewServeMux()
	Auth{
		Users:  s.users,
		Codes:  s.codes,
		Hasher: auth.Codes{Secret: []byte("a secret of at least thirty-two bytes")},
		Mail:   s.box,
		Tokens: s.tokens,
	}.Register(mux)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	s.url = srv.URL + "/api/v1"
	return s
}

func TestSignUpWithCode(t *testing.T) {
	s := newAuthServer(t)

	var p pending
	res := post(t, s.url+"/auth/register", `{"name":"Аян","email":" Ayan@Mail.com ","password":"flowscore1"}`, &p)
	if res.StatusCode != http.StatusAccepted || p.Email != "ayan@mail.com" || len(s.box.sent) != 1 || s.box.sent[0].To != "ayan@mail.com" {
		t.Fatalf("sign up: %d %+v, %d emails", res.StatusCode, p, len(s.box.sent))
	}
	code := s.box.code()

	// Before the code: the right password is not enough
	var body map[string]any
	if res := post(t, s.url+"/auth/login", `{"email":"ayan@mail.com","password":"flowscore1"}`, &body); res.StatusCode != http.StatusForbidden {
		t.Errorf("sign in before the code: %d, want 403", res.StatusCode)
	}

	// A wrong code tells the tries left
	wrong := "000000"
	if code == wrong {
		wrong = "111111"
	}
	var w struct{ AttemptsLeft int }
	res = post(t, s.url+"/auth/verify", `{"email":"ayan@mail.com","code":"`+wrong+`"}`, &w)
	if res.StatusCode != http.StatusUnprocessableEntity || w.AttemptsLeft != db.MaxAttempts-1 {
		t.Errorf("wrong code: %d, %d left", res.StatusCode, w.AttemptsLeft)
	}

	// The right code signs in
	var sess session
	res = post(t, s.url+"/auth/verify", `{"email":"ayan@mail.com","code":"`+code+`"}`, &sess)
	if res.StatusCode != http.StatusOK || sess.User.Name != "Аян" {
		t.Fatalf("right code: %d %+v", res.StatusCode, sess)
	}
	if id, err := s.tokens.UserID(sess.Token); err != nil || id != sess.User.ID {
		t.Errorf("token names %q, %v", id, err)
	}
	// It worked once; from now on the password is enough
	if res := post(t, s.url+"/auth/verify", `{"email":"ayan@mail.com","code":"`+code+`"}`, &body); res.StatusCode != http.StatusGone {
		t.Errorf("the code again: %d, want 410", res.StatusCode)
	}
	if res := post(t, s.url+"/auth/login", `{"email":"AYAN@mail.com","password":"flowscore1"}`, &sess); res.StatusCode != http.StatusOK {
		t.Errorf("sign in after the code: %d", res.StatusCode)
	}

	for _, c := range []struct {
		name, path, body string
		want             int
	}{
		{"same email", "register", `{"name":"Другой","email":"ayan@mail.com","password":"flowscore1"}`, http.StatusConflict},
		{"weak password", "register", `{"name":"Аян","email":"new@mail.com","password":"short"}`, http.StatusUnprocessableEntity},
		{"no name", "register", `{"email":"new@mail.com","password":"flowscore1"}`, http.StatusUnprocessableEntity},
		{"wrong password", "login", `{"email":"ayan@mail.com","password":"flowscore2"}`, http.StatusUnauthorized},
		{"unknown email", "login", `{"email":"nobody@mail.com","password":"flowscore1"}`, http.StatusUnauthorized},
		{"unknown email's code", "verify", `{"email":"nobody@mail.com","code":"123456"}`, http.StatusGone},
		{"not JSON", "login", `email=ayan`, http.StatusBadRequest},
	} {
		var body map[string]any
		if res := post(t, s.url+"/auth/"+c.path, c.body, &body); res.StatusCode != c.want {
			t.Errorf("%s: %d, want %d: %v", c.name, res.StatusCode, c.want, body)
		}
	}
}

func TestResend(t *testing.T) {
	s := newAuthServer(t)
	var body map[string]any

	// An unknown email gets the same answer, and no email goes out
	if res := post(t, s.url+"/auth/resend", `{"email":"nobody@mail.com"}`, &body); res.StatusCode != http.StatusAccepted || len(s.box.sent) != 0 {
		t.Errorf("unknown email: %d, %d emails", res.StatusCode, len(s.box.sent))
	}

	post(t, s.url+"/auth/register", `{"name":"Аян","email":"ayan@mail.com","password":"flowscore1"}`, &body)
	// Right after the first code: wait a minute
	res := post(t, s.url+"/auth/resend", `{"email":"ayan@mail.com"}`, &body)
	if res.StatusCode != http.StatusTooManyRequests || res.Header.Get("Retry-After") != "60" {
		t.Errorf("resend at once: %d, Retry-After %q", res.StatusCode, res.Header.Get("Retry-After"))
	}

	// A minute later (the code is gone from the fake): a new one goes out
	delete(s.codes, "u-ayan@mail.com:"+db.Verify)
	if res := post(t, s.url+"/auth/resend", `{"email":"ayan@mail.com"}`, &body); res.StatusCode != http.StatusAccepted || len(s.box.sent) != 2 {
		t.Errorf("resend later: %d, %d emails", res.StatusCode, len(s.box.sent))
	}
}

func TestMe(t *testing.T) {
	s := newAuthServer(t)
	u, _ := s.users.Create(context.Background(), "ayan@mail.com", "Аян", "hash")

	fresh, _ := s.tokens.Issue(u.ID, time.Now())
	expired, _ := s.tokens.Issue(u.ID, time.Now().Add(-time.Hour))
	gone, _ := s.tokens.Issue("deleted", time.Now())

	for _, c := range []struct {
		name, header string
		want         int
	}{
		{"signed in", "Bearer " + fresh, http.StatusOK},
		{"no token", "", http.StatusUnauthorized},
		{"not Bearer", "Basic " + fresh, http.StatusUnauthorized},
		{"expired", "Bearer " + expired, http.StatusUnauthorized},
		{"account gone", "Bearer " + gone, http.StatusUnauthorized},
	} {
		req, _ := http.NewRequest(http.MethodGet, s.url+"/me", nil)
		if c.header != "" {
			req.Header.Set("Authorization", c.header)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		var body user
		json.NewDecoder(res.Body).Decode(&body)
		res.Body.Close()
		if res.StatusCode != c.want {
			t.Errorf("%s: %d, want %d", c.name, res.StatusCode, c.want)
		}
		if c.want == http.StatusOK && body.Name != "Аян" {
			t.Errorf("%s: %+v", c.name, body)
		}
	}
}
