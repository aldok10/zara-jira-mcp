package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	clmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/clockify/mcp"
)

func RegisterClockifyTools(s *server.MCPServer, h *clmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_time_entries",
			mcp.WithDescription("Recent time entries from Clockify. Shows who worked on what."),
			mcp.WithInteger("days", mcp.Description("Days to look back (default: 1 for today)")),
		),
		h.TimeEntries,
	)
	s.AddTool(
		mcp.NewTool("pm_time_report",
			mcp.WithDescription("Time tracked by team member per project for a date range."),
			mcp.WithInteger("days", mcp.Description("Days to look back (default: 7)")),
		),
		h.SummaryReport,
	)
}
