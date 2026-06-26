package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerStoryPointsTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_story_points",
		mcp.WithDescription("Calculate total story points from sprint, epic, or JQL query. Groups by status/assignee/type. Auto-reads custom field."),
		mcp.WithNumber("board_id", mcp.Description("Board ID (uses active sprint)")),
		mcp.WithString("epic_key", mcp.Description("Epic key to sum points for")),
		mcp.WithString("jql", mcp.Description("Custom JQL query")),
		mcp.WithString("group_by", mcp.Description("Group results: status, assignee, type (default: all)")),
	), h.StoryPointsSummary)

	s.AddTool(mcp.NewTool("pm_sprint_points_burndown",
		mcp.WithDescription("Story points burndown for active sprint. Shows done/in-progress/todo/blocked points with percentages."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.SprintPointsBurndown)
}
