package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerGitHubTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("jira_link_pr",
		mcp.WithDescription("Link a PR/commit URL to a Jira issue by adding a comment with the link."),
		mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		mcp.WithString("pr_url", mcp.Required(), mcp.Description("Pull request or commit URL")),
		mcp.WithString("title", mcp.Description("PR title (default: 'Pull Request')")),
	), h.LinkPR)

	s.AddTool(mcp.NewTool("jira_from_branch",
		mcp.WithDescription("Extract Jira issue key from a git branch name (e.g. feature/PROJ-123-desc) and fetch issue details."),
		mcp.WithString("branch", mcp.Required(), mcp.Description("Git branch name")),
	), h.IssueFromBranch)

	s.AddTool(mcp.NewTool("jira_smart_commit",
		mcp.WithDescription("Apply Jira smart commit actions from a commit message. Supports #done, #time Xh, #comment text."),
		mcp.WithString("message", mcp.Required(), mcp.Description("Commit message (e.g. 'PROJ-123 fix login #done #time 2h')")),
	), h.SmartCommit)
}
