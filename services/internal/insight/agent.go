package insight

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

// Text is the analysis of one moment: a headline and two or three sentences.
type Text struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// Answer is the text of a moment, sent back to the match that asked.
type Answer struct {
	Key    string
	Minute int
	Text   Text
}

// Writer turns the facts of a moment into text: an LLM in production, a fake
// in tests.
type Writer interface {
	Write(ctx context.Context, b Brief) (Text, error)
}

// LimitError is a Writer's error when the API has run out of quota: the
// agent asks nothing until Retry has passed.
type LimitError struct {
	Retry time.Duration
}

func (e *LimitError) Error() string {
	return fmt.Sprintf("out of quota, retry in %s", e.Retry)
}

// A longer backlog is stale before it is written: the match has moved on
const queueSize = 8

// How long one answer may take. A demo match lasts two and a half minutes:
// an answer later than this is about a moment long gone.
const writeTimeout = 20 * time.Second

// Agent writes the analysis for every match, one request at a time. A demo
// match replays the same moments round after round, so each text is written
// once and kept, in memory and in a file: the server asks the API only about
// moments it has never seen.
type Agent struct {
	writer Writer        // nil: the agent only answers from what it has
	gap    time.Duration // the pause after each request, to stay under the API's limit
	jobs   chan job

	mu    sync.Mutex
	done  map[string]Text // by Brief.Key
	file  string          // where done is kept, "" for nowhere
	until time.Time       // no requests before it: the quota is out
}

type job struct {
	brief Brief
	reply chan<- Answer
}

func NewAgent(w Writer, gap time.Duration) *Agent {
	return &Agent{writer: w, gap: gap, jobs: make(chan job, queueSize), done: map[string]Text{}}
}

// UseFile loads the texts kept in path and keeps every new one there. A
// missing file is an empty one: it appears with the first text.
func (a *Agent) UseFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(data) > 0 {
		if err := json.Unmarshal(data, &a.done); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
	}
	a.file = path
	return nil
}

// Len is how many moments the agent has a text for.
func (a *Agent) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.done)
}

// Ask asks for the text of a moment; it arrives on reply later, or at once
// when it is already written. Ask never blocks the match that calls it.
func (a *Agent) Ask(b Brief, reply chan<- Answer) {
	t, ok, limited := a.lookup(b.Key)
	switch {
	case ok:
		send(reply, Answer{Key: b.Key, Minute: b.Minute, Text: t})
		return
	case a.writer == nil || limited:
		return // nobody to ask, or not now
	}
	select {
	case a.jobs <- job{brief: b, reply: reply}:
	default:
		log.Printf("insight %s: the queue is full, skipped", b.Key)
	}
}

// Run writes the asked moments until ctx is done.
func (a *Agent) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case j := <-a.jobs:
			if wait := a.answer(ctx, j); wait > 0 && !pause(ctx, wait) {
				return
			}
		}
	}
}

// answer sends the text of the moment and returns how long to wait before
// the next request: nothing after a text it already had, the gap after a
// request, until the quota is back after a LimitError.
func (a *Agent) answer(ctx context.Context, j job) time.Duration {
	b := j.brief
	t, ok, _ := a.lookup(b.Key) // the moment may have been asked twice
	if !ok {
		ctx, cancel := context.WithTimeout(ctx, writeTimeout)
		defer cancel()
		var err error
		if t, err = a.writer.Write(ctx, b); err != nil {
			log.Printf("insight %s: %v", b.Key, err)
			var limit *LimitError
			if errors.As(err, &limit) {
				a.mu.Lock()
				a.until = time.Now().Add(limit.Retry)
				a.mu.Unlock()
				return limit.Retry
			}
			return a.gap
		}
		a.keep(b.Key, t)
	}
	send(j.reply, Answer{Key: b.Key, Minute: b.Minute, Text: t})
	if ok {
		return 0
	}
	return a.gap
}

// lookup returns the text of the moment if there is one, and whether the
// quota is out right now.
func (a *Agent) lookup(key string) (t Text, ok, limited bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	t, ok = a.done[key]
	return t, ok, time.Now().Before(a.until)
}

// keep remembers the text and writes the whole file again. A failed write is
// logged once and the agent goes on in memory: on a host with a read-only
// disk the texts it has are still served.
func (a *Agent) keep(key string, t Text) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.done[key] = t
	if a.file == "" {
		return
	}
	// Map keys come out sorted, so the file changes only where a text is new
	data, err := json.MarshalIndent(a.done, "", "  ")
	if err == nil {
		// A new file renamed over the old one: never a half-written file
		tmp := a.file + ".tmp"
		if err = os.WriteFile(tmp, data, 0o644); err == nil {
			err = os.Rename(tmp, a.file)
		}
	}
	if err != nil {
		log.Printf("insight: keeping the texts in %s: %v; from now on in memory only", a.file, err)
		a.file = ""
	}
}

// send drops the answer rather than wait: a match that is over reads no more.
func send(reply chan<- Answer, a Answer) {
	select {
	case reply <- a:
	default:
	}
}

// pause waits for d and reports false if ctx is done first.
func pause(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
	}
}
