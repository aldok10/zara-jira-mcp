package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerEpicSprintTools(s *server.MCPServer, h *tools.Handlers) {
	// Epic tools
	s.AddTool(
		mcp.NewTool("jira_epic_issues",
			mcp.WithDescription("List all issues in an epic. Uses JQL to find linked issues."),
			mcp.WithString("epic_key", mcp.Required(), mcp.Description("Epic issue key (e.g. PROJ-100)")),
			mcp.WithNumber("max_results", mcp.Description("Maximum results (default 50)")),
		),
		h.EpicIssues,
	)

	s.AddTool(
		mcp.NewTool("jira_epic_add",
			mcp.WithDescription("Add issues to an epic by setting the parent/epic link."),
			mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated issue keys to add")),
			mcp.WithString("epic_key", mcp.Required(), mcp.Description("Epic key to add issues to")),
		),
		h.EpicAdd,
	)

	s.AddTool(
		mcp.NewTool("jira_epic_remove",
			mcp.WithDescription("Remove issues from their epic by clearing the parent link."),
			mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated issue keys to remove from epic")),
		),
		h.EpicRemove,
	)

	// Sprint tools
	s.AddTool(
		mcp.NewTool("jira_sprints",
			mcp.WithDescription("List sprints for a board. Filter by state: active, future, closed."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID (get from jira_boards)")),
			mcp.WithString("state", mcp.Description("Filter by state: active, future, closed (default: all)")),
		),
		h.ListSprints,
	)

	s.AddTool(
		mcp.NewTool("jira_sprint_create",
			mcp.WithDescription("Create a new sprint on a board."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Sprint name")),
			mcp.WithString("goal", mcp.Description("Sprint goal")),
		),
		h.CreateSprintTool,
	)

	s.AddTool(
		mcp.NewTool("jira_sprint_start",
			mcp.WithDescription("Start a sprint with start and end dates."),
			mcp.WithNumber("sprint_id", mcp.Required(), mcp.Description("Sprint ID")),
			mcp.WithString("start_date", mcp.Required(), mcp.Description("Start date (YYYY-MM-DD)")),
			mcp.WithString("end_date", mcp.Required(), mcp.Description("End date (YYYY-MM-DD)")),
		),
		h.StartSprintTool,
	)

	s.AddTool(
		mcp.NewTool("jira_sprint_close",
			mcp.WithDescription("Close/complete a sprint."),
			mcp.WithNumber("sprint_id", mcp.Required(), mcp.Description("Sprint ID")),
		),
		h.CloseSprintTool,
	)

	s.AddTool(
		mcp.NewTool("jira_sprint_move_issues",
			mcp.WithDescription("Move issues into a sprint."),
			mcp.WithNumber("sprint_id", mcp.Required(), mcp.Description("Sprint ID")),
			mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated issue keys")),
		),
		h.MoveIssuesToSprintTool,
	)
}
