// Command export writes a season of StatsBomb matches as one CSV of normalized
// events: the dataset Flow is calibrated on (ai/calibration). It goes through
// the same mapper as the live replay, so the model learns from exactly the
// events the Flow Engine sees.
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/statsbomb"
)

var header = []string{"match_id", "date", "period", "second", "side", "type", "xg", "home_share"}

// match is one loaded match, or the reason it could not be loaded.
type match struct {
	info   event.Match
	events []event.Event
	err    error
}

func main() {
	dir := flag.String("dir", "data/statsbomb", "files from ai/calibration/download.py")
	season := flag.String("season", "2-27", "competition-season of the matches file")
	out := flag.String("out", "../ai/data/events.csv", "CSV to write")
	flag.Parse()

	matchesPath := filepath.Join(*dir, "matches-"+*season+".json")
	ids, err := statsbomb.MatchIDs(matchesPath)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	matches := loadAll(*dir, matchesPath, ids)

	rows, err := write(*out, matches)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d matches, %d events → %s in %v", len(matches), rows, *out, time.Since(start).Round(time.Millisecond))
}

// loadAll parses the matches on every CPU. Each worker writes only its own
// slot of the result, so the slice needs no lock and keeps the season order.
func loadAll(dir, matchesPath string, ids []int) []match {
	results := make([]match, len(ids))
	jobs := make(chan int)

	var wg sync.WaitGroup
	for range runtime.NumCPU() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				results[i] = load(dir, matchesPath, ids[i])
			}
		}()
	}
	for i := range ids {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return results
}

func load(dir, matchesPath string, id int) match {
	info, err := statsbomb.LoadMatchInfo(matchesPath, id)
	if err != nil {
		return match{err: err}
	}
	name := strconv.Itoa(id) + ".json"
	events, err := statsbomb.LoadMatch(matchesPath,
		filepath.Join(dir, "events-"+name), filepath.Join(dir, "lineups-"+name), id)
	if err != nil {
		return match{err: fmt.Errorf("match %d: %w", id, err)}
	}
	return match{info: info, events: events}
}

func write(path string, matches []match) (rows int, err error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return 0, err
	}
	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.Write(header); err != nil {
		return 0, err
	}
	for _, m := range matches {
		if m.err != nil {
			return rows, m.err
		}
		date := m.info.KickoffAt.Format(time.DateOnly)
		for _, e := range m.events {
			if err := w.Write(row(m.info.ID, date, e)); err != nil {
				return rows, err
			}
			rows++
		}
	}
	w.Flush()
	return rows, w.Error()
}

func row(matchID, date string, e event.Event) []string {
	return []string{
		matchID,
		date,
		strconv.Itoa(e.Period),
		strconv.Itoa(int(e.Elapsed.Seconds())),
		string(e.Side),
		string(e.Type),
		optional(e.XG),
		optional(e.HomeShare),
	}
}

// optional writes a missing number as an empty cell, which pandas reads as NaN.
func optional(v *float64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatFloat(*v, 'f', 4, 64)
}
