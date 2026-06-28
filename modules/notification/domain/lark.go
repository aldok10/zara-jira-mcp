package domain

import "context"

// LarkNotifier defines the interface for sending messages to Lark.
type LarkNotifier interface {
	SendText(ctx context.Context, text string) error
	SendMarkdown(ctx context.Context, title, content string) error
}
