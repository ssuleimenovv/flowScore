package live

import (
	"math"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/flow"
	"github.com/ssuleimenovv/flowscore/services/internal/xg"
)

func TestShotsCarryOurXG(t *testing.T) {
	model, err := xg.Load("../../../ai/xg/model.json")
	if err != nil {
		t.Fatal(err)
	}
	hub := NewHub()
	client := hub.Subscribe("m1")
	p := NewPublisher(hub, NewStore(), "m1", flow.DefaultParams())
	p.UseXG(&model)
	p.Start(event.Match{ID: "m1"})

	// A header from the six-yard box after a corner, on the 0–100 pitch
	sbXG := 0.3
	at := 10 * time.Minute
	header := event.Event{
		ID: "s1", Type: event.ShotOnTarget, Side: event.Home, Period: 1, Elapsed: at,
		Pos: &event.Position{X: 95, Y: 50}, XG: &sbXG, Header: true, Situation: "set_piece",
	}
	p.Run(closed(flow.Update{At: at, Cause: &header}))

	got := read[MatchEvent](t, client, "match.event", nil)
	want := model.XG(xg.Shot{X: 114, Y: 40, Header: true, Situation: xg.SetPiece})
	if got.ModelXG == nil || math.Abs(*got.ModelXG-want) > 1e-12 {
		t.Errorf("modelXG = %v, want %.4f", got.ModelXG, want)
	}
	if *got.XG != sbXG {
		t.Errorf("xG = %v, want StatsBomb's %v kept", *got.XG, sbXG)
	}
}
