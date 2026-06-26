package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aldok10/zara-jira-mcp/internal/calendar"
	"github.com/mark3labs/mcp-go/mcp"
)

// CalendarListEvents lists upcoming calendar events.
func (h *Handlers) CalendarListEvents(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Calendar == nil || !h.Calendar.Available() {
		return errorResult("Calendar not configured. Set LARK_APP_ID + LARK_APP_SECRET."), nil
	}

	days := req.GetInt("days", 7)
	from := time.Now()
	to := from.Add(time.Duration(days) * 24 * time.Hour)

	events, err := h.Calendar.ListEvents(ctx, from, to)
	if err != nil {
		return errorResult("Failed to list events: " + err.Error()), nil
	}

	if len(events) == 0 {
		return textResult(fmt.Sprintf("No events in the next %d days.", days)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Upcoming events (next %d days): %d\n\n", days, len(events)))
	for _, ev := range events {
		sb.WriteString(fmt.Sprintf("- %s — %s to %s\n",
			ev.Summary,
			ev.StartTime.Format("Mon Jan 2 15:04"),
			ev.EndTime.Format("15:04")))
		if ev.Location != "" {
			sb.WriteString(fmt.Sprintf("  Location: %s\n", ev.Location))
		}
	}
	return textResult(sb.String()), nil
}

// CalendarCreateEvent creates a calendar event.
func (h *Handlers) CalendarCreateEvent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Calendar == nil || !h.Calendar.Available() {
		return errorResult("Calendar not configured."), nil
	}

	summary, err := req.RequireString("summary")
	if err != nil {
		return errorResult("summary parameter required"), nil
	}
	startStr, err := req.RequireString("start")
	if err != nil {
		return errorResult("start parameter required (format: 2006-01-02 15:04)"), nil
	}
	endStr := req.GetString("end", "")
	description := req.GetString("description", "")

	start, err := parseFlexTime(startStr)
	if err != nil {
		return errorResult("Invalid start time: " + err.Error()), nil
	}

	var end time.Time
	if endStr != "" {
		end, err = parseFlexTime(endStr)
		if err != nil {
			return errorResult("Invalid end time: " + err.Error()), nil
		}
	} else {
		end = start.Add(1 * time.Hour) // default 1 hour
	}

	ev := &calendar.Event{
		Summary:     summary,
		Description: description,
		StartTime:   start,
		EndTime:     end,
	}

	created, err := h.Calendar.CreateEvent(ctx, ev)
	if err != nil {
		return errorResult("Failed to create event: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Event created: %s (ID: %s)", summary, created.ID)), nil
}

// CalendarScheduleMeeting creates a meeting with video conference and attendees.
func (h *Handlers) CalendarScheduleMeeting(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Calendar == nil || !h.Calendar.Available() {
		return errorResult("Calendar not configured."), nil
	}

	summary, err := req.RequireString("summary")
	if err != nil {
		return errorResult("summary parameter required"), nil
	}
	startStr, err := req.RequireString("start")
	if err != nil {
		return errorResult("start parameter required (format: 2006-01-02 15:04)"), nil
	}

	start, err := parseFlexTime(startStr)
	if err != nil {
		return errorResult("Invalid start time: " + err.Error()), nil
	}

	duration := req.GetInt("duration_minutes", 60)
	end := start.Add(time.Duration(duration) * time.Minute)
	description := req.GetString("description", "")
	attendeesStr := req.GetString("attendees", "")

	var attendees []string
	if attendeesStr != "" {
		attendees = strings.Split(attendeesStr, ",")
		for i := range attendees {
			attendees[i] = strings.TrimSpace(attendees[i])
		}
	}

	meeting, err := h.Calendar.ScheduleMeeting(ctx, summary, description, start, end, attendees)
	if err != nil {
		return errorResult("Failed to schedule meeting: " + err.Error()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Meeting scheduled: %s\n", summary))
	sb.WriteString(fmt.Sprintf("Time: %s — %s (%d min)\n", start.Format("Mon Jan 2 15:04"), end.Format("15:04"), duration))
	if len(attendees) > 0 {
		sb.WriteString(fmt.Sprintf("Attendees: %s\n", strings.Join(attendees, ", ")))
	}
	sb.WriteString(fmt.Sprintf("ID: %s\n", meeting.ID))
	sb.WriteString("Video conference: Lark VC (auto-created)")
	return textResult(sb.String()), nil
}

func parseFlexTime(s string) (time.Time, error) {
	formats := []string{
		"2006-01-02 15:04",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		time.RFC3339,
	}
	loc, _ := time.LoadLocation("Asia/Jakarta")
	for _, f := range formats {
		if t, err := time.ParseInLocation(f, s, loc); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("use format: 2006-01-02 15:04")
}
