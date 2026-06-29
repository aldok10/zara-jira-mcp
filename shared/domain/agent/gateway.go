package agent

import "context"

// Gateway is the channel-agnostic interface for sending messages back to users.
// Inspired by OpenClaw's multi-channel gateway and Hermes' 20 platform adapters.
type Gateway interface {
	// SendText sends a plain text reply to a channel.
	SendText(ctx context.Context, channelID, text string) error

	// SendMarkdown sends a formatted reply.
	SendMarkdown(ctx context.Context, channelID, title, content string) error

	// Channel returns the name of this gateway (lark, slack, etc.)
	Channel() string
}

// Config holds the personality and behavior configuration.
// This is the "SOUL" — inspired by OpenClaw's SOUL.md and Hermes' identity system.
type Config struct {
	Name        string `json:"name"`
	Personality string `json:"personality"` // system prompt personality
	MaxTurns    int    `json:"max_turns"`   // max agent loop iterations (default 10)
	MaxHistory  int    `json:"max_history"` // messages to keep in context (default 20)
	Model       string `json:"model"`       // AI model override
	SoulPath    string `json:"soul_path"`   // path to SOUL.md file (optional)
}

// DefaultConfig returns sensible defaults for the agent.
// If soulPath is provided, it will be loaded from file instead of defaults.
func DefaultConfig() Config {
	return Config{
		Name:        "Zara",
		Personality: "You are Zara, a warm, sharp, empathetic AI project manager assistant. You help teams stay organized, track blockers, manage risks, and improve their process. You are concise but warm. Answer in 2-3 sentences unless asked for detail. Use bullet points for lists.",
		MaxTurns:    10,
		MaxHistory:  20,
		Model:       "",
	}
}

// ApplySOUL merges a SOUL definition into this config.
// SOUL values override config defaults.
func (c *Config) ApplySOUL(soul *SOUL) {
	if soul == nil {
		return
	}
	if soul.Name != "" {
		c.Name = soul.Name
	}
	if soul.Personality != "" {
		c.Personality = soul.Personality
	}
	if soul.Backstory != "" {
		// Prepend backstory to personality
		c.Personality = soul.Backstory + "\n\n" + c.Personality
	}
}
