package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	calmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/calendar/mcp"
)

func RegisterCalendarTools(s *server.MCPServer, h *calmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_calendar_create",
			mcp.WithDescription("Create a calendar event in Lark. Use for sprint ceremonies, deadlines, reminders."),
			mcp.WithString("summary", mcp.Required(), mcp.Description("Event title")),
			mcp.WithString("start", mcp.Required(), mcp.Description("Start time (format: 2006-01-02 15:04)")),
			mcp.WithString("end", mcp.Description("End time (default: 1 hour after start)")),
			mcp.WithString("description", mcp.Description("Event description")),
		),
		h.CreateEvent,
	)
	s.AddTool(
		mcp.NewTool("pm_calendar_events",
			mcp.WithDescription("List upcoming calendar events for the next N days. Shows meetings, ceremonies, deadlines."),
			mcp.WithInteger("days", mcp.Description("Number of days to look ahead (default: 7)")),
		),
		h.ListEvents,
	)
	s.AddTool(
		mcp.NewTool("pm_calendar_schedule_meeting",
			mcp.WithDescription("Schedule a meeting with Lark video conference and attendees. Creates event + VC link + invites."),
			mcp.WithString("summary", mcp.Required(), mcp.Description("Meeting title")),
			mcp.WithString("start", mcp.Required(), mcp.Description("Start time (format: 2006-01-02 15:04)")),
			mcp.WithString("description", mcp.Description("Meeting agenda/description")),
			mcp.WithInteger("duration_minutes", mcp.Description("Duration in minutes (default: 60)")),
			mcp.WithString("attendees", mcp.Description("Comma-separated email addresses")),
		),
		h.ScheduleMeeting,
	)
}
