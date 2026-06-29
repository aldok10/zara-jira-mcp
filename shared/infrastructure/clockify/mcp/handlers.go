// Package mcp provides MCP tool handlers for Clockify time tracking.
package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/clockify"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

// Handlers exposes Clockify time tracking operations as MCP tool handlers.
type Handlers struct {
	client *clockify.Client
	errH   *mcputil.ErrorHandler
}

func NewHandlers(client *clockify.Client) *Handlers {
	return &Handlers{
		client: client,
		errH:   mcputil.NewErrorHandler(nil),
	}
}

// TimeEntries lists recent time entries.
func (h *Handlers) TimeEntries(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Clockify not configured. Set CLOCKIFY_API_KEY and CLOCKIFY_WORKSPACE_ID."), nil
	}

	days := req.GetInt("days", 1)
	now := time.Now()
	from := now.AddDate(0, 0, -days)
	to := now

	entries, err := h.client.GetTimeEntries(ctx, from, to)
	if err != nil {
		return h.errH.Wrap("get time entries", err), nil
	}

	if len(entries) == 0 {
		return mcputil.TextResult("No time entries found."), nil
	}

	var b strings.Builder
	for _, e := range entries {
		b.WriteString(fmt.Sprintf("%s — %s (%s)\n  Project: %s | Duration: %s\n",
			e.UserName, e.Description, e.Start.Format("Mon 15:04"),
			e.ProjectName, fmtDuration(e.Duration)))
	}
	return mcputil.TextResult(b.String()), nil
}

// SummaryReport shows time summary grouped by user and project.
func (h *Handlers) SummaryReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Clockify not configured. Set CLOCKIFY_API_KEY and CLOCKIFY_WORKSPACE_ID."), nil
	}

	days := req.GetInt("days", 7)
	now := time.Now()
	from := now.AddDate(0, 0, -days)
	to := now

	report, err := h.client.GetSummaryReport(ctx, from, to)
	if err != nil {
		return h.errH.Wrap("get summary report", err), nil
	}

	if len(report) == 0 {
		return mcputil.TextResult("No time data found for the period."), nil
	}

	var b strings.Builder
	for user, projects := range report {
		b.WriteString(fmt.Sprintf("%s:\n", user))
		total := time.Duration(0)
		for proj, dur := range projects {
			b.WriteString(fmt.Sprintf("  %s: %s\n", proj, fmtDuration(dur)))
			total += dur
		}
		b.WriteString(fmt.Sprintf("  Total: %s\n\n", fmtDuration(total)))
	}
	return mcputil.TextResult(b.String()), nil
}

func fmtDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
