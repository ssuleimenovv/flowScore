package live

import (
	"fmt"
	"math"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/insight"
)

// The factor groups (flow.Group) as the brief names them
var factorNames = map[string]string{
	"shots":         "удары",
	"goals":         "голы",
	"key_passes":    "ключевые передачи",
	"possession":    "владение",
	"corners":       "угловые",
	"cards":         "жёлтые карточки соперника",
	"red_cards":     "свои удаления",
	"substitutions": "замены",
	"other":         "другие события",
}

// briefOf tells the agent what the snapshot knows at this moment: the score,
// each team's Flow and what it is made of, the goals and red cards so far
// and the chances. The snapshot must have a prediction.
func briefOf(s Snapshot, key string) insight.Brief {
	return insight.Brief{
		Key:         key,
		Minute:      int(s.At.Minutes()),
		Competition: s.Match.Competition,
		Home:        teamOf(s, event.Home),
		Away:        teamOf(s, event.Away),
		Events:      keyEvents(s),
		Chances:     chancesOf(s.Prediction.Current),
		PreMatch:    chancesOf(s.Prediction.PreMatch),
	}
}

func teamOf(s Snapshot, side event.Side) insight.Team {
	t := insight.Team{
		Name:    s.Match.Home.Name,
		Goals:   s.Score.Home,
		Flow:    round(s.Flow.Home),
		Delta10: round(s.Delta10.Home),
	}
	if side == event.Away {
		t = insight.Team{
			Name:    s.Match.Away.Name,
			Goals:   s.Score.Away,
			Flow:    round(s.Flow.Away),
			Delta10: round(s.Delta10.Away),
		}
	}
	for _, f := range s.Factors {
		if f.Side == side {
			t.Factors = append(t.Factors, factorLine(f))
		}
	}
	return t
}

// factorLine is "удары +14 (4 за 7 мин)"; the count only when there are
// events in the last minutes, the share only for possession.
func factorLine(f FlowFactor) string {
	line := fmt.Sprintf("%s %+d", factorNames[f.Key], round(f.Value))
	if f.Share != nil {
		line += fmt.Sprintf(", доля мяча %d%%", round(*f.Share*100))
	}
	if f.Count > 0 {
		line += fmt.Sprintf(" (%d за %d мин)", f.Count, f.Minutes)
	}
	return line
}

// keyEvents are the goals and red cards so far: "23′ гол — Arsenal (Mesut Özil)".
func keyEvents(s Snapshot) []string {
	var out []string
	for _, e := range s.Events {
		var kind string
		switch e.Type {
		case event.Goal:
			kind = "гол"
		case event.RedCard:
			kind = "удаление"
		default:
			continue
		}

		minute := fmt.Sprintf("%d′", e.Minute)
		if e.AddedTime != nil {
			minute = fmt.Sprintf("%d+%d′", e.Minute, *e.AddedTime)
		}
		team := s.Match.Home.Name
		if e.Side != nil && *e.Side == event.Away {
			team = s.Match.Away.Name
		}
		line := fmt.Sprintf("%s %s — %s", minute, kind, team)
		if e.Player != nil {
			line += fmt.Sprintf(" (%s)", e.Player.Name)
		}
		out = append(out, line)
	}
	return out
}

func chancesOf(p Probabilities) insight.Chances {
	return insight.Chances{Home: p.Home, Draw: p.Draw, Away: p.Away}
}

func round(x float64) int {
	return int(math.Round(x))
}
