package mail

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBrevoSends(t *testing.T) {
	var got struct {
		Sender  struct{ Email string }
		To      []struct{ Email string }
		Subject string
	}
	var key string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key = r.Header.Get("api-key")
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusCreated)
	}))
	t.Cleanup(srv.Close)

	b := Brevo{Key: "k1", From: "flowscore@gmail.com", URL: srv.URL}
	if err := b.Send(context.Background(), CodeMail("ayan@mail.kz", "482193", false)); err != nil {
		t.Fatal(err)
	}
	if key != "k1" || got.Sender.Email != "flowscore@gmail.com" || len(got.To) != 1 ||
		got.To[0].Email != "ayan@mail.kz" || got.Subject != "Код FlowScore: 482193" {
		t.Errorf("sent %+v with key %q", got, key)
	}
}

func TestBrevoReportsRefusal(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message":"sender not validated"}`, http.StatusBadRequest)
	}))
	t.Cleanup(srv.Close)

	err := Brevo{Key: "k1", From: "x@gmail.com", URL: srv.URL}.Send(context.Background(), CodeMail("a@b.kz", "1", false))
	if err == nil || !strings.Contains(err.Error(), "sender not validated") {
		t.Errorf("err = %v", err)
	}
}

func TestCodeMail(t *testing.T) {
	m := CodeMail("a@b.kz", "482193", true)
	if !strings.Contains(m.Text, "482193") || !strings.Contains(m.HTML, "482193") ||
		!strings.Contains(m.Text, "новый пароль") {
		t.Errorf("mail %+v", m)
	}
}
