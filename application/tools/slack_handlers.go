package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// SlackSendMessage sends a message to a Slack channel.
func (h *Handlers) SlackSendMessage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.CheckNotificationBudget() {
		return errorResult("Daily notification budget exceeded (5/day)."), nil
	}
	if h.Slack == nil || !h.Slack.Available() {
		return errorResult("Slack not configured. Set SLACK_BOT_TOKEN or SLACK_WEBHOOK_URL."), nil
	}
	text, err := req.RequireString("text")
	if err != nil {
		return errorResult("text parameter is required"), nil
	}
	channel := req.GetString("channel", "")
	title := req.GetString("title", "")

	if title != "" {
		if err := h.Slack.SendRichMessage(ctx, channel, title, text); err != nil {
			return errorResult("Slack send failed: " + err.Error()), nil
		}
	} else {
		if err := h.Slack.SendMessage(ctx, channel, text); err != nil {
			return errorResult("Slack send failed: " + err.Error()), nil
		}
	}
	h.LogNotification("slack", "medium", title)
	return textResult("Message sent to Slack."), nil
}

// SlackListChannels lists accessible Slack channels.
func (h *Handlers) SlackListChannels(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Slack == nil || !h.Slack.Available() {
		return errorResult("Slack not configured."), nil
	}
	channels, err := h.Slack.ListChannels(ctx)
	if err != nil {
		return errorResult("Failed to list channels: " + err.Error()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Slack channels (%d):\n\n", len(channels)))
	for _, ch := range channels {
		topic := ch.Topic
		if len(topic) > 50 {
			topic = topic[:50] + "..."
		}
		sb.WriteString(fmt.Sprintf("- #%s (%s) [%d members] %s\n", ch.Name, ch.ID, ch.MemberCount, topic))
	}
	return textResult(sb.String()), nil
}

// SlackChannelHistory gets recent messages from a channel.
func (h *Handlers) SlackChannelHistory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Slack == nil || !h.Slack.Available() {
		return errorResult("Slack not configured."), nil
	}
	channel, err := req.RequireString("channel")
	if err != nil {
		return errorResult("channel parameter is required"), nil
	}
	limit := req.GetInt("limit", 20)

	messages, err := h.Slack.GetChannelHistory(ctx, channel, limit)
	if err != nil {
		return errorResult("Failed to get history: " + err.Error()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Recent messages in channel (%d):\n\n", len(messages)))
	for _, m := range messages {
		text := m.Text
		if len(text) > 200 {
			text = text[:200] + "..."
		}
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", m.Timestamp, m.User, text))
	}
	return textResult(sb.String()), nil
}

// SlackNotifyTeam sends a formatted PM notification to the team channel.
func (h *Handlers) SlackNotifyTeam(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.CheckNotificationBudget() {
		return errorResult("Daily notification budget exceeded (5/day)."), nil
	}
	if h.Slack == nil || !h.Slack.Available() {
		return errorResult("Slack not configured."), nil
	}
	content, err := req.RequireString("content")
	if err != nil {
		return errorResult("content parameter is required"), nil
	}
	title := req.GetString("title", "PM Update")
	channel := req.GetString("channel", "")

	if err := h.Slack.SendRichMessage(ctx, channel, title, content); err != nil {
		return errorResult("Slack notify failed: " + err.Error()), nil
	}
	h.LogNotification("slack_team", "high", title)
	return textResult("Team notified on Slack."), nil
}
