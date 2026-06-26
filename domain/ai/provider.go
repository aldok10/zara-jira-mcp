package ai

import "context"

// Provider defines the interface for AI completions.
type Provider interface {
	Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}
