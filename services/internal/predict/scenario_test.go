package predict

import (
	"testing"
	"time"
)

// A red at the current minute is the same as a red already shown
func TestRedNowIsPredict(t *testing.T) {
	m := load(t)
	s := Situation{Period: 2, At: 60 * time.Minute, HomeGoals: 1}
	got := m.Simulate(0.2, s, []Red{{Home: true, Minute: 60}})

	shown := s
	shown.HomeReds = 1
	if want := m.Predict(0.2, shown); got != want {
		t.Errorf("red at 60' simulated %+v, shown %+v", got, want)
	}
}

func TestLaterRedHurtsLess(t *testing.T) {
	m := load(t)
	s := Situation{Period: 2, At: 50 * time.Minute}
	none := m.Predict(0, s).Home
	early := m.Simulate(0, s, []Red{{Home: true, Minute: 55}}).Home
	late := m.Simulate(0, s, []Red{{Home: true, Minute: 85}}).Home
	if !(early < late && late < none) {
		t.Errorf("home win: red at 55' %.3f, at 85' %.3f, none %.3f; want them rising", early, late, none)
	}
}

func TestOtherTeamsRedHelps(t *testing.T) {
	m := load(t)
	s := Situation{Period: 1, At: 20 * time.Minute}
	if m.Simulate(0, s, []Red{{Home: false, Minute: 30}}).Home <= m.Predict(0, s).Home {
		t.Error("an away red should raise the home win chance")
	}
}

// Two reds in any order give the same match
func TestRedOrderDoesNotMatter(t *testing.T) {
	m := load(t)
	s := Situation{Period: 2, At: 46 * time.Minute}
	a := m.Simulate(0, s, []Red{{Home: true, Minute: 80}, {Home: false, Minute: 60}})
	b := m.Simulate(0, s, []Red{{Home: false, Minute: 60}, {Home: true, Minute: 80}})
	if a != b {
		t.Errorf("order changed the result: %+v vs %+v", a, b)
	}
}
