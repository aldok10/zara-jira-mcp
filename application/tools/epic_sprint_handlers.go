package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// EpicIssues lists issues belonging to an epic via JQL.
func (h *Handlers) EpicIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	epicKey, err := req.RequireString("epic_key")
	if err != nil {
		return errorResult("epic_key is required"), nil
	}
	maxResults := req.GetInt("max_results", 50)

	jql := fmt.Sprintf(`"Epic Link" = %s ORDER BY rank ASC`, epicKey)
	result, err := h.Jira.SearchIssues(ctx, jql, maxResults, 0)
	if err != nil {
		return sanitizedError("epic search failed", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Epic %s has %d issues:\n\n", epicKey, result.Total))
	for _, issue := range result.Issues {
		sb.WriteString(fmt.Sprintf("**%s** [%s] %s\n  Status: %s | Assignee: %s\n\n",
			issue.Key, issue.Type, issue.Summary, issue.Status, issue.Assignee))
	}
	return textResult(sb.String()), nil
}

// EpicAdd adds issues to an epic by setting the parent field.
func (h *Handlers) EpicAdd(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	issueKeysRaw, err := req.RequireString("issue_keys")
	if err != nil {
		return errorResult("issue_keys is required"), nil
	}
	epicKey, err := req.RequireString("epic_key")
	if err != nil {
		return errorResult("epic_key is required"), nil
	}

	keys := strings.Split(issueKeysRaw, ",")
	var errors []string
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		if err := h.Jira.SetEpicLink(ctx, k, epicKey); err != nil {
			errors = append(errors, fmt.Sprintf("%s: operation failed", k))
		}
	}

	if len(errors) > 0 {
		return errorResult("Some issues failed:\n" + strings.Join(errors, "\n")), nil
	}
	return textResult(fmt.Sprintf("Added %d issue(s) to epic %s", len(keys), epicKey)), nil
}

// EpicRemove removes issues from their epic by clearing the parent field.
func (h *Handlers) EpicRemove(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	issueKeysRaw, err := req.RequireString("issue_keys")
	if err != nil {
		return errorResult("issue_keys is required"), nil
	}

	keys := strings.Split(issueKeysRaw, ",")
	var errors []string
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		if err := h.Jira.RemoveEpicLink(ctx, k); err != nil {
			errors = append(errors, fmt.Sprintf("%s: operation failed", k))
		}
	}

	if len(errors) > 0 {
		return errorResult("Some issues failed:\n" + strings.Join(errors, "\n")), nil
	}
	return textResult(fmt.Sprintf("Removed %d issue(s) from their epic", len(keys))), nil
}

// ListSprints lists sprints for a board.
func (h *Handlers) ListSprints(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id is required"), nil
	}
	state := req.GetString("state", "")

	sprints, err := h.Jira.GetSprints(ctx, boardID, state)
	if err != nil {
		return sanitizedError("failed to get sprints", err), nil
	}

	if len(sprints) == 0 {
		return textResult("No sprints found."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprints for board %d:\n\n", boardID))
	for _, s := range sprints {
		sb.WriteString(fmt.Sprintf("**#%d** %s [%s]", s.ID, s.Name, s.State))
		if s.StartDate != "" {
			sb.WriteString(fmt.Sprintf(" | %s - %s", s.StartDate, s.EndDate))
		}
		if s.Goal != "" {
			sb.WriteString(fmt.Sprintf("\n  Goal: %s", s.Goal))
		}
		sb.WriteString("\n\n")
	}
	return textResult(sb.String()), nil
}

// CreateSprintTool creates a new sprint on a board.
func (h *Handlers) CreateSprintTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id is required"), nil
	}
	name, err := req.RequireString("name")
	if err != nil {
		return errorResult("name is required"), nil
	}
	goal := req.GetString("goal", "")

	sprint, err := h.Jira.CreateSprint(ctx, boardID, name, goal)
	if err != nil {
		return sanitizedError("failed to create sprint", err), nil
	}
	return textResult(fmt.Sprintf("Created sprint #%d: %s [%s]", sprint.ID, sprint.Name, sprint.State)), nil
}

// StartSprintTool starts a sprint with given dates.
func (h *Handlers) StartSprintTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintID, err := req.RequireInt("sprint_id")
	if err != nil {
		return errorResult("sprint_id is required"), nil
	}
	startDate, err := req.RequireString("start_date")
	if err != nil {
		return errorResult("start_date is required"), nil
	}
	endDate, err := req.RequireString("end_date")
	if err != nil {
		return errorResult("end_date is required"), nil
	}

	if err := h.Jira.StartSprint(ctx, sprintID, startDate, endDate); err != nil {
		return sanitizedError("failed to start sprint", err), nil
	}
	return textResult(fmt.Sprintf("Sprint %d started (%s to %s)", sprintID, startDate, endDate)), nil
}

// CloseSprintTool closes/completes a sprint.
func (h *Handlers) CloseSprintTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintID, err := req.RequireInt("sprint_id")
	if err != nil {
		return errorResult("sprint_id is required"), nil
	}

	if err := h.Jira.CloseSprint(ctx, sprintID); err != nil {
		return sanitizedError("failed to close sprint", err), nil
	}
	return textResult(fmt.Sprintf("Sprint %d closed", sprintID)), nil
}

// MoveIssuesToSprintTool moves issues into a sprint.
func (h *Handlers) MoveIssuesToSprintTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintID, err := req.RequireInt("sprint_id")
	if err != nil {
		return errorResult("sprint_id is required"), nil
	}
	issueKeysRaw, err := req.RequireString("issue_keys")
	if err != nil {
		return errorResult("issue_keys is required"), nil
	}

	keys := strings.Split(issueKeysRaw, ",")
	trimmed := make([]string, 0, len(keys))
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k != "" {
			trimmed = append(trimmed, k)
		}
	}

	if err := h.Jira.MoveIssuesToSprint(ctx, sprintID, trimmed); err != nil {
		return sanitizedError("failed to move issues between sprints", err), nil
	}
	return textResult(fmt.Sprintf("Moved %d issue(s) to sprint %d", len(trimmed), sprintID)), nil
}
