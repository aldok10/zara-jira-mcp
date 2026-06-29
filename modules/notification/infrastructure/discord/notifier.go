package discord

import (
	"context"

	"github.com/aldok10/zara-jira-mcp/modules/notification/domain"
)

// NotifierAdapter implements domain.Notifier for Discord.
type NotifierAdapter struct {
	client *Client
}

var _ domain.Notifier = (*NotifierAdapter)(nil)

// NewNotifierAdapter creates a Discord notifier adapter.
func NewNotifierAdapter(client *Client) *NotifierAdapter {
	return &NotifierAdapter{client: client}
}

// Channel returns the channel name.
func (a *NotifierAdapter) Channel() string { return "discord" }

// SendMessage sends a message to a Discord channel.
func (a *NotifierAdapter) SendMessage(ctx context.Context, channelID, title, content string) error {
	if title != "" {
		return a.client.SendEmbed(ctx, channelID, title, content, 0x00AFFF)
	}
	return a.client.SendMessage(ctx, channelID, content)
}
