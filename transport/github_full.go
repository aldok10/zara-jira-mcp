package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerGitHubFullTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_github_prs",
		mcp.WithDescription("List open PRs with review status, age, and assignees."),
		mcp.WithString("state", mcp.Description("PR state: open, closed, all (default: open)")),
		mcp.WithNumber("limit", mcp.Description("Max results (default: 30)")),
	), h.GitHubPRs)

	s.AddTool(mcp.NewTool("pm_github_pr_metrics",
		mcp.WithDescription("PR metrics: avg age of open PRs, stale PR count."),
		mcp.WithNumber("stale_days", mcp.Description("Days without update to consider stale (default: 7)")),
	), h.GitHubPRMetrics)

	s.AddTool(mcp.NewTool("pm_github_releases",
		mcp.WithDescription("Recent releases/tags. Correlate with sprint delivery."),
		mcp.WithNumber("limit", mcp.Description("Max releases (default: 10)")),
	), h.GitHubReleases)

	s.AddTool(mcp.NewTool("pm_github_activity",
		mcp.WithDescription("Repo activity summary: commits, PRs merged, issues closed in last N days."),
		mcp.WithNumber("days", mcp.Description("Number of days to look back (default: 7)")),
	), h.GitHubActivity)
}
