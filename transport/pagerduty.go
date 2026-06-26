package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPagerDutyTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_incidents",
		mcp.WithDescription("List recent PagerDuty incidents. Shows severity, status, assignee, service."),
		mcp.WithString("status", mcp.Description("Filter: triggered, acknowledged, resolved (default: all)")),
	), h.PagerDutyIncidents)

	s.AddTool(mcp.NewTool("pm_incident_summary",
		mcp.WithDescription("Summarize incident impact: counts by status/urgency, average resolution time."),
	), h.PagerDutyIncidentSummary)

	s.AddTool(mcp.NewTool("pm_oncall",
		mcp.WithDescription("Who is on call right now across all schedules."),
	), h.PagerDutyOnCall)
}
