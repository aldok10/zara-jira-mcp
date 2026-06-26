package lark

import "context"

// Notifier defines the interface for sending messages to Lark.
type Notifier interface {
	SendText(ctx context.Context, text string) error
	SendMarkdown(ctx context.Context, title, content string) error
}
