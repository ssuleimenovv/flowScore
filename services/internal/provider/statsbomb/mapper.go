package statsbomb

import (
	"cmp"
	"strconv"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

const (
	pitchLength = 120.0
	pitchWidth  = 80.0
)

// mapEvent converts one StatsBomb event into zero, one or two of our events.
// A foul with a card becomes two: the foul and the card. names comes from
// loadNames; a player missing from it keeps the name in the event.
func mapEvent(r rawEvent, matchID string, homeTeamID int, names map[int]string) []event.Event {
	base := event.Event{
		ID:      r.ID,
		MatchID: matchID,
		Side:    event.Away,
		Period:  r.Period,
		Elapsed: elapsed(r),
	}
	if r.Team.ID == homeTeamID {
		base.Side = event.Home
	}
	if r.Player != nil {
		base.PlayerID = strconv.Itoa(r.Player.ID)
		base.Player = cmp.Or(names[r.Player.ID], r.Player.Name)
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
		e.Header = r.Shot.BodyPart.Name == "Head"
		e.Situation = situation(r)
		return []event.Event{e}

	case "Pass":
		if r.Pass == nil {
			return nil
		}
		var out []event.Event
		if r.Pass.Type != nil && r.Pass.Type.Name == "Corner" {
			out = append(out, with(base, event.Corner))
		}
		// A corner that finds a head is a key pass too, so both can come out
		if r.Pass.ShotAssist || r.Pass.GoalAssist {
			key := with(base, event.KeyPass)
			key.ID = base.ID + ":key"
			out = append(out, key)
		}
		return out

	case "Duel":
		if r.Duel != nil && r.Duel.Type.Name == "Tackle" {
			return []event.Event{with(base, event.Tackle)}
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

	// StatsBomb gives an own goal to the team it counts for as "Own Goal For",
	// with no shot and no player
	case "Own Goal For":
		e := with(base, event.Goal)
		e.OwnGoal = true
		return []event.Event{e}

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

// situation is what led to the shot, the way ai/xg/shots.py names it for
// the xG model: a penalty, a shot straight from a free kick, play that
// started with a corner or a free kick, or open play.
func situation(r rawEvent) string {
	switch {
	case r.Shot.Type.Name == "Penalty":
		return "penalty"
	case r.Shot.Type.Name == "Free Kick":
		return "free_kick"
	case r.PlayPattern.Name == "From Corner" || r.PlayPattern.Name == "From Free Kick":
		return "set_piece"
	default:
		return "open_play"
	}
}
