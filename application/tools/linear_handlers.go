package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handlers) LinearIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Linear == nil || !h.Linear.Available() {
		return errorResult("Linear not configured. Set LINEAR_API_KEY."), nil
	}
	team := req.GetString("team", "")
	state := req.GetString("state", "")

	issues, err := h.Linear.ListIssues(ctx, team, state)
	if err != nil {
		return errorResult("Linear API error: " + err.Error()), nil
	}
	if len(issues) == 0 {
		return textResult("No issues found."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Linear issues (%d):\n\n", len(issues)))
	for _, i := range issues {
		sb.WriteString(fmt.Sprintf("- [P%d] %s | %s | Assignee: %s | Team: %s\n",
			i.Priority, i.Title, i.State, i.Assignee, i.Team))
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) LinearCycles(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Linear == nil || !h.Linear.Available() {
		return errorResult("Linear not configured. Set LINEAR_API_KEY."), nil
	}

	cycles, err := h.Linear.ListCycles(ctx)
	if err != nil {
		return errorResult("Linear API error: " + err.Error()), nil
	}
	if len(cycles) == 0 {
		return textResult("No cycles found."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Linear cycles (%d):\n\n", len(cycles)))
	for _, c := range cycles {
		name := c.Name
		if name == "" {
			name = fmt.Sprintf("Cycle %d", c.Number)
		}
		sb.WriteString(fmt.Sprintf("- %s | %s to %s | Progress: %.0f%%\n",
			name, c.StartsAt[:10], c.EndsAt[:10], c.Progress*100))
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) LinearActivity(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Linear == nil || !h.Linear.Available() {
		return errorResult("Linear not configured. Set LINEAR_API_KEY."), nil
	}

	activities, err := h.Linear.RecentActivity(ctx)
	if err != nil {
		return errorResult("Linear API error: " + err.Error()), nil
	}
	if len(activities) == 0 {
		return textResult("No recent activity."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Recent activity (%d):\n\n", len(activities)))
	for _, a := range activities {
		sb.WriteString(fmt.Sprintf("- %s | %s | %s | by %s\n",
			a.CreatedAt[:10], a.Type, a.Issue, a.Actor))
	}
	return textResult(sb.String()), nil
}
