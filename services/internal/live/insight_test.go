package live

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/insight"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
)

func TestBriefOf(t *testing.T) {
	home := event.Home
	added := 2
	share := 0.58
	snap := Snapshot{
		Match:   event.Match{Competition: "Premier League", Home: event.Team{Name: "Arsenal"}, Away: event.Team{Name: "Chelsea"}},
		Score:   Score{Home: 1},
		At:      47*time.Minute + 10*time.Second,
		Flow:    FlowValues{Home: 71.6, Away: 30.8},
		Delta10: FlowValues{Home: 14.6, Away: -4.2},
		Factors: []FlowFactor{
			{Side: event.Home, Key: "shots", Value: 13.7, Count: 4, Minutes: 7},
			{Side: event.Away, Key: "possession", Value: 3.2, Share: &share},
		},
		Events: []MatchEvent{
			{Type: event.ShotOnTarget, Side: &home, Minute: 20},
			{Type: event.Goal, Side: &home, Minute: 45, AddedTime: &added, Player: &PersonRef{Name: "Mesut Özil"}},
		},
		Prediction: &Prediction{Current: Probabilities{61, 24, 15}, PreMatch: Probabilities{45, 28, 27}},
	}

	b := briefOf(snap, "m1:m40")
	if b.Key != "m1:m40" || b.Minute != 47 || b.Competition != "Premier League" {
		t.Errorf("brief = %+v", b)
	}
	wantHome := insight.Team{Name: "Arsenal", Goals: 1, Flow: 72, Delta10: 15, Factors: []string{"удары +14 (4 за 7 мин)"}}
	if !slices.Equal(b.Home.Factors, wantHome.Factors) || b.Home.Name != wantHome.Name ||
		b.Home.Goals != 1 || b.Home.Flow != 72 || b.Home.Delta10 != 15 {
		t.Errorf("home = %+v, want %+v", b.Home, wantHome)
	}
	if !slices.Equal(b.Away.Factors, []string{"владение +3, доля мяча 58%"}) || b.Away.Delta10 != -4 {
		t.Errorf("away = %+v", b.Away)
	}
	// Only the goals and the red cards, with added time
	if !slices.Equal(b.Events, []string{"45+2′ гол — Arsenal (Mesut Özil)"}) {
		t.Errorf("events = %q", b.Events)
	}
	if b.Chances != (insight.Chances{Home: 61, Draw: 24, Away: 15}) || b.PreMatch.Home != 45 {
		t.Errorf("chances = %+v, %+v", b.Chances, b.PreMatch)
	}
}

// echoWriter writes the key of the moment as the title
type echoWriter struct{}

func (echoWriter) Write(_ context.Context, b insight.Brief) (insight.Text, error) {
	return insight.Text{Title: b.Key, Text: "разбор"}, nil
}

func TestPublisherExplainsKeyMoments(t *testing.T) {
	model, err := predict.Load("../../../ai/prediction/model.json")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	agent := insight.NewAgent(echoWriter{}, 0)
	go agent.Run(ctx)

	hub := NewHub()
	client := hub.Subscribe("m1")
	store := NewStore()
	p := NewPublisher(hub, store, "m1", flow.DefaultParams())
	p.UseModel(&model)
	p.UseAgent(agent)
	p.Start(event.Match{ID: "m1", Home: event.Team{Name: "Arsenal"}, Away: event.Team{Name: "Chelsea"}})

	updates := make(chan flow.Update)
	done := make(chan struct{})
	go func() {
		p.Run(updates)
		close(done)
	}()

	// A goal in the 23rd minute, then the ticks up to the 30th: two moments
	goalAt := 23 * time.Minute
	updates <- flow.Update{At: goalAt, Cause: &event.Event{ID: "g1", Type: event.Goal, Side: event.Home, Period: 1, Elapsed: goalAt}}
	updates <- flow.Update{At: 29 * time.Minute}
	updates <- flow.Update{At: 30 * time.Minute}

	// The answers come back while the match goes on
	deadline := time.Now().Add(2 * time.Second)
	for {
		snap, _ := store.Get("m1")
		if snap.Explanation != nil && snap.Explanation.Title == "m1:m30" {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("explanation = %+v, want the one of the 30th minute", snap.Explanation)
		}
		time.Sleep(10 * time.Millisecond)
	}
	close(updates)
	<-done

	first := read[Explanation](t, client, "insight.update", nil)
	if first.Title != "m1:goal:g1" || first.Minute != 23 {
		t.Errorf("first insight.update = %+v, want the goal", first)
	}
}

func TestAnswerFromAheadIsDropped(t *testing.T) {
	store := NewStore()
	p := NewPublisher(NewHub(), store, "m1", flow.DefaultParams())
	p.Start(event.Match{ID: "m1"})
	p.at = 5 * time.Minute

	// The last round's answer about the 30th minute, at the 5th of this one
	p.explain(insight.Answer{Minute: 30, Text: insight.Text{Title: "рано"}})
	if snap, _ := store.Get("m1"); snap.Explanation != nil {
		t.Errorf("explanation = %+v, want none yet", snap.Explanation)
	}

	p.explain(insight.Answer{Minute: 3, Text: insight.Text{Title: "вовремя"}})
	if snap, _ := store.Get("m1"); snap.Explanation == nil || snap.Explanation.Title != "вовремя" {
		t.Errorf("explanation = %+v", snap.Explanation)
	}
}
