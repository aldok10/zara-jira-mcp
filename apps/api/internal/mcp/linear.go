package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	lmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/linear/mcp"
)

func RegisterLinearTools(s *server.MCPServer, h *lmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_linear_issues",
			mcp.WithDescription("List issues from Linear. Filter by team key or state name."),
			mcp.WithString("team", mcp.Description("Team key (e.g. ENG, PROD)")),
			mcp.WithString("state", mcp.Description("State name filter (e.g. In Progress, Done, Todo)")),
		),
		h.ListIssues,
	)
	s.AddTool(
		mcp.NewTool("pm_linear_cycles",
			mcp.WithDescription("List current and recent Linear cycles (sprints)."),
		),
		h.ListCycles,
	)
	s.AddTool(
		mcp.NewTool("pm_linear_activity",
			mcp.WithDescription("Recent activity feed from Linear: state changes, assignments."),
		),
		h.RecentActivity,
	)
}
