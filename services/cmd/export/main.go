// Command export writes a season of StatsBomb matches for Flow calibration
// (ai/calibration):
//
//   - events.csv: every normalized event. It goes through the same mapper as the
//     live replay, so the model learns from exactly the events the Flow Engine sees.
//   - flow.csv: Flow of both teams at every whole minute with the current
//     params, the reference the Python port of Flow is checked against.
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
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/provider/statsbomb"
)

// match is one loaded match, or the reason it could not be loaded.
type match struct {
	info   event.Match
	events []event.Event
	err    error
}

func main() {
	dir := flag.String("dir", "data/statsbomb", "files from ai/calibration/download.py")
	season := flag.String("season", "2-27", "competition-season of the matches file")
	out := flag.String("out", "../ai/data", "folder for events.csv and flow.csv")
	flag.Parse()

	matchesPath := filepath.Join(*dir, "matches-"+*season+".json")
	ids, err := statsbomb.MatchIDs(matchesPath)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	matches := loadAll(*dir, matchesPath, ids)

	events := [][]string{{"match_id", "date", "period", "second", "side", "type", "xg", "home_share"}}
	samples := [][]string{{"match_id", "period", "second", "home", "away"}}
	params := flow.DefaultParams()

	for _, m := range matches {
		if m.err != nil {
			log.Fatal(m.err)
		}
		date := m.info.KickoffAt.Format(time.DateOnly)
		for _, e := range m.events {
			events = append(events, eventRow(m.info.ID, date, e))
		}
		for _, s := range flow.Trace(m.events, params) {
			samples = append(samples, []string{
				m.info.ID, strconv.Itoa(s.Period), seconds(s.At), exact(s.Home), exact(s.Away),
			})
		}
	}

	for name, rows := range map[string][][]string{"events.csv": events, "flow.csv": samples} {
		if err := writeCSV(filepath.Join(*out, name), rows); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("%d matches: %d events, %d Flow samples → %s in %v",
		len(matches), len(events)-1, len(samples)-1, *out, time.Since(start).Round(time.Millisecond))
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

func writeCSV(path string, rows [][]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// WriteAll flushes and reports the first error, so nothing is lost silently
	return csv.NewWriter(f).WriteAll(rows)
}

func eventRow(matchID, date string, e event.Event) []string {
	return []string{
		matchID,
		date,
		strconv.Itoa(e.Period),
		seconds(e.Elapsed),
		string(e.Side),
		string(e.Type),
		optional(e.XG),
		optional(e.HomeShare),
	}
}

func seconds(d time.Duration) string {
	return strconv.Itoa(int(d.Seconds()))
}

// optional writes a missing number as an empty cell, which pandas reads as NaN.
func optional(v *float64) string {
	if v == nil {
		return ""
	}
	return exact(*v)
}

// exact writes the shortest text that reads back as the same float64, so the
// Python side works with the very numbers Go had.
func exact(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}
