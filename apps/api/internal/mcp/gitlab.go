package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	glmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/gitlab/mcp"
)

func RegisterGitLabTools(s *server.MCPServer, h *glmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_gitlab_issues",
			mcp.WithDescription("List GitLab issues. Filter by state and labels."),
			mcp.WithString("state", mcp.Description("opened, closed, all (default: opened)")),
			mcp.WithString("labels", mcp.Description("Comma-separated label filter")),
			mcp.WithNumber("limit", mcp.Description("Max results (default: 20)")),
		),
		h.ListIssues,
	)
	s.AddTool(
		mcp.NewTool("pm_gitlab_create_issue",
			mcp.WithDescription("Create a GitLab issue. Bridge PM planning to developer work."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Issue title")),
			mcp.WithString("description", mcp.Description("Issue description (markdown)")),
			mcp.WithString("labels", mcp.Description("Comma-separated labels")),
			mcp.WithNumber("assignee_id", mcp.Description("Assignee user ID")),
			mcp.WithNumber("milestone_id", mcp.Description("Milestone ID")),
		),
		h.CreateIssue,
	)
	s.AddTool(
		mcp.NewTool("pm_gitlab_merge_requests",
			mcp.WithDescription("List GitLab merge requests."),
			mcp.WithString("state", mcp.Description("opened, closed, merged, all (default: opened)")),
			mcp.WithNumber("limit", mcp.Description("Max results (default: 20)")),
		),
		h.ListMRs,
	)
	s.AddTool(
		mcp.NewTool("pm_gitlab_milestones",
			mcp.WithDescription("List GitLab milestones."),
			mcp.WithString("state", mcp.Description("active, closed (default: active)")),
		),
		h.ListMilestones,
	)
	s.AddTool(
		mcp.NewTool("pm_gitlab_create_milestone",
			mcp.WithDescription("Create a GitLab milestone. Map sprint goals to milestones."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Milestone title")),
			mcp.WithString("description", mcp.Description("Description")),
			mcp.WithString("due_date", mcp.Description("Due date (YYYY-MM-DD)")),
		),
		h.CreateMilestone,
	)
	s.AddTool(
		mcp.NewTool("pm_gitlab_search_branches",
			mcp.WithDescription("Search GitLab branches matching a pattern (e.g. issue key)."),
			mcp.WithString("pattern", mcp.Required(), mcp.Description("Branch pattern to search for")),
		),
		h.SearchBranches,
	)
	s.AddTool(
		mcp.NewTool("pm_gitlab_search_mrs_by_branch",
			mcp.WithDescription("Find merge requests for a specific branch."),
			mcp.WithString("branch", mcp.Required(), mcp.Description("Branch name")),
		),
		h.SearchMRsByBranch,
	)
	s.AddTool(
		mcp.NewTool("pm_gitlab_read_file",
			mcp.WithDescription("Read a file from the GitLab repo."),
			mcp.WithString("path", mcp.Required(), mcp.Description("File path")),
			mcp.WithString("ref", mcp.Description("Branch/tag (default: main)")),
		),
		h.GetFileContent,
	)
	s.AddTool(
		mcp.NewTool("pm_gitlab_list_files",
			mcp.WithDescription("List files/directories in a GitLab repo path."),
			mcp.WithString("path", mcp.Description("Directory path (default: root)")),
			mcp.WithString("ref", mcp.Description("Branch/tag (default: main)")),
		),
		h.ListFiles,
	)
}
