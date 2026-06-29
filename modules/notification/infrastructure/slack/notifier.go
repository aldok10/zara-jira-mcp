package slack

import (
	"context"

	"github.com/aldok10/zara-jira-mcp/modules/notification/domain"
)

// NotifierAdapter implements domain.Notifier for Slack.
type NotifierAdapter struct {
	client *Client
}

var _ domain.Notifier = (*NotifierAdapter)(nil)

// NewNotifierAdapter creates a Slack notifier adapter.
func NewNotifierAdapter(client *Client) *NotifierAdapter {
	return &NotifierAdapter{client: client}
}

// Channel returns the channel name.
func (a *NotifierAdapter) Channel() string { return "slack" }

// SendMessage sends a message to a Slack channel.
func (a *NotifierAdapter) SendMessage(ctx context.Context, channelID, title, content string) error {
	if title != "" {
		return a.client.SendRichMessage(ctx, channelID, title, content)
	}
	return a.client.SendMessage(ctx, channelID, content)
}
