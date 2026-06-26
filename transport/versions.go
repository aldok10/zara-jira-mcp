package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerVersionTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_attachments",
			mcp.WithDescription("List attachments on a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		),
		h.GetAttachments,
	)

	s.AddTool(
		mcp.NewTool("jira_versions",
			mcp.WithDescription("List project versions/releases."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project key")),
		),
		h.GetVersions,
	)

	s.AddTool(
		mcp.NewTool("jira_version_create",
			mcp.WithDescription("Create a new project version for release tracking."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project key")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Version name (e.g. v1.2.0)")),
			mcp.WithString("description", mcp.Description("Version description")),
		),
		h.CreateVersion,
	)

	s.AddTool(
		mcp.NewTool("jira_version_release",
			mcp.WithDescription("Mark a version as released."),
			mcp.WithString("version_id", mcp.Required(), mcp.Description("Version ID (from jira_versions)")),
		),
		h.ReleaseVersion,
	)

	s.AddTool(
		mcp.NewTool("jira_components",
			mcp.WithDescription("List project components with leads."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project key")),
		),
		h.GetComponents,
	)

	s.AddTool(
		mcp.NewTool("jira_fields",
			mcp.WithDescription("List all available Jira fields (system + custom). Useful for finding custom field IDs."),
			mcp.WithBoolean("custom_only", mcp.Description("Show only custom fields (default: false)")),
		),
		h.GetFields,
	)

	s.AddTool(
		mcp.NewTool("pm_tech_debt_ratio",
			mcp.WithDescription("Calculate tech debt ratio: bugs/debt items vs feature items. Flags health status and recommends action when >20%."),
			mcp.WithString("project", mcp.Description("Project key (all projects if empty)")),
			mcp.WithNumber("sprints", mcp.Description("Number of sprints for trend (default: 3)")),
		),
		h.TechDebtRatio,
	)

	s.AddTool(
		mcp.NewTool("pm_priority_churn",
			mcp.WithDescription("Detect priority instability — issues with priority changes in last N days. High churn = team burnout risk (DORA 2024)."),
			mcp.WithString("project", mcp.Description("Project key (all projects if empty)")),
			mcp.WithNumber("days", mcp.Description("Lookback window in days (default: 14)")),
		),
		h.PriorityChurn,
	)
}
