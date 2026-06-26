package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerOutcomeTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_impediment_aging",
		mcp.WithDescription("Track how long blockers stay alive. Shows avg resolution time, chronic blockers (>3 days), and aging distribution."),
	), h.PMImpedimentAging)

	s.AddTool(mcp.NewTool("pm_sm_impact",
		mcp.WithDescription("Track Scrum Master's measurable impact: blockers resolved, avg resolution time, action items completed, risks mitigated. Returns SM Impact Score."),
		mcp.WithString("sprint_name", mcp.Description("Sprint name (optional, shows all if empty)")),
	), h.PMSMImpact)

	s.AddTool(mcp.NewTool("pm_stakeholder_pulse",
		mcp.WithDescription("Record stakeholder satisfaction score (1-5). Track how happy stakeholders are over time."),
		mcp.WithString("stakeholder", mcp.Required(), mcp.Description("Stakeholder name")),
		mcp.WithNumber("score", mcp.Required(), mcp.Description("Satisfaction 1-5 (1=unhappy, 5=delighted)")),
		mcp.WithString("sprint_name", mcp.Description("Sprint context")),
		mcp.WithString("feedback", mcp.Description("Qualitative feedback")),
	), h.PMStakeholderPulse)

	s.AddTool(mcp.NewTool("pm_stakeholder_trend",
		mcp.WithDescription("Show stakeholder satisfaction trends over time. Grouped by stakeholder with averages."),
	), h.PMStakeholderTrend)

	s.AddTool(mcp.NewTool("pm_improvement_velocity",
		mcp.WithDescription("Are retro actions getting done faster? Compares action item completion rate across sprints."),
	), h.PMImprovementVelocity)

	s.AddTool(mcp.NewTool("pm_team_autonomy",
		mcp.WithDescription("AI-powered team self-organization assessment. Looks at blocker resolution spread, action ownership, and SM dependency. Returns autonomy score (1-5)."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMTeamAutonomy)

	s.AddTool(mcp.NewTool("pm_outcome_map",
		mcp.WithDescription("Connect sprint work to business outcomes. Map which OKR/objective this sprint serves."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		mcp.WithString("objective", mcp.Required(), mcp.Description("Business objective this sprint serves")),
		mcp.WithString("key_results", mcp.Description("Measurable key results (newline-separated)")),
	), h.PMOutcomeMap)

	s.AddTool(mcp.NewTool("pm_outcome_history",
		mcp.WithDescription("Show OKR/objective alignment history. Which sprints served which business outcomes."),
	), h.PMOutcomeHistory)
}
