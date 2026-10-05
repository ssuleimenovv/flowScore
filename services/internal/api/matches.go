package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
	"github.com/ssuleimenovv/flowscore/services/internal/live"
	"github.com/ssuleimenovv/flowscore/services/internal/predict"
)

// Register adds the match routes of api/openapi.yaml.
func Register(mux *http.ServeMux, store *live.Store) {
	mux.HandleFunc("GET /api/v1/matches", listHandler(store))
	mux.HandleFunc("GET /api/v1/matches/{matchId}", withSnapshot(store, matchResponse))
	mux.HandleFunc("GET /api/v1/matches/{matchId}/flow", withSnapshot(store, flowResponse))
	mux.HandleFunc("GET /api/v1/matches/{matchId}/events", withSnapshot(store, eventsResponse))
	mux.HandleFunc("GET /api/v1/matches/{matchId}/insight", insightHandler(store))
	mux.HandleFunc("POST /api/v1/matches/{matchId}/simulate", simulateHandler(store))

}

// withSnapshot loads the match once for every route and answers 404 if it is unknown.
func withSnapshot(store *live.Store, render func(live.Snapshot) any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap, ok := store.Get(r.PathValue("matchId"))
		if !ok {
			writeProblem(w, http.StatusNotFound, "Match not found")
			return
		}
		writeJSON(w, http.StatusOK, render(snap))
	}
}

type team struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type competition struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Round string `json:"round"`
}

type venue struct {
	Name string `json:"name"`
}

type match struct {
	ID            string         `json:"id"`
	Status        string         `json:"status"`
	KickoffAt     time.Time      `json:"kickoffAt"`
	Competition   competition    `json:"competition"`
	Venue         venue          `json:"venue"`
	Home          team           `json:"home"`
	Away          team           `json:"away"`
	Score         live.Score     `json:"score"`
	HalftimeScore *live.Score    `json:"halftimeScore"`
	Clock         live.Clock     `json:"clock"`
	Stats         []live.StatRow `json:"stats"`
	Seq           int64          `json:"seq"`
}

func matchResponse(s live.Snapshot) any {
	return toMatch(s)
}

func toMatch(s live.Snapshot) match {
	m := s.Match
	return match{
		ID:            m.ID,
		Status:        s.Status,
		KickoffAt:     m.KickoffAt,
		Competition:   competition{ID: m.CompetitionID, Name: m.Competition, Round: m.Round},
		Venue:         venue{Name: m.Venue},
		Home:          toTeam(m.Home),
		Away:          toTeam(m.Away),
		Score:         s.Score,
		HalftimeScore: s.Halftime,
		Clock:         live.ClockOf(s.Period, s.At, s.UpdatedAt),
		Stats:         orEmpty(s.Stats),
		Seq:           s.Seq,
	}
}

// matchSummary is one card of the home screen: the match header plus what the
// card draws, the flow wave, its explanation and the chances.
type matchSummary struct {
	match
	Flow       live.FlowValues   `json:"flow"`
	Delta10    live.FlowValues   `json:"delta10"`
	Factors    []live.FlowFactor `json:"factors"`
	Points     []live.FlowPoint  `json:"points"`
	Prediction *live.Prediction  `json:"prediction"`
}

type matchList struct {
	Items []matchSummary `json:"items"`
}

// listHandler answers the home screen: every match, by kick-off time.
func listHandler(store *live.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		snaps := store.List()
		items := make([]matchSummary, 0, len(snaps))
		for _, s := range snaps {
			items = append(items, matchSummary{
				match:      toMatch(s),
				Flow:       s.Flow,
				Delta10:    s.Delta10,
				Factors:    orEmpty(s.Factors),
				Points:     orEmpty(s.Points),
				Prediction: s.Prediction,
			})
		}
		writeJSON(w, http.StatusOK, matchList{Items: items})
	}
}

type flowSeries struct {
	MatchID   string            `json:"matchId"`
	Current   live.FlowValues   `json:"current"`
	Delta10   live.FlowValues   `json:"delta10"`
	Factors   []live.FlowFactor `json:"factors"`
	Points    []live.FlowPoint  `json:"points"`
	UpdatedAt time.Time         `json:"updatedAt"`
	Seq       int64             `json:"seq"`
}

func flowResponse(s live.Snapshot) any {
	return flowSeries{
		MatchID:   s.Match.ID,
		Current:   s.Flow,
		Delta10:   s.Delta10,
		Factors:   orEmpty(s.Factors),
		Points:    orEmpty(s.Points),
		UpdatedAt: s.UpdatedAt,
		Seq:       s.Seq,
	}
}

type eventList struct {
	Items []live.MatchEvent `json:"items"`
	Seq   int64             `json:"seq"`
}

func eventsResponse(s live.Snapshot) any {
	items := orEmpty(s.Events)
	slices.Reverse(items) // the contract returns newest first
	return eventList{Items: items, Seq: s.Seq}
}

type insight struct {
	MatchID     string           `json:"matchId"`
	Explanation *struct{}        `json:"explanation"` // the Explainability agent comes later
	Prediction  *live.Prediction `json:"prediction"`
	Sentiment   *struct{}        `json:"sentiment"`
	Seq         int64            `json:"seq"`
}

// insightHandler answers 404 for a match without a prediction as well: the
// contract has no insight without one, and the gateway may run without a model.
func insightHandler(store *live.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap, ok := store.Get(r.PathValue("matchId"))
		if !ok || snap.Prediction == nil {
			writeProblem(w, http.StatusNotFound, "No insight for this match")
			return
		}
		writeJSON(w, http.StatusOK, insight{MatchID: snap.Match.ID, Prediction: snap.Prediction, Seq: snap.Seq})
	}
}

type simulationRequest struct {
	Reds []struct {
		Side   event.Side `json:"side"`
		Minute int        `json:"minute"`
	} `json:"reds"`
}

type simulation struct {
	Minute   int                `json:"minute"`
	Current  live.Probabilities `json:"current"`
	Scenario live.Probabilities `json:"scenario"`
}

// A match has two teams of 11, and a team with fewer than 7 is abandoned
const maxReds = 8

// simulateHandler answers "what if": the chances with reds added later in the
// match, next to the chances as it stands, both from the same moment.
func simulateHandler(store *live.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap, ok := store.Get(r.PathValue("matchId"))
		if !ok || snap.Outlook == nil {
			writeProblem(w, http.StatusNotFound, "No outcome model for this match")
			return
		}

		var req simulationRequest
		body := http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(body).Decode(&req); err != nil || len(req.Reds) > maxReds {
			writeProblem(w, http.StatusBadRequest, "Invalid simulation")
			return
		}
		reds := make([]predict.Red, 0, len(req.Reds))
		for _, red := range req.Reds {
			if (red.Side != event.Home && red.Side != event.Away) || red.Minute < 1 || red.Minute > 120 {
				writeProblem(w, http.StatusBadRequest, "Invalid red card")
				return
			}
			reds = append(reds, predict.Red{Home: red.Side == event.Home, Minute: red.Minute})
		}

		now, scenario := snap.Outlook.Simulate(reds)
		writeJSON(w, http.StatusOK, simulation{Minute: snap.Outlook.Minute(), Current: now, Scenario: scenario})
	}
}

func toTeam(t event.Team) team {
	return team{ID: t.ID, Code: t.Code, Name: t.Name}
}

// orEmpty turns a nil slice into an empty one: JSON gets [] instead of null.
func orEmpty[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}

// problem is RFC 9457 Problem Details (Problem in the contract).
type problem struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
}

func writeProblem(w http.ResponseWriter, status int, title string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(problem{Type: "about:blank", Title: title, Status: status}); err != nil {
		log.Printf("write problem: %v", err)
	}
}
