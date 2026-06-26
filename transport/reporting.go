package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerReportingTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("report_to_po",
		mcp.WithDescription("Generate PO briefing: value delivered, blocked items needing decisions, sprint goal status. Written for Product Owner, not engineers."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.ReportToPO)

	s.AddTool(mcp.NewTool("report_escalation_brief",
		mcp.WithDescription("[DEPRECATED: use pm_escalation_draft or pm_escalation_report] Structured impediment escalation for management."),
		mcp.WithNumber("board_id", mcp.Description("Board ID for sprint impact context")),
	), h.EscalationBrief)

	s.AddTool(mcp.NewTool("report_cross_team_deps",
		mcp.WithDescription("[DEPRECATED: use pm_dependencies] Cross-team dependency status report."),
	), h.CrossTeamDependencyReport)

	s.AddTool(mcp.NewTool("report_delivery_confidence",
		mcp.WithDescription("Delivery confidence report (GREEN/AMBER/RED) with data-backed assessment. Shows management whether sprint goal will be hit."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.DeliveryConfidenceReport)

	s.AddTool(mcp.NewTool("report_resource_planning",
		mcp.WithDescription("Team capacity and resource planning report for management. Throughput trends, utilization, workload per person."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.ResourcePlanningReport)
}
