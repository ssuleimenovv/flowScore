package insight

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"path/filepath"
)

var brief = Brief{
	Key:         "1:m30",
	Minute:      30,
	Competition: "Premier League",
	Home:        Team{Name: "Arsenal", Goals: 1, Flow: 72, Delta10: 15, Factors: []string{"удары +14 (4 за 7 мин)"}},
	Away:        Team{Name: "Chelsea", Flow: 31, Delta10: -4},
	Events:      []string{"23′ гол — Arsenal (Mesut Özil)"},
	Chances:     Chances{Home: 61, Draw: 24, Away: 15},
	PreMatch:    Chances{Home: 45, Draw: 28, Away: 27},
}

func TestPromptHasTheFacts(t *testing.T) {
	p := brief.prompt()
	for _, want := range []string{
		"Минута: 30",
		"Счёт: Arsenal 1:0 Chelsea",
		"Flow Arsenal: 72 (за 10 минут +15)",
		"- удары +14 (4 за 7 мин)",
		"Flow Chelsea: 31 (за 10 минут -4)",
		"- 23′ гол — Arsenal (Mesut Özil)",
		"Шансы сейчас: П1 61%, X 24%, П2 15%",
	} {
		if !strings.Contains(p, want) {
			t.Errorf("prompt has no %q:\n%s", want, p)
		}
	}
}

// countingWriter answers every moment with its key and counts the requests
type countingWriter struct{ calls atomic.Int32 }

func (w *countingWriter) Write(_ context.Context, b Brief) (Text, error) {
	w.calls.Add(1)
	return Text{Title: b.Key}, nil
}

func TestAgentWritesEachMomentOnce(t *testing.T) {
	w := &countingWriter{}
	a := NewAgent(w, 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go a.Run(ctx)

	reply := make(chan Answer, 1)
	a.Ask(brief, reply)
	first := receive(t, reply)

	// The next replay asks again: the answer comes at once, without a request
	a.Ask(brief, reply)
	second := receive(t, reply)

	if first.Text.Title != "1:m30" || second != first {
		t.Errorf("answers %+v and %+v", first, second)
	}
	if n := w.calls.Load(); n != 1 {
		t.Errorf("%d requests, want 1", n)
	}
}

func TestAgentNeverBlocksTheMatch(t *testing.T) {
	a := NewAgent(&countingWriter{}, 0) // not running: the queue only fills
	reply := make(chan Answer, 1)
	done := make(chan struct{})
	go func() {
		for i := range queueSize + 5 {
			b := brief
			b.Key = string(rune('a' + i))
			a.Ask(b, reply)
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Ask blocked on a full queue")
	}
}

func receive(t *testing.T, reply <-chan Answer) Answer {
	t.Helper()
	select {
	case a := <-reply:
		return a
	case <-time.After(time.Second):
		t.Fatal("no answer")
		return Answer{}
	}
}

// fakeGemini answers every request with body and keeps what it was sent
func fakeGemini(t *testing.T, body string) (url string, sent *map[string]any, path *string) {
	t.Helper()
	sent, path = new(map[string]any), new(string)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*path = r.URL.Path
		raw, _ := io.ReadAll(r.Body)
		json.Unmarshal(raw, sent)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv.URL, sent, path
}

// The request Gemini gets, and the answer read back, against a fake API
func TestGeminiRequestAndAnswer(t *testing.T) {
	url, sent, path := fakeGemini(t, `{
		"candidates": [{
			"content": {"role": "model", "parts": [{"text": "{\"title\":\"Арсенал давит\",\"text\":\"Четыре удара за семь минут.\"}"}]},
			"finishReason": "STOP"
		}],
		"usageMetadata": {"promptTokenCount": 900, "thoughtsTokenCount": 40, "candidatesTokenCount": 120}
	}`)

	g, err := NewGemini(context.Background(), "test-key", "gemini-test", url)
	if err != nil {
		t.Fatal(err)
	}
	got, err := g.Write(context.Background(), brief)
	if err != nil {
		t.Fatal(err)
	}
	if got != (Text{Title: "Арсенал давит", Text: "Четыре удара за семь минут."}) {
		t.Errorf("got %+v", got)
	}

	if !strings.HasSuffix(*path, "/models/gemini-test:generateContent") {
		t.Errorf("path %s", *path)
	}
	config := (*sent)["generationConfig"].(map[string]any)
	if config["responseMimeType"] != "application/json" || config["responseJsonSchema"] == nil {
		t.Errorf("generationConfig %v", config)
	}
	if (*sent)["systemInstruction"] == nil {
		t.Error("no system instruction")
	}
}

func TestGeminiCutOffIsAnError(t *testing.T) {
	url, _, _ := fakeGemini(t, `{"candidates": [{
		"content": {"role": "model", "parts": [{"text": "{\"title\":\"Арс"}]},
		"finishReason": "MAX_TOKENS"
	}]}`)

	g, err := NewGemini(context.Background(), "test-key", "gemini-test", url)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Write(context.Background(), brief); err == nil {
		t.Error("a cut-off answer came back as text")
	}
}

func TestAgentKeepsTextsInAFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "insights.json")
	w := &countingWriter{}
	a := NewAgent(w, 0)
	if err := a.UseFile(path); err != nil { // no file yet: an empty one
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go a.Run(ctx)

	reply := make(chan Answer, 1)
	a.Ask(brief, reply)
	receive(t, reply)

	// The next start of the server, without a key: the text comes from the file
	again := NewAgent(nil, 0)
	if err := again.UseFile(path); err != nil {
		t.Fatal(err)
	}
	again.Ask(brief, reply)
	if got := receive(t, reply); got.Text.Title != "1:m30" || again.Len() != 1 {
		t.Errorf("answer %+v, %d texts", got, again.Len())
	}
}

// limitedWriter is out of quota for an hour
type limitedWriter struct{ calls atomic.Int32 }

func (w *limitedWriter) Write(context.Context, Brief) (Text, error) {
	w.calls.Add(1)
	return Text{}, &LimitError{Retry: time.Hour}
}

func TestAgentWaitsOutTheQuota(t *testing.T) {
	w := &limitedWriter{}
	a := NewAgent(w, 0)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go a.Run(ctx)

	reply := make(chan Answer, 1)
	a.Ask(brief, reply)
	deadline := time.Now().Add(time.Second)
	for {
		if _, _, limited := a.lookup(brief.Key); limited {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("the agent did not notice the limit")
		}
		time.Sleep(5 * time.Millisecond)
	}

	// Another moment while the quota is out: not even queued
	other := brief
	other.Key = "1:m40"
	a.Ask(other, reply)
	if n := len(a.jobs); n != 0 || w.calls.Load() != 1 {
		t.Errorf("%d queued, %d requests; want none after the limit", n, w.calls.Load())
	}
}

func TestGeminiOutOfQuota(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		io.WriteString(w, `{"error": {"code": 429, "status": "RESOURCE_EXHAUSTED", "message": "quota",
			"details": [{"@type": "type.googleapis.com/google.rpc.RetryInfo", "retryDelay": "44799s"}]}}`)
	}))
	defer srv.Close()

	g, err := NewGemini(context.Background(), "test-key", "gemini-test", srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	_, err = g.Write(context.Background(), brief)
	var limit *LimitError
	if !errors.As(err, &limit) || limit.Retry != 44799*time.Second {
		t.Errorf("error %v, want a LimitError for 44799s", err)
	}
}
