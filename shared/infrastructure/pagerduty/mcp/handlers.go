package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/pagerduty"
)

// Handlers exposes PagerDuty operations as MCP tool handlers.
type Handlers struct {
	client *pagerduty.Client
	errH   *mcputil.ErrorHandler
}

func NewHandlers(client *pagerduty.Client) *Handlers {
	return &Handlers{
		client: client,
		errH:   mcputil.NewErrorHandler(nil),
	}
}

// ListIncidents lists recent PagerDuty incidents.
func (h *Handlers) ListIncidents(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("PagerDuty not configured. Set PAGERDUTY_API_KEY."), nil
	}
	status := req.GetString("status", "")

	incidents, err := h.client.ListIncidents(ctx, status)
	if err != nil {
		return h.errH.Wrap("list incidents", err), nil
	}
	if len(incidents) == 0 {
		return mcputil.TextResult("No incidents found."), nil
	}
	var b strings.Builder
	for _, inc := range incidents {
		assignee := "unassigned"
		if inc.Assignee != "" {
			assignee = inc.Assignee
		}
		b.WriteString(fmt.Sprintf("[%s/%s] %s\n  Service: %s | Assigned: %s | Created: %s\n",
			inc.Status, inc.Urgency, inc.Title, inc.ServiceName, assignee, inc.CreatedAt[:10]))
	}
	return mcputil.TextResult(b.String()), nil
}

// GetOnCalls shows who is on call right now.
func (h *Handlers) GetOnCalls(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("PagerDuty not configured. Set PAGERDUTY_API_KEY."), nil
	}

	oncalls, err := h.client.GetOnCalls(ctx)
	if err != nil {
		return h.errH.Wrap("get on-calls", err), nil
	}
	if len(oncalls) == 0 {
		return mcputil.TextResult("No on-call schedules found."), nil
	}
	var b strings.Builder
	for _, o := range oncalls {
		b.WriteString(fmt.Sprintf("%s — %s (%s → %s)\n", o.UserName, o.Schedule, o.Start[:10], o.End[:10]))
	}
	return mcputil.TextResult(b.String()), nil
}
