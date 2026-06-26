package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerOKRKPITools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_okr_define",
		mcp.WithDescription("Define an OKR (Objective + Key Results). Links strategic goals to sprint work for outcome tracking."),
		mcp.WithString("objective", mcp.Required(), mcp.Description("Objective statement")),
		mcp.WithString("key_results", mcp.Required(), mcp.Description("Key Results (comma or newline separated)")),
		mcp.WithString("epic_keys", mcp.Description("Linked epic keys (comma separated)")),
		mcp.WithString("measurement", mcp.Description("How to measure: epic_completion, jql_count, manual (default: epic_completion)")),
		mcp.WithString("quarter", mcp.Description("Quarter (default: current)")),
	), h.PMOKRDefine)

	s.AddTool(mcp.NewTool("pm_okr_list",
		mcp.WithDescription("List all defined OKRs with current progress status."),
		mcp.WithString("quarter", mcp.Description("Filter by quarter (default: current)")),
	), h.PMOKRList)

	s.AddTool(mcp.NewTool("pm_kr_link",
		mcp.WithDescription("Link a Key Result to Jira signal (JQL query for auto-progress tracking)."),
		mcp.WithString("kr_id", mcp.Required(), mcp.Description("Key Result identifier")),
		mcp.WithString("jql", mcp.Required(), mcp.Description("JQL that measures this KR")),
		mcp.WithString("metric_type", mcp.Description("count_increase, count_decrease, completion_rate, time_metric")),
		mcp.WithNumber("baseline", mcp.Description("Starting value")),
		mcp.WithNumber("target", mcp.Description("Target value")),
	), h.PMKRLink)

	s.AddTool(mcp.NewTool("pm_kr_progress",
		mcp.WithDescription("Calculate Key Result progress from live Jira data. Auto-queries linked JQL signals."),
		mcp.WithString("quarter", mcp.Description("Quarter (default: current)")),
	), h.PMKRProgress)

	s.AddTool(mcp.NewTool("pm_outcome_review",
		mcp.WithDescription("AI analysis: how did sprint work contribute to OKR progress? Bridges output→outcome gap."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMOutcomeReview)

	s.AddTool(mcp.NewTool("pm_hypothesis",
		mcp.WithDescription("Record a product/process hypothesis to validate. Tracks experiments tied to outcomes."),
		mcp.WithString("hypothesis", mcp.Required(), mcp.Description("If we do X, then Y will happen")),
		mcp.WithString("measure", mcp.Description("How we'll know it worked")),
		mcp.WithString("deadline", mcp.Description("When to evaluate")),
	), h.PMHypothesis)

	s.AddTool(mcp.NewTool("pm_hypothesis_validate",
		mcp.WithDescription("Validate or invalidate a hypothesis with evidence."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Hypothesis ID")),
		mcp.WithString("result", mcp.Required(), mcp.Description("validated, invalidated, or inconclusive")),
		mcp.WithString("evidence", mcp.Description("What data showed")),
	), h.PMHypothesisValidate)

	s.AddTool(mcp.NewTool("pm_kpi_define",
		mcp.WithDescription("Define a KPI to track over time. Simpler than OKR — just a metric with target."),
		mcp.WithString("name", mcp.Required(), mcp.Description("KPI name (e.g. Sprint Velocity, Bug Escape Rate)")),
		mcp.WithNumber("target", mcp.Description("Target value")),
		mcp.WithString("direction", mcp.Description("higher_is_better or lower_is_better (default: higher)")),
	), h.PMKPIDefine)

	s.AddTool(mcp.NewTool("pm_kpi_snapshot",
		mcp.WithDescription("Record a KPI data point. Call regularly (weekly/per-sprint) to build trend."),
		mcp.WithString("name", mcp.Required(), mcp.Description("KPI name")),
		mcp.WithNumber("value", mcp.Required(), mcp.Description("Current value")),
		mcp.WithString("note", mcp.Description("Context for this measurement")),
	), h.PMKPISnapshot)

	s.AddTool(mcp.NewTool("pm_kpi_dashboard",
		mcp.WithDescription("Show all KPIs with trend lines and target comparison. Quick health-of-team view."),
	), h.PMKPIDashboard)
}
