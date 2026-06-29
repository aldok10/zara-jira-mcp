// Package mcp provides MCP tool handlers for Lark Calendar.
package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/calendar"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

// Handlers exposes Lark Calendar operations as MCP tool handlers.
type Handlers struct {
	client *calendar.Client
	errH   *mcputil.ErrorHandler
}

func NewHandlers(client *calendar.Client) *Handlers {
	return &Handlers{
		client: client,
		errH:   mcputil.NewErrorHandler(nil),
	}
}

// CreateEvent creates a calendar event.
func (h *Handlers) CreateEvent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Calendar not configured. Set LARK_APP_ID and LARK_APP_SECRET."), nil
	}

	summary, err := req.RequireString("summary")
	if err != nil {
		return mcputil.ErrorResult("Missing required parameter: summary"), nil
	}
	startStr, err := req.RequireString("start")
	if err != nil {
		return mcputil.ErrorResult("Missing required parameter: start"), nil
	}
	endStr := req.GetString("end", "")
	desc := req.GetString("description", "")

	start, err := time.Parse("2006-01-02 15:04", startStr)
	if err != nil {
		return mcputil.ErrorResult(fmt.Sprintf("Invalid start time %q. Use format: 2006-01-02 15:04", startStr)), nil
	}

	var end time.Time
	if endStr == "" {
		end = start.Add(1 * time.Hour)
	} else {
		end, err = time.Parse("2006-01-02 15:04", endStr)
		if err != nil {
			return mcputil.ErrorResult(fmt.Sprintf("Invalid end time %q. Use format: 2006-01-02 15:04", endStr)), nil
		}
	}

	ev := &calendar.Event{
		Summary:     summary,
		Description: desc,
		StartTime:   start,
		EndTime:     end,
	}

	created, err := h.client.CreateEvent(ctx, ev)
	if err != nil {
		return h.errH.Wrap("create calendar event", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Event created: %s (ID: %s)", created.Summary, created.ID)), nil
}

// ListEvents lists upcoming calendar events.
func (h *Handlers) ListEvents(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Calendar not configured. Set LARK_APP_ID and LARK_APP_SECRET."), nil
	}

	days := req.GetInt("days", 7)
	now := time.Now()
	from := now
	to := now.Add(time.Duration(days) * 24 * time.Hour)

	events, err := h.client.ListEvents(ctx, from, to)
	if err != nil {
		return h.errH.Wrap("list calendar events", err), nil
	}

	if len(events) == 0 {
		return mcputil.TextResult("No upcoming events found."), nil
	}

	var b strings.Builder
	for _, ev := range events {
		b.WriteString(fmt.Sprintf("%s (%s → %s)\n  %s\n",
			ev.Summary,
			ev.StartTime.Format("Mon 15:04"),
			ev.EndTime.Format("15:04"),
			ev.Description,
		))
	}
	return mcputil.TextResult(b.String()), nil
}

// ScheduleMeeting creates a meeting event with video conference.
func (h *Handlers) ScheduleMeeting(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Calendar not configured. Set LARK_APP_ID and LARK_APP_SECRET."), nil
	}

	summary, err := req.RequireString("summary")
	if err != nil {
		return mcputil.ErrorResult("Missing required parameter: summary"), nil
	}
	startStr, err := req.RequireString("start")
	if err != nil {
		return mcputil.ErrorResult("Missing required parameter: start"), nil
	}
	desc := req.GetString("description", "")
	durationMin := req.GetInt("duration_minutes", 60)
	attendeesStr := req.GetString("attendees", "")

	start, err := time.Parse("2006-01-02 15:04", startStr)
	if err != nil {
		return mcputil.ErrorResult(fmt.Sprintf("Invalid start time %q. Use format: 2006-01-02 15:04", startStr)), nil
	}
	end := start.Add(time.Duration(durationMin) * time.Minute)

	var attendees []string
	if attendeesStr != "" {
		attendees = strings.Split(attendeesStr, ",")
		for i := range attendees {
			attendees[i] = strings.TrimSpace(attendees[i])
		}
	}

	created, err := h.client.ScheduleMeeting(ctx, summary, desc, start, end, attendees)
	if err != nil {
		return h.errH.Wrap("schedule meeting", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Meeting scheduled: %s (ID: %s)", created.Summary, created.ID)), nil
}
