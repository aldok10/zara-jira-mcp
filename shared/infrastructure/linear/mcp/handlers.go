package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/linear"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

type Handlers struct {
	client *linear.Client
	errH   *mcputil.ErrorHandler
}

func NewHandlers(client *linear.Client) *Handlers {
	return &Handlers{
		client: client,
		errH:   mcputil.NewErrorHandler(nil),
	}
}

func (h *Handlers) ListIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Linear not configured. Set LINEAR_API_KEY."), nil
	}
	team := req.GetString("team", "")
	state := req.GetString("state", "")

	issues, err := h.client.ListIssues(ctx, team, state)
	if err != nil {
		return h.errH.Wrap("list issues", err), nil
	}
	if len(issues) == 0 {
		return mcputil.TextResult("No issues found."), nil
	}
	var b strings.Builder
	for _, iss := range issues {
		assignee := "unassigned"
		if iss.Assignee != "" {
			assignee = iss.Assignee
		}
		priority := fmt.Sprintf("P%d", iss.Priority)
		if iss.Priority == 0 {
			priority = "No priority"
		}
		b.WriteString(fmt.Sprintf("%s [%s] %s — %s (%s)\n", iss.ID[:8], iss.State, iss.Title, assignee, priority))
	}
	return mcputil.TextResult(b.String()), nil
}

func (h *Handlers) ListCycles(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Linear not configured."), nil
	}

	cycles, err := h.client.ListCycles(ctx)
	if err != nil {
		return h.errH.Wrap("list cycles", err), nil
	}
	if len(cycles) == 0 {
		return mcputil.TextResult("No cycles found."), nil
	}
	var b strings.Builder
	for _, c := range cycles {
		b.WriteString(fmt.Sprintf("Cycle %d: %s (%s → %s) — %.0f%%\n", c.Number, c.Name, c.StartsAt[:10], c.EndsAt[:10], c.Progress*100))
	}
	return mcputil.TextResult(b.String()), nil
}

func (h *Handlers) RecentActivity(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Linear not configured."), nil
	}

	activities, err := h.client.RecentActivity(ctx)
	if err != nil {
		return h.errH.Wrap("get recent activity", err), nil
	}
	if len(activities) == 0 {
		return mcputil.TextResult("No recent activity."), nil
	}
	var b strings.Builder
	for _, a := range activities {
		actor := a.Actor
		if actor == "" {
			actor = "system"
		}
		b.WriteString(fmt.Sprintf("[%s] %s — %s by %s\n", a.CreatedAt[:10], a.Type, a.Issue, actor))
	}
	return mcputil.TextResult(b.String()), nil
}
