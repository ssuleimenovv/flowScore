package live

import (
	"strings"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
)

func TestToPercentsAddUpTo100(t *testing.T) {
	tests := []struct {
		in   predict.Outcome
		want Probabilities
	}{
		{predict.Outcome{Home: 1.0 / 3, Draw: 1.0 / 3, Away: 1.0 / 3}, Probabilities{34, 33, 33}},
		{predict.Outcome{Home: 0.3814, Draw: 0.2950, Away: 0.3236}, Probabilities{38, 30, 32}},
		{predict.Outcome{Home: 0.005, Draw: 0.006, Away: 0.989}, Probabilities{0, 1, 99}},
		{predict.Outcome{Home: 1}, Probabilities{100, 0, 0}},
	}
	for _, tt := range tests {
		if got := toPercents(tt.in); got != tt.want {
			t.Errorf("toPercents(%+v) = %+v, want %+v", tt.in, got, tt.want)
		}
	}
}

func TestPublisherPredicts(t *testing.T) {
	model, err := predict.Load("../../../ai/prediction/model.json")
	if err != nil {
		t.Fatal(err)
	}
	hub := NewHub()
	client := hub.Subscribe("m1")
	store := NewStore()
	p := NewPublisher(hub, store, "m1", flow.DefaultParams())
	p.UseModel(&model)

	match := event.Match{ID: "m1", Home: event.Team{Name: "Manchester City"}, Away: event.Team{Name: "Arsenal"}}
	p.Start(match)
	before, _ := store.Get("m1")
	if before.Prediction == nil || before.Prediction.Current != before.Prediction.PreMatch {
		t.Fatalf("prediction at kick-off = %+v, want the pre-match chances", before.Prediction)
	}
	pre := before.Prediction.PreMatch

	goalAt := 30 * time.Minute
	p.Run(closed(
		flow.Update{At: goalAt, Cause: &event.Event{Type: event.Goal, Side: event.Home, Period: 1, Elapsed: goalAt}},
		flow.Update{At: goalAt + 5*time.Second}, // a tick: the chances barely move
	))

	afterGoal := read[Probabilities](t, client, "prediction.update", nil)
	if afterGoal.Home <= pre.Home {
		t.Errorf("home chance after 1:0 = %d, not above %d before", afterGoal.Home, pre.Home)
	}
	final := read[Probabilities](t, client, "prediction.update", func(p Probabilities) bool { return p.Home == 100 })
	if final != (Probabilities{Home: 100}) {
		t.Errorf("after the final whistle = %+v, want 100 · 0 · 0", final)
	}

	snap, _ := store.Get("m1")
	if snap.Prediction.Current != final || snap.Prediction.PreMatch != pre || snap.Prediction.Model != model.Name {
		t.Errorf("snapshot prediction = %+v", snap.Prediction)
	}
}

func TestWithoutModelNoPrediction(t *testing.T) {
	hub := NewHub()
	client := hub.Subscribe("m1")
	store := NewStore()
	p := NewPublisher(hub, store, "m1", flow.DefaultParams())

	p.Start(event.Match{ID: "m1"})
	p.Run(closed(flow.Update{At: time.Minute}))

	snap, _ := store.Get("m1")
	if snap.Prediction != nil {
		t.Errorf("prediction = %+v, want nil", snap.Prediction)
	}
	// Run has finished, so every message it sent is in the buffer
	for len(client.Messages()) > 0 {
		if msg := <-client.Messages(); strings.Contains(string(msg), `"prediction.update"`) {
			t.Errorf("sent %s without a model", msg)
		}
	}
}
