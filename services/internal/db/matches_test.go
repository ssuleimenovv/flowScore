package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
)

func TestMatchesKeepsSnapshot(t *testing.T) {
	url := os.Getenv("FLOWSCORE_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("FLOWSCORE_TEST_DATABASE_URL is not set")
	}
	ctx := context.Background()
	pool, err := Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	t.Cleanup(func() { pool.Exec(ctx, `DELETE FROM matches WHERE id = 'test-1'`) })

	m := Matches{Pool: pool}
	snap := live.Snapshot{
		Match: event.Match{
			ID: "test-1", Source: "espn", CompetitionID: "eng.1",
			KickoffAt: time.Date(2026, 10, 17, 14, 0, 0, 0, time.UTC),
			Home:      event.Team{ID: "368", Name: "Everton"},
			Away:      event.Team{ID: "363", Name: "Chelsea"},
		},
		Status: "finished",
		Score:  live.Score{Away: 2},
		At:     95 * time.Minute,
		Events: []live.MatchEvent{{ID: "e1", Minute: 30}},
	}
	// Twice: a match saved again replaces itself
	for range 2 {
		if err := m.Save(ctx, snap); err != nil {
			t.Fatal(err)
		}
	}

	got, ok := m.Get("test-1")
	if !ok || got.Status != "finished" || got.Score != snap.Score || got.At != snap.At ||
		got.Match.Away.Name != "Chelsea" || !got.Match.KickoffAt.Equal(snap.Match.KickoffAt) ||
		len(got.Events) != 1 || got.Events[0].Minute != 30 {
		t.Errorf("got %+v, %v", got, ok)
	}
	if _, ok := m.Get("none"); ok {
		t.Error("a match never saved was found")
	}
}
