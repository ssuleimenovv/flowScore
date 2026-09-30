package flow

import (
	"math"
	"testing"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

func TestTraceSamplesEveryMinute(t *testing.T) {
	p := DefaultParams()
	samples := Trace([]event.Event{
		ev(event.Goal, event.Home, 1, 44),  // on the minute: already in the 44:00 sample
		ev(event.Foul, event.Away, 2, 46)}, // the second half runs to 46:00
		p)

	// 44 minutes of the first half, then 46:00 of the second
	if len(samples) != 45 {
		t.Fatalf("got %d samples, want 45", len(samples))
	}

	near(t, "before the goal", samples[42].Home, floor)
	near(t, "44:00", samples[43].Home, 75.10)

	// The break keeps 70% of the impulse, then one minute of decay from 45:00
	second := samples[44]
	if second.Period != 2 || second.At != 46*time.Minute {
		t.Errorf("last sample = %+v, want period 2 at 46:00", second)
	}
	impulse := 23 * p.HalftimeKeep * math.Exp(-1/p.Tau)
	near(t, "46:00", second.Home, 100*(1-math.Exp(-(p.Base+impulse)/p.K)))
	near(t, "away", second.Away, floor)
}
