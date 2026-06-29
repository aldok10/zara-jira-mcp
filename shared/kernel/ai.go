package kernel

import "context"

// Provider defines the interface for AI completions.
// Used by sprint and handler modules for AI-powered analysis.
type Provider interface {
	Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}
