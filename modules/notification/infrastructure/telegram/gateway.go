package telegram

import (
	"context"
	"fmt"
	"strconv"
)

// GatewayAdapter wraps the Telegram client as an agent domain.Gateway.
type GatewayAdapter struct {
	client *Client
}

// NewGatewayAdapter creates a Telegram gateway adapter.
func NewGatewayAdapter(client *Client) *GatewayAdapter {
	return &GatewayAdapter{client: client}
}

// Channel returns the channel name for this gateway.
func (g *GatewayAdapter) Channel() string { return "telegram" }

// SendText sends a plain text message to a Telegram chat.
func (g *GatewayAdapter) SendText(ctx context.Context, channelID, text string) error {
	chatID := parseChatID(channelID)
	return g.client.SendMessage(ctx, chatID, text)
}

// SendMarkdown sends a formatted message to a Telegram chat.
func (g *GatewayAdapter) SendMarkdown(ctx context.Context, channelID, title, content string) error {
	chatID := parseChatID(channelID)
	msg := fmt.Sprintf("*%s*\n\n%s", title, content)
	return g.client.SendMessage(ctx, chatID, msg)
}

func parseChatID(s string) int64 {
	if s == "" {
		return 0
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return id
}
