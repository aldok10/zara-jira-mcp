package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerRoutingTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("notify_routed",
			mcp.WithDescription("Smart notification routing. Automatically sends to the optimal channel(s) based on severity and audience."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Notification content")),
			mcp.WithString("severity", mcp.Description("critical, high, medium, low, info (default: medium)")),
			mcp.WithString("audience", mcp.Description("individual, team, stakeholder, executive (default: team)")),
			mcp.WithString("title", mcp.Description("Notification title")),
		),
		h.NotifyRouted,
	)
}
