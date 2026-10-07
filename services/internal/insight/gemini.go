package insight

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/genai"
)

// Gemini is the Writer that asks the Gemini API.
type Gemini struct {
	client *genai.Client
	model  string
}

// NewGemini makes a writer with the API key; baseURL is empty except in tests.
func NewGemini(ctx context.Context, apiKey, model, baseURL string) (*Gemini, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:      apiKey,
		Backend:     genai.BackendGeminiAPI,
		HTTPOptions: genai.HTTPOptions{BaseURL: baseURL},
	})
	if err != nil {
		return nil, err
	}
	return &Gemini{client: client, model: model}, nil
}

func (g *Gemini) Write(ctx context.Context, b Brief) (Text, error) {
	resp, err := g.client.Models.GenerateContent(ctx, g.model, genai.Text(b.prompt()), &genai.GenerateContentConfig{
		SystemInstruction:  genai.NewContentFromText(instructions, genai.RoleUser),
		ResponseMIMEType:   "application/json",
		ResponseJsonSchema: textSchema,
		// Retelling the facts needs next to no thought, and thinking is most of the wait
		ThinkingConfig: &genai.ThinkingConfig{ThinkingLevel: genai.ThinkingLevelMinimal},
	})
	if err != nil {
		return Text{}, limited(err)
	}

	if u := resp.UsageMetadata; u != nil {
		log.Printf("insight %s: %d tokens in, %d thinking, %d out",
			b.Key, u.PromptTokenCount, u.ThoughtsTokenCount, u.CandidatesTokenCount)
	}

	if len(resp.Candidates) == 0 {
		return Text{}, errors.New("no answer")
	}
	if reason := resp.Candidates[0].FinishReason; reason != genai.FinishReasonStop {
		return Text{}, fmt.Errorf("stopped: %s", reason)
	}
	var out Text
	if err := json.Unmarshal([]byte(resp.Text()), &out); err != nil {
		return Text{}, fmt.Errorf("answer: %w", err)
	}
	if out.Title == "" || out.Text == "" {
		return Text{}, errors.New("empty answer")
	}
	return out, nil
}

// limited turns the API's 429 into a LimitError with the wait the API asks
// for: "retryDelay": "44799s" when the quota for the day is out
func limited(err error) error {
	var apiErr genai.APIError
	if !errors.As(err, &apiErr) || apiErr.Code != http.StatusTooManyRequests {
		return err
	}
	retry := time.Minute // a limit per minute, when the API does not say
	for _, d := range apiErr.Details {
		if s, ok := d["retryDelay"].(string); ok {
			if r, err := time.ParseDuration(s); err == nil {
				retry = r
			}
		}
	}
	return &LimitError{Retry: retry}
}
