package predict

// RatingFor is the rating at which the model, at kick-off, gives the home
// win and the away win the chances o does: what the bookmakers think of the
// two teams, in the model's own terms. The model's draw stays its own.
func (m Model) RatingFor(o Outcome) float64 {
	want := o.Home - o.Away
	// The home chance minus the away chance grows with the rating
	lo, hi := -4.0, 4.0
	for range 40 {
		mid := (lo + hi) / 2
		if p := m.Predict(mid, Situation{Period: 1}); p.Home-p.Away < want {
			lo = mid
		} else {
			hi = mid
		}
	}
	return (lo + hi) / 2
}
