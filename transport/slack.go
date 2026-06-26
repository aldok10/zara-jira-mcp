package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSlackTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("slack_send",
			mcp.WithDescription("Send a message to a Slack channel. Supports plain text or rich format with title."),
			mcp.WithString("text", mcp.Required(), mcp.Description("Message text (supports Slack mrkdwn format)")),
			mcp.WithString("channel", mcp.Description("Channel ID or name (uses default if empty)")),
			mcp.WithString("title", mcp.Description("Optional title for rich message card")),
		),
		h.SlackSendMessage,
	)

	s.AddTool(
		mcp.NewTool("slack_channels",
			mcp.WithDescription("List accessible Slack channels with member count and topic."),
		),
		h.SlackListChannels,
	)

	s.AddTool(
		mcp.NewTool("slack_history",
			mcp.WithDescription("Get recent messages from a Slack channel. Useful for context gathering."),
			mcp.WithString("channel", mcp.Required(), mcp.Description("Channel ID")),
			mcp.WithNumber("limit", mcp.Description("Number of messages (default 20)")),
		),
		h.SlackChannelHistory,
	)

	s.AddTool(
		mcp.NewTool("slack_notify_team",
			mcp.WithDescription("Send a formatted PM notification to the team Slack channel. Use for sprint updates, standup summaries, alerts."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Notification content in mrkdwn")),
			mcp.WithString("title", mcp.Description("Notification title (default: PM Update)")),
			mcp.WithString("channel", mcp.Description("Target channel (uses default if empty)")),
		),
		h.SlackNotifyTeam,
	)
}
