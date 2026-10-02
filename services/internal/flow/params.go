package flow

import "github.com/ssuleimenovv/flowscore/services/internal/event"

// Params are the v1 values from docs/FLOW.md, section 4: the mockup's v0 with
// the three weights that season validation found stable (section 8).
type Params struct {
	Tau          float64 // impulse memory, match minutes
	K            float64 // scale of the 0–100 curve
	Base         float64 // S₀: impulse of a team that is just playing, the Flow floor
	HalftimeKeep float64 // share of impulse kept after the break (r)
	ShotBase     float64
	ShotPerXG    float64
	DefaultXG    float64 // used when a shot has no xG
	Possession   float64 // w_p: weight per minute for each point of share above 50%
	Weights      map[event.Type]float64
}

func DefaultParams() Params {
	return Params{
		Tau:          5,
		K:            18.7,
		Base:         3,
		HalftimeKeep: 0.7,
		ShotBase:     6,
		ShotPerXG:    7.3,
		DefaultXG:    0.1,
		Possession:   9,
		Weights: map[event.Type]float64{
			event.Goal:         23, // kept by design, see docs/FLOW.md, section 8
			event.YellowCard:   5,
			event.Corner:       4,
			event.Substitution: 5,
			event.KeyPass:      10,
			event.RedCard:      -18, // to the team that lost a player
		},
	}
}

// Weight returns which team the event benefits and by how much
func (p Params) Weight(e event.Event) (event.Side, float64) {
	switch e.Type {
	case event.ShotOnTarget, event.ShotOffTarget, event.ShotBlocked:
		xg := p.DefaultXG
		if e.XG != nil {
			xg = *e.XG
		}
		return e.Side, p.ShotBase + p.ShotPerXG*xg
	case event.YellowCard:
		return e.Side.Opponent(), p.Weights[e.Type]
	default:
		return e.Side, p.Weights[e.Type]
	}
}
