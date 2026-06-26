package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerManagementTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_dependency_report",
			mcp.WithDescription("Cross-team dependency report: waiting items, aging, overdue."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for sprint-blocked items")),
		),
		h.DependencyReport,
	)

	s.AddTool(
		mcp.NewTool("pm_escalation_report",
			mcp.WithDescription("Items needing management attention: long blockers, critical risks."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for sprint context")),
		),
		h.EscalationReport,
	)

	s.AddTool(
		mcp.NewTool("pm_resource_utilization",
			mcp.WithDescription("Workload per member: assigned, done, WIP, blocked. Flags overloaded."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.ResourceUtilization,
	)
}
