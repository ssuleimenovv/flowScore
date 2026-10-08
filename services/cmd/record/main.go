// Command record saves what ESPN says about real matches while they are
// played, so the live provider can be built and tested on real data later.
//
// Usage (from services/):
//
//	go run ./cmd/record
//	go run ./cmd/record -leagues eng.1,esp.1 -every 20s
//
// Every poll it reads each league's scoreboard and, for every match being
// played, the summary. A summary is kept only when it changed since the last
// one, gzipped, as data/espn/<league>/<match>/<time>.json.gz. A finished
// match is kept once more; with -pre, a match not started yet once as well.
// Stop it with Ctrl+C.
package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/provider/espn"
)

func main() {
	leagues := flag.String("leagues", "eng.1", "ESPN leagues, comma-separated: eng.1, esp.1, ger.1, ita.1, fra.1")
	out := flag.String("out", "data/espn", "folder for the recordings")
	every := flag.Duration("every", 30*time.Second, "how often to poll")
	pre := flag.Bool("pre", false, "also keep each match not started yet, once")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	r := recorder{
		client: espn.NewClient(),
		out:    *out,
		pre:    *pre,
		last:   map[string][32]byte{},
		kept:   map[string]bool{},
	}
	log.Printf("recording %s every %s into %s", *leagues, *every, *out)
	for {
		for _, league := range strings.Split(*leagues, ",") {
			r.poll(ctx, strings.TrimSpace(league))
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(*every):
		}
	}
}

type recorder struct {
	client *espn.Client
	out    string
	pre    bool
	last   map[string][32]byte // the hash of each match's last kept summary
	kept   map[string]bool     // matches kept once and done with: finished, or not started
}

func (r *recorder) poll(ctx context.Context, league string) {
	games, err := r.client.Scoreboard(ctx, league)
	if err != nil {
		log.Printf("%s: %v", league, err)
		return
	}
	live := 0
	for _, g := range games {
		switch {
		case g.State == "in":
			live++
			r.keep(ctx, league, g)
		case g.State == "post" && !r.kept[g.ID+"post"]:
			r.kept[g.ID+"post"] = r.keep(ctx, league, g)
		case g.State == "pre" && r.pre && !r.kept[g.ID+"pre"]:
			r.kept[g.ID+"pre"] = r.keep(ctx, league, g)
		}
	}
	log.Printf("%s: %d matches, %d live", league, len(games), live)
}

// keep saves the match's summary if it changed and reports whether the
// summary is on disk now.
func (r *recorder) keep(ctx context.Context, league string, g espn.Game) bool {
	body, err := r.client.Summary(ctx, league, g.ID)
	if err != nil {
		log.Printf("%s %s: %v", league, g.Name, err)
		return false
	}
	sum := sha256.Sum256(body)
	if sum == r.last[g.ID] {
		return true
	}

	dir := filepath.Join(r.out, league, g.ID)
	name := filepath.Join(dir, time.Now().UTC().Format("20060102T150405Z")+".json.gz")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Print(err)
		return false
	}
	var zipped bytes.Buffer
	w := gzip.NewWriter(&zipped)
	w.Write(body)
	w.Close()
	if err := os.WriteFile(name, zipped.Bytes(), 0o644); err != nil {
		log.Print(err)
		return false
	}
	r.last[g.ID] = sum
	log.Printf("%s %s (%s): kept %d KB", league, g.Name, g.State, zipped.Len()/1024)
	return true
}
