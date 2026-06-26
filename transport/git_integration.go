package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerGitIntegrationTools(s *server.MCPServer, h *tools.Handlers) {
	// GitHub Issues & Milestones
	s.AddTool(mcp.NewTool("github_create_issue",
		mcp.WithDescription("Create a GitHub issue. Bridge PM planning to developer task tracking."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Issue title")),
		mcp.WithString("body", mcp.Description("Issue body (markdown)")),
		mcp.WithString("labels", mcp.Description("Comma-separated labels")),
		mcp.WithString("assignees", mcp.Description("Comma-separated GitHub usernames")),
		mcp.WithNumber("milestone", mcp.Description("Milestone number")),
	), h.GitHubCreateIssue)

	s.AddTool(mcp.NewTool("github_issues",
		mcp.WithDescription("List GitHub issues. Filter by state and labels."),
		mcp.WithString("state", mcp.Description("open, closed, all (default: open)")),
		mcp.WithString("labels", mcp.Description("Comma-separated label filter")),
		mcp.WithNumber("limit", mcp.Description("Max results (default: 20)")),
	), h.GitHubListIssues)

	s.AddTool(mcp.NewTool("github_create_milestone",
		mcp.WithDescription("Create a GitHub milestone. Map sprint goals to milestones."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Milestone title")),
		mcp.WithString("description", mcp.Description("Milestone description")),
		mcp.WithString("due_date", mcp.Description("Due date (YYYY-MM-DD)")),
	), h.GitHubCreateMilestone)

	s.AddTool(mcp.NewTool("github_milestones",
		mcp.WithDescription("List GitHub milestones with progress (open/closed issues)."),
		mcp.WithString("state", mcp.Description("open, closed, all (default: open)")),
	), h.GitHubListMilestones)

	s.AddTool(mcp.NewTool("github_read_file",
		mcp.WithDescription("Read a file from the GitHub repo. Scan code, configs, docs."),
		mcp.WithString("path", mcp.Required(), mcp.Description("File path (e.g. 'src/main.go')")),
		mcp.WithString("ref", mcp.Description("Branch/tag/commit (default: main)")),
	), h.GitHubReadFile)

	s.AddTool(mcp.NewTool("github_list_files",
		mcp.WithDescription("List files/directories in a GitHub repo path."),
		mcp.WithString("path", mcp.Description("Directory path (default: root)")),
		mcp.WithString("ref", mcp.Description("Branch/tag (default: main)")),
	), h.GitHubListFiles)

	// GitLab Issues & Milestones
	s.AddTool(mcp.NewTool("gitlab_create_issue",
		mcp.WithDescription("Create a GitLab issue. Bridge PM planning to developer work."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Issue title")),
		mcp.WithString("description", mcp.Description("Issue description (markdown)")),
		mcp.WithString("labels", mcp.Description("Comma-separated labels")),
		mcp.WithNumber("assignee_id", mcp.Description("Assignee user ID")),
		mcp.WithNumber("milestone_id", mcp.Description("Milestone ID")),
	), h.GitLabCreateIssue)

	s.AddTool(mcp.NewTool("gitlab_issues",
		mcp.WithDescription("List GitLab issues. Filter by state and labels."),
		mcp.WithString("state", mcp.Description("opened, closed, all (default: opened)")),
		mcp.WithString("labels", mcp.Description("Comma-separated label filter")),
		mcp.WithNumber("limit", mcp.Description("Max results (default: 20)")),
	), h.GitLabListIssues)

	s.AddTool(mcp.NewTool("gitlab_create_milestone",
		mcp.WithDescription("Create a GitLab milestone. Map sprint goals to milestones."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Milestone title")),
		mcp.WithString("description", mcp.Description("Description")),
		mcp.WithString("due_date", mcp.Description("Due date (YYYY-MM-DD)")),
	), h.GitLabCreateMilestone)

	s.AddTool(mcp.NewTool("gitlab_milestones",
		mcp.WithDescription("List GitLab milestones."),
		mcp.WithString("state", mcp.Description("active, closed (default: active)")),
	), h.GitLabListMilestones)

	s.AddTool(mcp.NewTool("gitlab_merge_requests",
		mcp.WithDescription("List GitLab merge requests."),
		mcp.WithString("state", mcp.Description("opened, closed, merged, all (default: opened)")),
		mcp.WithNumber("limit", mcp.Description("Max results (default: 20)")),
	), h.GitLabListMRs)

	s.AddTool(mcp.NewTool("gitlab_read_file",
		mcp.WithDescription("Read a file from the GitLab repo."),
		mcp.WithString("path", mcp.Required(), mcp.Description("File path")),
		mcp.WithString("ref", mcp.Description("Branch/tag (default: main)")),
	), h.GitLabReadFile)

	s.AddTool(mcp.NewTool("gitlab_list_files",
		mcp.WithDescription("List files/directories in a GitLab repo path."),
		mcp.WithString("path", mcp.Description("Directory path (default: root)")),
		mcp.WithString("ref", mcp.Description("Branch/tag (default: main)")),
	), h.GitLabListFiles)
}
