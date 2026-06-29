package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aldok10/zara-jira-mcp/modules/notification/domain"
)

// NotifierAdapter implements domain.Notifier for Telegram.
type NotifierAdapter struct {
	client *Client
}

var _ domain.Notifier = (*NotifierAdapter)(nil)

// NewNotifierAdapter creates a Telegram notifier adapter.
func NewNotifierAdapter(client *Client) *NotifierAdapter {
	return &NotifierAdapter{client: client}
}

// Channel returns the channel name.
func (a *NotifierAdapter) Channel() string { return "telegram" }

// SendMessage sends a message to a Telegram chat.
func (a *NotifierAdapter) SendMessage(ctx context.Context, channelID, title, content string) error {
	var chatID int64
	if channelID != "" {
		id, err := strconv.ParseInt(channelID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid telegram chat_id: %w", err)
		}
		chatID = id
	}

	msg := content
	if title != "" {
		msg = fmt.Sprintf("*%s*\n\n%s", title, content)
	}
	return a.client.SendMessage(ctx, chatID, msg)
}
