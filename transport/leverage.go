package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerLeverageTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_tech_debt_ratio",
		mcp.WithDescription("Calculate tech debt ratio: bugs/debt items vs feature items. Flags health status and recommends action when >20%."),
		mcp.WithString("project", mcp.Description("Project key (all projects if empty)")),
		mcp.WithNumber("sprints", mcp.Description("Number of sprints for trend (default: 3)")),
	), h.TechDebtRatio)

	s.AddTool(mcp.NewTool("pm_priority_churn",
		mcp.WithDescription("Detect priority instability — issues with priority changes in last N days. High churn = team burnout risk (DORA 2024)."),
		mcp.WithString("project", mcp.Description("Project key (all projects if empty)")),
		mcp.WithNumber("days", mcp.Description("Lookback window in days (default: 14)")),
	), h.PriorityChurn)

	s.AddTool(mcp.NewTool("jira_trace_branch",
		mcp.WithDescription("Trace a Jira ticket to its branch in GitHub/GitLab. Shows if branch exists, PRs/MRs status, and whether it's been merged to target branch. Use to verify ticket implementation status."),
		mcp.WithString("key", mcp.Required(), mcp.Description("Jira issue key (e.g. SIT-3658)")),
	), h.JiraTraceBranch)

	s.AddTool(mcp.NewTool("pm_incident_summary",
		mcp.WithDescription("Summarize incident impact: counts by status/urgency, average resolution time."),
	), h.IncidentImpact)

	s.AddTool(mcp.NewTool("pm_sprint_forecast_simple",
		mcp.WithDescription("Simple sprint forecast: will we finish? Based on current burn rate and remaining items."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.SprintForecastSimple)

	s.AddTool(mcp.NewTool("pm_backlog_health",
		mcp.WithDescription("Backlog quality check: finds stale items, recommends grooming actions."),
		mcp.WithString("project", mcp.Description("Project key")),
		mcp.WithNumber("days", mcp.Description("Days to consider stale (default: 90)")),
	), h.BacklogHealthCheck)
}
