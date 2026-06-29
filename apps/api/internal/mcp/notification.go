package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	nmcp "github.com/aldok10/zara-jira-mcp/modules/notification/interfaces/mcp"
)

func RegisterNotificationTools(s *server.MCPServer, h *nmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_notify_lark",
			mcp.WithDescription("Send a markdown message to the configured Lark group."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message content in markdown")),
			mcp.WithString("title", mcp.Description("Card title (default: PM Update)")),
		),
		h.NotifyLark,
	)
	s.AddTool(
		mcp.NewTool("jira_notify_slack",
			mcp.WithDescription("Send a message to a Slack channel."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message content in mrkdwn")),
			mcp.WithString("title", mcp.Description("Message title (default: PM Update)")),
			mcp.WithString("channel", mcp.Description("Channel ID or name (uses default if empty)")),
		),
		h.NotifySlack,
	)
	s.AddTool(
		mcp.NewTool("jira_notify_discord",
			mcp.WithDescription("Send a message to a Discord channel."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message content")),
			mcp.WithString("title", mcp.Description("Embed title (sends as rich embed)")),
			mcp.WithString("channel", mcp.Description("Channel ID (uses default if empty)")),
		),
		h.NotifyDiscord,
	)
	s.AddTool(
		mcp.NewTool("jira_notify_telegram",
			mcp.WithDescription("Send a message to a Telegram chat."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message text")),
			mcp.WithString("title", mcp.Description("Bold title prepended to message")),
			mcp.WithString("chat_id", mcp.Description("Chat ID (uses default if empty)")),
		),
		h.NotifyTelegram,
	)
	s.AddTool(
		mcp.NewTool("notify_routed",
			mcp.WithDescription("Smart notification routing. Auto-sends to best channel based on severity and audience."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Notification content")),
			mcp.WithString("severity", mcp.Description("critical, high, medium, low, info (default: medium)")),
			mcp.WithString("audience", mcp.Description("individual, team, stakeholder, executive (default: team)")),
			mcp.WithString("title", mcp.Description("Notification title")),
		),
		h.NotifyRouted,
	)
}
