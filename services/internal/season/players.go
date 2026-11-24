package season

import (
	"math"
	"strings"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/xg"
)

// PlayerMatch is what one player did in a match the season has read
type PlayerMatch struct {
	ID       string
	Name     string
	Side     event.Side
	Starter  bool
	Minutes  int
	Goals    int // own goals are not counted
	Assists  int
	Shots    int // a goal is a shot on target too
	OnTarget int
	XG       float64 // our model's, over the player's shots
}

// Players is what each player who got on the pitch did: the starters and
// the substitutes who came on. ESPN's commentary names the players without
// their IDs, so a player is found by the ID when there is one, else by the
// side and the name. nil without the lineups, which ESPN gives for most
// matches but not all.

func Players(events []event.Event, l *event.Lineups, model *xg.Model) []PlayerMatch {
	if l == nil {
		return nil
	}
	end := matchEnd(events)
	var players []PlayerMatch
	index := map[string]int{}  // a player's place in players, by ID
	byName := map[string]int{} // and by side and name: ESPN's commentary names a player without the ID
	for _, side := range []struct {
		side    event.Side
		players []event.Player
	}{{event.Home, l.Home.Players}, {event.Away, l.Away.Players}} {
		for _, p := range side.players {
			if !p.Starter && p.On == nil {
				continue // stayed on the bench
			}
			index[p.ID] = len(players)
			byName[string(side.side)+" "+p.Name] = len(players)
			players = append(players, PlayerMatch{
				ID: p.ID, Name: p.Name, Side: side.side, Starter: p.Starter,
				Minutes: minutes(p, end),
			})
		}
	}

	// "Assisted by": the key pass is the goal's ID with ":key" after it
	goals := map[string]bool{}
	for _, e := range events {
		if e.Type == event.Goal && !e.OwnGoal {
			goals[e.ID] = true
		}
	}
	for _, e := range events {
		i, ok := index[e.PlayerID]
		if !ok {
			i, ok = byName[string(e.Side)+" "+e.Player]
		}
		if !ok {
			continue
		}

		p := &players[i]
		switch {
		case e.Type == event.KeyPass && goals[strings.TrimSuffix(e.ID, ":key")]:
			p.Assists++
		case isShot(e):
			p.Shots++
			if e.Type == event.Goal || e.Type == event.ShotOnTarget {
				p.OnTarget++
			}
			if e.Type == event.Goal {
				p.Goals++
			}
			if v, ok := shotXG(e, model); ok {
				p.XG += v
			}
		}
	}
	return players
}

// matchEnd is the last moment of the match: 90 minutes, or more with
// stoppage time and extra time
func matchEnd(events []event.Event) time.Duration {
	end := 90 * time.Minute
	for _, e := range events {
		end = max(end, e.Elapsed)
	}
	return end
}

// minutes is how long a player was on the pitch, to the whole minute
func minutes(p event.Player, end time.Duration) int {
	from, to := time.Duration(0), end
	if p.On != nil {
		from = p.On.Elapsed
	}
	if p.Off != nil {
		to = p.Off.Elapsed
	}
	return int(math.Round(max(0, to-from).Minutes()))
}
