package statsbomb

import (
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

const (
	pitchLength = 120.0
	pitchWidth  = 80.0
)

// mapEvent converts one StatsBomb event into zero, one or two of our events.
// A foul with a card becomes two: the foul and the card.
func mapEvent(r rawEvent, matchID string, homeTeamID int) []event.Event {
	base := event.Event{
		ID:      r.ID,
		MatchID: matchID,
		Side:    event.Away,
		Period:  r.Period,
		Elapsed: time.Duration(r.Minute)*time.Minute + time.Duration(r.Second)*time.Second,
	}
	if r.Team.ID == homeTeamID {
		base.Side = event.Home
	}
	if r.Player != nil {
		base.Player = r.Player.Name
	}
	if len(r.Location) == 2 {
		base.Pos = &event.Position{
			X: r.Location[0] / pitchLength * 100,
			Y: r.Location[1] / pitchWidth * 100,
		}
	}

	switch r.Type.Name {
	case "Shot":
		if r.Shot == nil {
			return nil
		}
		e := base
		xg := r.Shot.XG
		e.XG = &xg
		e.Type = shotType(r.Shot.Outcome.Name)
		return []event.Event{e}

	case "Pass":
		if r.Pass != nil && r.Pass.Type != nil && r.Pass.Type.Name == "Corner" {
			return []event.Event{with(base, event.Corner)}
		}

	case "Dribble":
		if r.Dribble != nil && r.Dribble.Outcome != nil && r.Dribble.Outcome.Name == "Complete" {
			return []event.Event{with(base, event.TakeOn)}
		}

	case "Foul Committed":
		out := []event.Event{with(base, event.Foul)}
		if r.FoulCommitted != nil && r.FoulCommitted.Card != nil {
			out = append(out, cardEvent(base, r.FoulCommitted.Card.Name))
		}
		return out

	case "Bad Behaviour":
		if r.BadBehaviour != nil && r.BadBehaviour.Card != nil {
			return []event.Event{cardEvent(base, r.BadBehaviour.Card.Name)}
		}

	case "Substitution":
		return []event.Event{with(base, event.Substitution)}

	case "Offside":
		return []event.Event{with(base, event.Offside)}
	}
	return nil
}

func shotType(outcome string) event.Type {
	switch outcome {
	case "Goal":
		return event.Goal
	case "Saved", "Saved to Post":
		return event.ShotOnTarget
	case "Blocked":
		return event.ShotBlocked
	default: // Off T, Post, Wayward, Saved Off T
		return event.ShotOffTarget
	}
}

func cardEvent(base event.Event, card string) event.Event {
	e := with(base, event.YellowCard)
	if card == "Red Card" || card == "Second Yellow" {
		e.Type = event.RedCard
	}
	e.ID = base.ID + ":card"
	return e
}

func with(base event.Event, t event.Type) event.Event {
	base.Type = t
	return base
}
