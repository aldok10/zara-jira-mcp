package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerOKRKPITools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_okr_define",
		mcp.WithDescription("Define an OKR (Objective + Key Results). Key Results format per line: 'KR title | target_value | unit'."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Objective title")),
		mcp.WithString("level", mcp.Description("company, team, sprint, individual (default: team)")),
		mcp.WithString("owner", mcp.Description("Who owns this OKR")),
		mcp.WithString("cycle", mcp.Description("Period (e.g. Q3-2026)")),
		mcp.WithString("description", mcp.Description("Context")),
		mcp.WithString("key_results", mcp.Description("Key Results, newline-separated: 'title | target | unit'")),
	), h.PMOKRDefine)

	s.AddTool(mcp.NewTool("pm_okr_list",
		mcp.WithDescription("List OKRs with Key Result progress."),
		mcp.WithString("status", mcp.Description("active (default), achieved, abandoned")),
	), h.PMOKRList)

	s.AddTool(mcp.NewTool("pm_kr_link",
		mcp.WithDescription("Link Jira issues to a Key Result for auto progress tracking."),
		mcp.WithNumber("kr_id", mcp.Required(), mcp.Description("Key Result ID")),
		mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated Jira keys")),
		mcp.WithString("link_type", mcp.Description("contributes (default), blocks, measures")),
	), h.PMKRLink)

	s.AddTool(mcp.NewTool("pm_kr_progress",
		mcp.WithDescription("Calculate KR progress from linked Jira issues (auto-updates based on Done status)."),
	), h.PMKRProgress)

	s.AddTool(mcp.NewTool("pm_outcome_review",
		mcp.WithDescription("AI assessment: did sprint work move the OKR needle?"),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMOutcomeReview)

	s.AddTool(mcp.NewTool("pm_kpi_define",
		mcp.WithDescription("Define a KPI with target and warning/danger thresholds."),
		mcp.WithString("name", mcp.Required(), mcp.Description("KPI name")),
		mcp.WithString("description", mcp.Description("What it measures")),
		mcp.WithString("formula", mcp.Description("How to calculate")),
		mcp.WithString("unit", mcp.Description("Unit (default: %)")),
		mcp.WithNumber("target_value", mcp.Description("Target (green)")),
		mcp.WithNumber("warning_threshold", mcp.Description("Warning level (yellow)")),
		mcp.WithNumber("danger_threshold", mcp.Description("Danger level (red)")),
	), h.PMKPIDefine)

	s.AddTool(mcp.NewTool("pm_kpi_snapshot",
		mcp.WithDescription("Record a KPI measurement. Shows status vs thresholds."),
		mcp.WithNumber("kpi_id", mcp.Description("KPI ID")),
		mcp.WithString("kpi_name", mcp.Description("KPI name (alternative to kpi_id)")),
		mcp.WithNumber("value", mcp.Required(), mcp.Description("Measured value")),
		mcp.WithString("sprint_name", mcp.Description("Sprint context")),
		mcp.WithString("notes", mcp.Description("Context")),
	), h.PMKPISnapshot)

	s.AddTool(mcp.NewTool("pm_kpi_dashboard",
		mcp.WithDescription("All KPIs with values, trends (^ v =), and status."),
	), h.PMKPIDashboard)

	s.AddTool(mcp.NewTool("pm_goal_hit_rate",
		mcp.WithDescription("Sprint goal success rate over time. Industry avg is 52%."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		mcp.WithNumber("limit", mcp.Description("Sprints to analyze (default: 10)")),
	), h.PMGoalHitRate)

	s.AddTool(mcp.NewTool("pm_okr_health",
		mcp.WithDescription("OKR risk assessment: progress vs time elapsed. Flags AT RISK objectives."),
	), h.PMOKRHealth)
}
