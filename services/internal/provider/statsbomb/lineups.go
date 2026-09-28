package statsbomb

import "fmt"

type rawLineup struct {
	Lineup []struct {
		ID       int     `json:"player_id"`
		Name     string  `json:"player_name"`
		Nickname *string `json:"player_nickname"`
	} `json:"lineup"`
}

// loadNames returns the name to show for every player of the match. StatsBomb
// names players in full ("Sergio Leonel Agüero del Castillo"); the nickname is
// the one fans know ("Sergio Agüero"), when there is one.
func loadNames(path string) (map[int]string, error) {
	var teams []rawLineup
	if err := readJSON(path, &teams); err != nil {
		return nil, fmt.Errorf("read lineups: %w", err)
	}

	names := map[int]string{}
	for _, team := range teams {
		for _, p := range team.Lineup {
			names[p.ID] = p.Name
			if p.Nickname != nil && *p.Nickname != "" {
				names[p.ID] = *p.Nickname
			}
		}
	}
	return names, nil
}
