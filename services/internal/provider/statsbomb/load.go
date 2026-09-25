package statsbomb

import (
	"cmp"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func LoadMatch(matchesPath, eventsPath string, matchID int) ([]event.Event, error) {
	var matches []rawMatch
	if err := readJSON(matchesPath, &matches); err != nil {
		return nil, fmt.Errorf("read matches: %w", err)
	}

	homeTeamID := 0
	for _, m := range matches {
		if m.MatchID == matchID {
			homeTeamID = m.HomeTeam.ID
			break
		}
	}
	if homeTeamID == 0 {
		return nil, fmt.Errorf("match %d not found in %s", matchID, matchesPath)
	}

	var raw []rawEvent
	if err := readJSON(eventsPath, &raw); err != nil {
		return nil, fmt.Errorf("read events: %w", err)
	}

	id := strconv.Itoa(matchID)
	var out []event.Event
	for _, r := range raw {
		out = append(out, mapEvent(r, id, homeTeamID)...)
	}
	out = append(out, possessionEvents(raw, id, homeTeamID)...)

	// Possession events are appended at the end; put everything back in match order.
	slices.SortStableFunc(out, func(a, b event.Event) int {
		if a.Period != b.Period {
			return cmp.Compare(a.Period, b.Period)
		}
		return cmp.Compare(a.Elapsed, b.Elapsed)
	})
	return out, nil
}

func readJSON(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}
