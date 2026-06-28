package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	jmcp "github.com/aldok10/zara-jira-mcp/modules/jira/interfaces/mcp"
)

func RegisterJiraTools(s *server.MCPServer, h *jmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_search",
			mcp.WithDescription("Search Jira issues using JQL."),
			mcp.WithString("jql", mcp.Required(), mcp.Description("JQL query string")),
			mcp.WithNumber("max_results", mcp.Description("Maximum results (default 20)")),
		),
		h.SearchIssues,
	)
	s.AddTool(
		mcp.NewTool("jira_get_issue",
			mcp.WithDescription("Get full details of a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
		),
		h.GetIssue,
	)
	s.AddTool(
		mcp.NewTool("jira_boards",
			mcp.WithDescription("List all accessible Jira boards."),
		),
		h.GetBoards,
	)
	s.AddTool(
		mcp.NewTool("jira_sprint_summary",
			mcp.WithDescription("Get active sprint status and issues for a board."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.GetSprintSummary,
	)
}
