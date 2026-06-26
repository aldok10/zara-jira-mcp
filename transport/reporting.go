package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerReportingTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("report_to_po",
		mcp.WithDescription("PO briefing: value delivered, blocked items, goal status."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.ReportToPO)



	s.AddTool(mcp.NewTool("report_delivery_confidence",
		mcp.WithDescription("Delivery confidence (GREEN/AMBER/RED). Data-backed sprint assessment."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.DeliveryConfidenceReport)

	s.AddTool(mcp.NewTool("report_resource_planning",
		mcp.WithDescription("Capacity planning: throughput, utilization, workload per person."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.ResourcePlanningReport)
}
