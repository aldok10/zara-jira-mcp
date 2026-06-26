package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCalendarTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("calendar_events",
		mcp.WithDescription("List upcoming calendar events for the next N days. Shows meetings, ceremonies, deadlines."),
		mcp.WithNumber("days", mcp.Description("Number of days to look ahead (default: 7)")),
	), h.CalendarListEvents)

	s.AddTool(mcp.NewTool("calendar_create",
		mcp.WithDescription("Create a calendar event in Lark. Use for sprint ceremonies, deadlines, reminders."),
		mcp.WithString("summary", mcp.Required(), mcp.Description("Event title")),
		mcp.WithString("start", mcp.Required(), mcp.Description("Start time (format: 2006-01-02 15:04)")),
		mcp.WithString("end", mcp.Description("End time (default: 1 hour after start)")),
		mcp.WithString("description", mcp.Description("Event description")),
	), h.CalendarCreateEvent)

	s.AddTool(mcp.NewTool("calendar_schedule_meeting",
		mcp.WithDescription("Schedule a meeting with Lark video conference and attendees. Creates event + VC link + invites."),
		mcp.WithString("summary", mcp.Required(), mcp.Description("Meeting title")),
		mcp.WithString("start", mcp.Required(), mcp.Description("Start time (format: 2006-01-02 15:04)")),
		mcp.WithNumber("duration_minutes", mcp.Description("Duration in minutes (default: 60)")),
		mcp.WithString("description", mcp.Description("Meeting agenda/description")),
		mcp.WithString("attendees", mcp.Description("Comma-separated email addresses")),
	), h.CalendarScheduleMeeting)
}
