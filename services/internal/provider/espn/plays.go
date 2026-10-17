package espn

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// playTypes are ESPN's play types that are events of ours. Kick-off, delays
// and the whistles are not here: the summary's clock already tells them.
// The own goal is named the way ESPN names the rest; no recording has had
// one yet.
var playTypes = map[string]event.Type{
	"goal":              event.Goal,
	"own-goal":          event.Goal,
	"penalty---scored":  event.Goal,
	"penalty---saved":   event.ShotOnTarget,
	"penalty---missed":  event.ShotOffTarget,
	"shot-on-target":    event.ShotOnTarget,
	"shot-off-target":   event.ShotOffTarget,
	"shot-blocked":      event.ShotBlocked,
	"shot-hit-woodwork": event.ShotOffTarget, // the bar or the post: not on target, as Opta counts it
	"corner-awarded":    event.Corner,
	"yellow-card":       event.YellowCard,
	"red-card":          event.RedCard,
	"substitution":      event.Substitution,
	"foul":              event.Foul,
	"handball":          event.Foul,
	"offside":           event.Offside,
}

// typeOf is our event for an ESPN play type. A goal comes in kinds,
// "goal---header", "goal---volley", and more than the recordings have shown.
func typeOf(play string) (event.Type, bool) {
	if strings.HasPrefix(play, "goal---") {
		return event.Goal, true
	}
	t, ok := playTypes[play]
	return t, ok
}

// ParsePlays reads the commentary of a summary: every play that is an event
// of ours, oldest first, with a key pass before each assisted shot. A play
// names its team only, so m tells the sides.
func ParsePlays(body []byte, m event.Match) ([]event.Event, error) {
	var raw struct {
		Commentary []struct {
			Play *struct {
				ID   string `json:"id"`
				Type struct {
					Type string `json:"type"` // "shot-on-target"
				} `json:"type"`
				Text   string `json:"text"`
				Period struct {
					Number int `json:"number"`
				} `json:"period"`
				Clock struct {
					Value        float64 `json:"value"`        // seconds; stands at 45:00 and 90:00 in added time
					DisplayValue string  `json:"displayValue"` // "90'+2'"
				} `json:"clock"`
				Team struct {
					DisplayName string `json:"displayName"`
				} `json:"team"`
				Participants []struct {
					Athlete struct {
						ID          string `json:"id"`
						DisplayName string `json:"displayName"`
					} `json:"athlete"`
				} `json:"participants"`
				X *float64 `json:"fieldPositionX"` // 0–100, attacking left to right, as ours
				Y *float64 `json:"fieldPositionY"`
			} `json:"play"`
		} `json:"commentary"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	var out []event.Event
	seen := map[string]bool{}
	for _, c := range raw.Commentary {
		p := c.Play
		if p == nil {
			continue // lineups, warm-up: words without a play
		}
		// A foul comes twice, as the free kick won and as the foul: one play
		if seen[p.ID] {
			continue
		}
		seen[p.ID] = true

		t, ok := typeOf(p.Type.Type)
		if !ok {
			continue
		}
		e := event.Event{
			ID:      m.ID + ":" + p.ID,
			MatchID: m.ID,
			Type:    t,
			Side:    event.Away,
			Period:  p.Period.Number,
			Elapsed: time.Duration(p.Clock.Value)*time.Second + added(p.Clock.DisplayValue),
			OwnGoal: p.Type.Type == "own-goal",
		}
		if p.Team.DisplayName == m.Home.Name {
			e.Side = event.Home
		}
		// The one who did it comes first, then the assist or the one fouled
		if len(p.Participants) > 0 {
			a := p.Participants[0].Athlete
			e.Player, e.PlayerID = a.DisplayName, a.ID
		}
		if p.X != nil && p.Y != nil {
			e.Pos = &event.Position{X: *p.X, Y: *p.Y}
		}
		if isShot(t) {
			e.Header = strings.Contains(p.Text, " header ")
			e.Situation = situation(p.Type.Type, p.Text)
			// "Assisted by Jacob Murphy with a cross": the second one set it up
			if len(p.Participants) > 1 && strings.Contains(p.Text, "Assisted by") {
				a := p.Participants[1].Athlete
				out = append(out, event.Event{
					ID:       e.ID + ":key",
					MatchID:  m.ID,
					Type:     event.KeyPass,
					Side:     e.Side,
					Period:   e.Period,
					Elapsed:  e.Elapsed,
					Player:   a.DisplayName,
					PlayerID: a.ID,
				})
			}
		}
		out = append(out, e)
	}
	return out, nil
}

// added is the stoppage time of a clock like "90'+2'": the minutes past
// 45:00 or 90:00 are only in the text. The second added minute starts at
// 91:00, which our clock shows as 90+2.
func added(clock string) time.Duration {
	_, after, ok := strings.Cut(clock, "+")
	if !ok {
		return 0
	}
	n, _ := strconv.Atoi(strings.TrimSuffix(after, "'"))
	return time.Duration(max(n-1, 0)) * time.Minute
}

func isShot(t event.Type) bool {
	return t == event.Goal || t == event.ShotOnTarget || t == event.ShotOffTarget || t == event.ShotBlocked
}

// situation is what led to a shot, as xg.Situation names it. The commentary
// only says so in words, Opta's: "from a direct free kick", "following a
// corner", "following a set piece situation".
func situation(play, text string) string {
	switch {
	case strings.HasPrefix(play, "penalty"):
		return "penalty"
	case strings.Contains(text, "direct free kick"):
		return "free_kick"
	case strings.Contains(text, "following a corner"), strings.Contains(text, "following a set piece"):
		return "set_piece"
	default:
		return "open_play"
	}
}
