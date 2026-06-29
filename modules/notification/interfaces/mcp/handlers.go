package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/modules/notification/domain"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

type Handlers struct {
	Notifiers map[string]domain.Notifier
	Router    domain.Router
	Error     *mcputil.ErrorHandler
}

func NewHandlers(notifiers map[string]domain.Notifier, router domain.Router) *Handlers {
	return &Handlers{
		Notifiers: notifiers,
		Router:    router,
		Error:     mcputil.NewErrorHandler(nil),
	}
}

// availableChannels returns a comma-separated list of configured channels.
func (h *Handlers) availableChannels() string {
	channels := make([]string, 0, len(h.Notifiers))
	for ch := range h.Notifiers {
		channels = append(channels, ch)
	}
	return strings.Join(channels, ", ")
}

// NotifyLark sends a markdown message to the configured Lark group.
func (h *Handlers) NotifyLark(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content, err := req.RequireString("content")
	if err != nil {
		return mcputil.ErrInvalid("content parameter is required"), nil
	}
	title := req.GetString("title", "PM Update")

	notifier, ok := h.Notifiers["lark"]
	if !ok {
		return mcputil.ErrorResult("Lark not configured. Set LARK_WEBHOOK_URL or LARK_APP_ID + LARK_APP_SECRET."), nil
	}

	if err := notifier.SendMessage(ctx, "", title, content); err != nil {
		return h.Error.Wrap("send lark message", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Message sent to Lark: %s", title)), nil
}

// NotifySlack sends a message to Slack.
func (h *Handlers) NotifySlack(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content, err := req.RequireString("content")
	if err != nil {
		return mcputil.ErrInvalid("content parameter is required"), nil
	}
	title := req.GetString("title", "PM Update")
	channel := req.GetString("channel", "")

	notifier, ok := h.Notifiers["slack"]
	if !ok {
		return mcputil.ErrorResult("Slack not configured. Set SLACK_BOT_TOKEN or SLACK_WEBHOOK_URL."), nil
	}

	if err := notifier.SendMessage(ctx, channel, title, content); err != nil {
		return h.Error.Wrap("send slack message", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Message sent to Slack: %s", title)), nil
}

// NotifyDiscord sends a message to Discord.
func (h *Handlers) NotifyDiscord(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content, err := req.RequireString("content")
	if err != nil {
		return mcputil.ErrInvalid("content parameter is required"), nil
	}
	title := req.GetString("title", "PM Update")
	channel := req.GetString("channel", "")

	notifier, ok := h.Notifiers["discord"]
	if !ok {
		return mcputil.ErrorResult("Discord not configured. Set DISCORD_BOT_TOKEN."), nil
	}

	if err := notifier.SendMessage(ctx, channel, title, content); err != nil {
		return h.Error.Wrap("send discord message", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Message sent to Discord: %s", title)), nil
}

// NotifyTelegram sends a message to Telegram.
func (h *Handlers) NotifyTelegram(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content, err := req.RequireString("content")
	if err != nil {
		return mcputil.ErrInvalid("content parameter is required"), nil
	}
	title := req.GetString("title", "PM Update")
	chatID := req.GetString("chat_id", "")

	notifier, ok := h.Notifiers["telegram"]
	if !ok {
		return mcputil.ErrorResult("Telegram not configured. Set TELEGRAM_BOT_TOKEN."), nil
	}

	if err := notifier.SendMessage(ctx, chatID, title, content); err != nil {
		return h.Error.Wrap("send telegram message", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Message sent to Telegram: %s", title)), nil
}

// NotifyRouted sends a message to the optimal channel based on severity and audience.
func (h *Handlers) NotifyRouted(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content, err := req.RequireString("content")
	if err != nil {
		return mcputil.ErrInvalid("content parameter is required"), nil
	}
	severity := req.GetString("severity", "medium")
	audience := req.GetString("audience", "team")
	title := req.GetString("title", "PM Update")

	if h.Router != nil {
		msg := domain.Message{
			ChannelID: "",
			Title:     title,
			Content:   content,
			Severity:  severity,
			Audience:  audience,
		}
		if err := h.Router.Route(ctx, msg); err != nil {
			return h.Error.Wrap("route notification", err), nil
		}
		return mcputil.TextResult(fmt.Sprintf("Notification routed: %s (severity: %s, audience: %s)", title, severity, audience)), nil
	}

	// Fallback: try Lark first, then Slack
	if n, ok := h.Notifiers["lark"]; ok {
		if err := n.SendMessage(ctx, "", title, content); err == nil {
			return mcputil.TextResult(fmt.Sprintf("Notification sent to Lark: %s", title)), nil
		}
	}
	if n, ok := h.Notifiers["slack"]; ok {
		if err := n.SendMessage(ctx, "", title, content); err == nil {
			return mcputil.TextResult(fmt.Sprintf("Notification sent to Slack: %s", title)), nil
		}
	}

	if len(h.Notifiers) == 0 {
		return mcputil.ErrorResult("No notification channels configured. Available: " + h.availableChannels()), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Notification recorded — no suitable channel found for severity=%s, audience=%s", severity, audience)), nil
}
