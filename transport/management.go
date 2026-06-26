package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerManagementTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_management_brief",
			mcp.WithDescription("Generate management-level status brief. Tailored by audience: manager (detailed), director/VP (3 lines), PO (what ships/what's blocked). No jargon, outcomes-focused."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithString("audience", mcp.Description("Target: manager (default), director, vp, executive, po, product_owner")),
		),
		h.ManagementBrief,
	)

	s.AddTool(
		mcp.NewTool("pm_dependency_report",
			mcp.WithDescription("Cross-team dependency report. Shows what we're waiting on from other teams, aging, and overdue items. Essential for cross-division communication."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for sprint-blocked items")),
		),
		h.DependencyReport,
	)

	s.AddTool(
		mcp.NewTool("pm_escalation_report",
			mcp.WithDescription("Items needing management attention: long blockers (>3d), unresolved critical risks, sprint-at-risk, high-priority stale items. Use before management sync."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for sprint context")),
		),
		h.EscalationReport,
	)

	s.AddTool(
		mcp.NewTool("pm_resource_utilization",
			mcp.WithDescription("Team workload distribution table: assigned, done, WIP, blocked per member. Flags overloaded and available members. For resource planning discussions."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.ResourceUtilization,
	)

	s.AddTool(
		mcp.NewTool("pm_blocker_aging",
			mcp.WithDescription("Blocker aging report with SLA tracking. Shows how long each blocker is stuck and who owns resolution. For accountability discussions."),
		),
		h.BlockerAgingReport,
	)

	s.AddTool(
		mcp.NewTool("pm_commitment_report",
			mcp.WithDescription("Sprint commitment vs delivery report. Shows current + historical delivery rates. Identifies over-commitment patterns. For sprint review with management."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintCommitmentReport,
	)
}
