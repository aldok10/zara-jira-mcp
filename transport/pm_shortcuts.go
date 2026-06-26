package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPMShortcuts(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm",
		mcp.WithDescription("Quick project status. THE main tool for PMs. Shows sprint progress, blockers, risks, pending actions. Auto-detects board."),
		mcp.WithNumber("board_id", mcp.Description("Board ID (auto-detected if not provided)")),
	), h.PMQuickStatus)

	s.AddTool(mcp.NewTool("pm_create",
		mcp.WithDescription("Create work item anywhere: Jira, GitHub, or GitLab. Simplest way to create tasks."),
		mcp.WithString("title", mcp.Required(), mcp.Description("What needs to be done")),
		mcp.WithString("platform", mcp.Description("jira (default), github, gitlab")),
		mcp.WithString("description", mcp.Description("Details")),
		mcp.WithString("project", mcp.Description("Project key (required for Jira)")),
		mcp.WithString("priority", mcp.Description("High, Medium, Low")),
		mcp.WithString("assignee", mcp.Description("Who should do this")),
		mcp.WithString("labels", mcp.Description("Comma-separated labels")),
		mcp.WithString("type", mcp.Description("Task, Bug, Story (default: Task)")),
	), h.PMCreate)

	s.AddTool(mcp.NewTool("pm_decide",
		mcp.WithDescription("Quick decision recording. Just say what was decided."),
		mcp.WithString("what", mcp.Required(), mcp.Description("What was decided")),
		mcp.WithString("why", mcp.Description("Why (optional)")),
		mcp.WithString("who", mcp.Description("Who decided (default: team)")),
	), h.PMDecide)

	s.AddTool(mcp.NewTool("pm_risk",
		mcp.WithDescription("Quick risk recording. Just say what could go wrong."),
		mcp.WithString("what", mcp.Required(), mcp.Description("What could go wrong")),
		mcp.WithString("severity", mcp.Description("critical, high, medium, low (default: medium)")),
		mcp.WithString("owner", mcp.Description("Who should handle this")),
	), h.PMRisk)

	s.AddTool(mcp.NewTool("pm_next",
		mcp.WithDescription("What should I do next? Suggests highest-priority PM action based on blockers, risks, and pending items."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.PMNext)
}
