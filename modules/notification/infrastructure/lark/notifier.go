package lark

import (
	"context"

	"github.com/aldok10/zara-jira-mcp/modules/notification/domain"
)

// NotifierAdapter implements domain.Notifier for Lark.
type NotifierAdapter struct {
	client *WebhookClient
}

var _ domain.Notifier = (*NotifierAdapter)(nil)

// NewNotifierAdapter creates a Lark notifier adapter.
func NewNotifierAdapter(client *WebhookClient) *NotifierAdapter {
	return &NotifierAdapter{client: client}
}

// Channel returns the channel name.
func (a *NotifierAdapter) Channel() string { return "lark" }

// SendMessage sends a message via Lark.
func (a *NotifierAdapter) SendMessage(ctx context.Context, channelID, title, content string) error {
	if title != "" {
		return a.client.SendMarkdown(ctx, title, content)
	}
	return a.client.SendText(ctx, content)
}
