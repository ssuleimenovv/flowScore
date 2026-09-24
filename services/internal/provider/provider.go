package provider

import (
	"context"

	"github.com/ssuleimenovv/flowscore/services/internal/event"
)

// Provider streams normalized events of the match
// The channel is closed when the match ends or ctx is cancelled
type Provider interface {
	Stream(ctx context.Context, matchID string) (<-chan event.Event, error)
}
