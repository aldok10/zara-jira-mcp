package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	pdmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/pagerduty/mcp"
)

func RegisterPagerDutyTools(s *server.MCPServer, h *pdmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_incidents",
			mcp.WithDescription("List recent PagerDuty incidents. Shows severity, status, assignee, service."),
			mcp.WithString("status", mcp.Description("Filter: triggered, acknowledged, resolved (default: all)")),
		),
		h.ListIncidents,
	)
	s.AddTool(
		mcp.NewTool("pm_oncall",
			mcp.WithDescription("Who is on call right now across all schedules."),
		),
		h.GetOnCalls,
	)
}
