package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerFacilitationTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_retro_format",
		mcp.WithDescription("AI recommends retro format based on team context."),
		mcp.WithNumber("board_id", mcp.Description("Board ID for context")),
	), h.PMRetroFormat)

	s.AddTool(mcp.NewTool("pm_meeting_audit",
		mcp.WithDescription("Assess if a meeting could be async."),
		mcp.WithString("meeting_name", mcp.Required(), mcp.Description("Meeting name")),
		mcp.WithString("meeting_type", mcp.Required(), mcp.Description("standup, planning, retro, review, grooming, status_update, decision, brainstorming, 1on1, allhands, other")),
		mcp.WithNumber("duration_minutes", mcp.Description("Duration in minutes (default: 30)")),
		mcp.WithNumber("attendees", mcp.Description("Number of attendees (default: 5)")),
		mcp.WithString("frequency", mcp.Description("daily, weekly, biweekly, monthly (default: weekly)")),
		mcp.WithNumber("agenda_items", mcp.Description("Agenda items count (default: 3)")),
	), h.PMMeetingAudit)
}
