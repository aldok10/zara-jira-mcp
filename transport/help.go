package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerHelpTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_help",
		mcp.WithDescription("Discover available tools by topic. Call with no arguments for full menu, or specify a topic like 'sprint', 'risks', 'notifications'."),
		mcp.WithString("topic", mcp.Description("Topic to search: sprint, risks, team, planning, reporting, notifications, integrations, daily, retro, decisions")),
	), h.PMHelp)

	s.AddTool(mcp.NewTool("pm_quickstart",
		mcp.WithDescription("First-time getting started guide. Shows which integrations are configured and suggests first actions."),
	), h.PMQuickstart)

	s.AddTool(mcp.NewTool("pm_workflow",
		mcp.WithDescription("Pre-built workflow recipes. Get step-by-step tool sequences for common PM ceremonies and activities."),
		mcp.WithString("workflow", mcp.Required(), mcp.Description("Workflow: standup, sprint_start, sprint_end, planning, retro, incident, weekly_review")),
	), h.PMWorkflow)
}
