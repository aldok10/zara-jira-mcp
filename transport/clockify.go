package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerClockifyTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_time_report",
		mcp.WithDescription("Time tracked by team member per project for a date range."),
		mcp.WithNumber("days", mcp.Description("Days to look back (default: 7)")),
	), h.ClockifyTimeReport)

	s.AddTool(mcp.NewTool("pm_time_entries",
		mcp.WithDescription("Recent time entries showing who worked on what."),
		mcp.WithNumber("days", mcp.Description("Days to look back (default: 1 for today)")),
	), h.ClockifyTimeEntries)
}
