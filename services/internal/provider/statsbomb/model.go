package statsbomb

type ref struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type rawEvent struct {
	ID       string    `json:"id"`
	Period   int       `json:"period"`
	Minute   int       `json:"minute"`
	Second   int       `json:"second"`
	Type     ref       `json:"type"`
	Team     ref       `json:"team"`
	Player   *ref      `json:"player"`
	Location []float64 `json:"location"`

	Shot *struct {
		XG      float64 `json:"statsbomb_xg"`
		Outcome ref     `json:"outcome"`
	} `json:"shot"`

	Pass *struct {
		Type *ref `json:"type"`
	} `json:"pass"`

	Dribble *struct {
		Outcome *ref `json:"outcome"`
	} `json:"dribble"`

	FoulCommitted *struct {
		Card *ref `json:"card"`
	} `json:"foul_committed"`

	BadBehaviour *struct {
		Card *ref `json:"card"`
	} `json:"bad_behaviour"`
}

type rawMatch struct {
	MatchID  int `json:"match_id"`
	HomeTeam struct {
		ID int `json:"home_team_id"`
	} `json:"home_team"`
}
