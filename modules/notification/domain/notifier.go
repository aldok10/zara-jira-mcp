// Package domain provides notification domain interfaces.
package domain

import "context"

// Notifier defines the interface for sending messages to a specific channel.
type Notifier interface {
	// SendMessage sends a message to the given channel/chat.
	SendMessage(ctx context.Context, channelID, title, content string) error

	// Channel returns the name of this notifier (lark, slack, etc.)
	Channel() string
}

// Message represents a structured notification payload.
type Message struct {
	Channel   string
	ChannelID string
	Title     string
	Content   string
	Severity  string // critical, high, medium, low, info
	Audience  string // individual, team, stakeholder, executive
}

// Router routes messages to the appropriate channel.
type Router interface {
	// Route sends a message to the best channel(s) based on severity and audience.
	Route(ctx context.Context, msg Message) error
}

// Gateway is the channel-agnostic interface for sending messages back to users.
type Gateway interface {
	// SendText sends a plain text reply to a channel.
	SendText(ctx context.Context, channelID, text string) error

	// SendMarkdown sends a formatted reply.
	SendMarkdown(ctx context.Context, channelID, title, content string) error

	// Channel returns the name of this gateway (lark, slack, etc.)
	Channel() string
}
