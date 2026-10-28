package api

import (
	"context"
	"encoding/json"
	"github.com/ssuleimenovv/flowscore/services/internal/auth"
	"github.com/ssuleimenovv/flowscore/services/internal/db"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// fakeUsers is the users table as a map, by email
type fakeUsers map[string]struct {
	user db.User
	hash string
}

func (f fakeUsers) Create(_ context.Context, email, name, hash string) (db.User, error) {
	if _, ok := f[email]; ok {
		return db.User{}, db.ErrEmailTaken
	}
	u := db.User{ID: "u" + email, Email: email, Name: name}
	f[email] = struct {
		user db.User
		hash string
	}{u, hash}
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

func TestAuth(t *testing.T) {
	tokens := auth.Tokens{Secret: []byte("a secret of at least thirty-two bytes"), TTL: 15 * time.Minute}
	mux := http.NewServeMux()
	RegisterAuth(mux, fakeUsers{}, tokens)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	var s session
	res := post(t, srv.URL+"/api/v1/auth/register", `{"name":"Аян","email":" Ayan@Mail.com ","password":"flowscore1"}`, &s)
	if res.StatusCode != http.StatusCreated || s.User.Email != "ayan@mail.com" || s.User.Name != "Аян" {
		t.Fatalf("sign up: %d %+v", res.StatusCode, s)
	}
	if id, err := tokens.UserID(s.Token); err != nil || id != s.User.ID {
		t.Errorf("token names %q, %v; want %q", id, err, s.User.ID)
	}

	for _, c := range []struct {
		name, path, body string
		want             int
	}{
		{"same email", "register", `{"name":"Другой","email":"ayan@mail.com","password":"flowscore1"}`, http.StatusConflict},
		{"weak password", "register", `{"name":"Аян","email":"new@mail.com","password":"short"}`, http.StatusUnprocessableEntity},
		{"no name", "register", `{"email":"new@mail.com","password":"flowscore1"}`, http.StatusUnprocessableEntity},
		{"sign in", "login", `{"email":"AYAN@mail.com","password":"flowscore1"}`, http.StatusOK},
		{"wrong password", "login", `{"email":"ayan@mail.com","password":"flowscore2"}`, http.StatusUnauthorized},
		{"unknown email", "login", `{"email":"nobody@mail.com","password":"flowscore1"}`, http.StatusUnauthorized},
		{"not JSON", "login", `email=ayan`, http.StatusBadRequest},
	} {
		var body map[string]any
		if res := post(t, srv.URL+"/api/v1/auth/"+c.path, c.body, &body); res.StatusCode != c.want {
			t.Errorf("%s: %d, want %d: %v", c.name, res.StatusCode, c.want, body)
		}
	}
}

func TestMe(t *testing.T) {
	tokens := auth.Tokens{Secret: []byte("a secret of at least thirty-two bytes"), TTL: 15 * time.Minute}
	users := fakeUsers{}
	u, _ := users.Create(context.Background(), "ayan@mail.com", "Аян", "hash")
	mux := http.NewServeMux()
	RegisterAuth(mux, users, tokens)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	fresh, _ := tokens.Issue(u.ID, time.Now())
	expired, _ := tokens.Issue(u.ID, time.Now().Add(-time.Hour))
	gone, _ := tokens.Issue("deleted", time.Now())

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
		req, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/me", nil)
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
