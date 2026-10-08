package statsbomb

type ref struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type rawShot struct {
	XG       float64 `json:"statsbomb_xg"`
	Outcome  ref     `json:"outcome"`
	BodyPart ref     `json:"body_part"` // "Head", "Right Foot"
	Type     ref     `json:"type"`      // "Open Play", "Free Kick", "Penalty"
}

type rawEvent struct {
	ID             string    `json:"id"`
	Period         int       `json:"period"`
	Minute         int       `json:"minute"`
	Second         int       `json:"second"`
	Type           ref       `json:"type"`
	Team           ref       `json:"team"`
	Player         *ref      `json:"player"`
	Location       []float64 `json:"location"`
	PossessionTeam ref       `json:"possession_team"`
	PlayPattern    ref       `json:"play_pattern"` // what the possession started with: "From Corner"

	Shot *rawShot `json:"shot"`

	Pass *struct {
		Type       *ref `json:"type"`
		ShotAssist bool `json:"shot_assist"`
		GoalAssist bool `json:"goal_assist"`
	} `json:"pass"`

	Duel *struct {
		Type ref `json:"type"`
	} `json:"duel"`

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
	MatchID     int    `json:"match_id"`
	MatchDate   string `json:"match_date"`
	KickOff     string `json:"kick_off"`
	MatchWeek   int    `json:"match_week"`
	Competition struct {
		ID   int    `json:"competition_id"`
		Name string `json:"competition_name"`
	} `json:"competition"`
	Stadium *struct {
		Name string `json:"name"`
	} `json:"stadium"`
	HomeTeam struct {
		ID   int    `json:"home_team_id"`
		Name string `json:"home_team_name"`
	} `json:"home_team"`
	AwayTeam struct {
		ID   int    `json:"away_team_id"`
		Name string `json:"away_team_name"`
	} `json:"away_team"`
}
