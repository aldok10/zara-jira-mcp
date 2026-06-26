package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerRoutingTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("notify_routed",
			mcp.WithDescription("Smart notification routing. Automatically sends to the optimal channel(s) based on severity and audience. Use this instead of platform-specific tools when unsure where to send."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Notification content")),
			mcp.WithString("severity", mcp.Description("critical, high, medium, low, info (default: medium)")),
			mcp.WithString("audience", mcp.Description("individual, team, stakeholder, executive (default: team)")),
			mcp.WithString("title", mcp.Description("Notification title")),
		),
		h.NotifyRouted,
	)

	s.AddTool(
		mcp.NewTool("daily_digest",
			mcp.WithDescription("Generate and send daily digest: active blockers, pending actions, open risks, stale items. Sends to primary notification channel. Call at start of day."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for sprint context")),
		),
		h.DailyDigest,
	)
}
