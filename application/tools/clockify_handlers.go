package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handlers) ClockifyTimeReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Clockify == nil || !h.Clockify.Available() {
		return errorResult("Clockify not configured. Set CLOCKIFY_API_KEY and CLOCKIFY_WORKSPACE_ID."), nil
	}

	daysBack := req.GetInt("days", 7)
	end := time.Now()
	start := end.AddDate(0, 0, -daysBack)

	report, err := h.Clockify.GetSummaryReport(ctx, start, end)
	if err != nil {
		return errorResult("Clockify API error: " + err.Error()), nil
	}
	if len(report) == 0 {
		return textResult("No time tracked in this period."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Time report (last %d days):\n\n", daysBack))
	for user, projects := range report {
		var total time.Duration
		for _, d := range projects {
			total += d
		}
		sb.WriteString(fmt.Sprintf("**%s** (%.1fh total)\n", user, total.Hours()))
		for proj, d := range projects {
			sb.WriteString(fmt.Sprintf("  - %s: %.1fh\n", proj, d.Hours()))
		}
		sb.WriteString("\n")
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) ClockifyTimeEntries(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Clockify == nil || !h.Clockify.Available() {
		return errorResult("Clockify not configured. Set CLOCKIFY_API_KEY and CLOCKIFY_WORKSPACE_ID."), nil
	}

	daysBack := req.GetInt("days", 1)
	end := time.Now()
	start := end.AddDate(0, 0, -daysBack)

	entries, err := h.Clockify.GetTimeEntries(ctx, start, end)
	if err != nil {
		return errorResult("Clockify API error: " + err.Error()), nil
	}
	if len(entries) == 0 {
		return textResult("No time entries in this period."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Time entries (last %d day(s)): %d\n\n", daysBack, len(entries)))
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("- %s | %s | %s | %.1fh\n",
			e.UserName, e.ProjectName, e.Description, e.Duration.Hours()))
	}
	return textResult(sb.String()), nil
}
