package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerLinearTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_linear_issues",
		mcp.WithDescription("List issues from Linear. Filter by team key or state name."),
		mcp.WithString("team", mcp.Description("Team key (e.g. ENG, PROD)")),
		mcp.WithString("state", mcp.Description("State name filter (e.g. In Progress, Done, Todo)")),
	), h.LinearIssues)

	s.AddTool(mcp.NewTool("pm_linear_cycles",
		mcp.WithDescription("List current and recent Linear cycles (sprints)."),
	), h.LinearCycles)

	s.AddTool(mcp.NewTool("pm_linear_activity",
		mcp.WithDescription("Recent activity feed from Linear: state changes, assignments."),
	), h.LinearActivity)
}
