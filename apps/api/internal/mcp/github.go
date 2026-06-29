package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	ghmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/github/mcp"
)

func RegisterGitHubTools(s *server.MCPServer, h *ghmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_github_prs",
			mcp.WithDescription("List open PRs with review status, age, and assignees."),
			mcp.WithString("state", mcp.Description("open, closed, all (default: open)")),
			mcp.WithNumber("limit", mcp.Description("Max results (default: 30)")),
		),
		h.ListPRs,
	)
	s.AddTool(
		mcp.NewTool("pm_github_releases",
			mcp.WithDescription("Recent releases/tags. Correlate with sprint delivery."),
			mcp.WithNumber("limit", mcp.Description("Max releases (default: 10)")),
		),
		h.ListReleases,
	)
	s.AddTool(
		mcp.NewTool("pm_github_activity",
			mcp.WithDescription("Repo activity summary: commits, PRs merged, issues closed in last N days."),
			mcp.WithNumber("days", mcp.Description("Number of days to look back (default: 7)")),
		),
		h.GetActivity,
	)
	s.AddTool(
		mcp.NewTool("pm_github_pr_metrics",
			mcp.WithDescription("PR metrics: avg age of open PRs, stale PR count."),
			mcp.WithNumber("stale_days", mcp.Description("Days without update to consider stale (default: 7)")),
		),
		h.GetPRMetrics,
	)
	s.AddTool(
		mcp.NewTool("pm_github_search_branches",
			mcp.WithDescription("Search GitHub branches matching a pattern (e.g. issue key)."),
			mcp.WithString("pattern", mcp.Required(), mcp.Description("Branch pattern to search for")),
		),
		h.SearchBranches,
	)
	s.AddTool(
		mcp.NewTool("pm_github_search_prs_by_branch",
			mcp.WithDescription("Find PRs for a specific branch."),
			mcp.WithString("branch", mcp.Required(), mcp.Description("Branch name")),
		),
		h.SearchPRsByBranch,
	)
	s.AddTool(
		mcp.NewTool("pm_github_create_issue",
			mcp.WithDescription("Create a GitHub issue. Bridge PM planning to developer task tracking."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Issue title")),
			mcp.WithString("body", mcp.Description("Issue body (markdown)")),
			mcp.WithString("labels", mcp.Description("Comma-separated labels")),
			mcp.WithString("assignees", mcp.Description("Comma-separated GitHub usernames")),
			mcp.WithNumber("milestone", mcp.Description("Milestone number")),
		),
		h.CreateIssue,
	)
	s.AddTool(
		mcp.NewTool("pm_github_issues",
			mcp.WithDescription("List GitHub issues. Filter by state and labels."),
			mcp.WithString("state", mcp.Description("open, closed, all (default: open)")),
			mcp.WithString("labels", mcp.Description("Comma-separated label filter")),
			mcp.WithNumber("limit", mcp.Description("Max results (default: 20)")),
		),
		h.ListIssues,
	)
	s.AddTool(
		mcp.NewTool("pm_github_create_milestone",
			mcp.WithDescription("Create a GitHub milestone. Map sprint goals to milestones."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Milestone title")),
			mcp.WithString("description", mcp.Description("Milestone description")),
			mcp.WithString("due_date", mcp.Description("Due date (YYYY-MM-DD)")),
		),
		h.CreateMilestone,
	)
	s.AddTool(
		mcp.NewTool("pm_github_milestones",
			mcp.WithDescription("List GitHub milestones with progress (open/closed issues)."),
			mcp.WithString("state", mcp.Description("open, closed, all (default: open)")),
		),
		h.ListMilestones,
	)
}
