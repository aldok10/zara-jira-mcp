// Package mcp provides MCP tool handlers for the jira module.
package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/modules/jira/application/port"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/validate"
)

// Handlers holds dependencies for jira MCP tool handlers.
type Handlers struct {
	Jira port.Inbound
}

// NewHandlers creates a new jira MCP handlers instance.
func NewHandlers(jiraService port.Inbound) *Handlers {
	return &Handlers{Jira: jiraService}
}

// Health returns server version and status.
func (h *Handlers) Health(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("zara-jira-mcp | status: ok | modular handlers"), nil
}

// SearchIssues searches Jira issues using JQL.
func (h *Handlers) SearchIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jql, err := req.RequireString("jql")
	if err != nil {
		return mcputil.ErrInvalid("jql parameter is required"), nil
	}
	if err := validate.JQL(jql); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	maxResults := req.GetInt("max_results", 20)

	results, err := h.Jira.SearchIssues(ctx, jql, int(maxResults))
	if err != nil {
		return mcputil.ErrJira("Jira search", err), nil
	}
	if results == nil {
		return mcputil.TextResult("No results found."), nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d issues\n", len(results.Issues)))
	for _, issue := range results.Issues {
		sb.WriteString(fmt.Sprintf("  %s - %s [%s]\n", issue.Key, issue.Summary, issue.Status))
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetIssue returns full details of a single Jira issue.
func (h *Handlers) GetIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	issue, err := h.Jira.GetIssue(ctx, key)
	if err != nil {
		return mcputil.ErrJira("get issue", err), nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s** - %s\n", issue.Key, issue.Summary))
	sb.WriteString(fmt.Sprintf("Type: %s | Status: %s | Priority: %s\n", issue.Type, issue.Status, issue.Priority))
	sb.WriteString(fmt.Sprintf("Assignee: %s\n", issue.Assignee))
	sb.WriteString(fmt.Sprintf("Description: %s\n", issue.Description))
	return mcputil.TextResult(sb.String()), nil
}

// GetBoards lists all accessible Jira boards.
func (h *Handlers) GetBoards(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boards, err := h.Jira.GetBoards(ctx)
	if err != nil {
		return mcputil.ErrJira("get boards", err), nil
	}

	var sb strings.Builder
	for _, b := range boards {
		sb.WriteString(fmt.Sprintf("%d: %s (%s)\n", b.ID, b.Name, b.Type))
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetSprintSummary returns the active sprint status for a board.
func (h *Handlers) GetSprintSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return mcputil.ErrInvalid("board_id parameter is required"), nil
	}
	if err := validate.BoardID(int(boardID)); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, int(boardID))
	if err != nil {
		return mcputil.ErrJira("get active sprints", err), nil
	}

	if len(sprints) == 0 {
		return mcputil.TextResult("No active sprints found for this board."), nil
	}

	sprint := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprint.ID)
	if err != nil {
		return mcputil.ErrJira("get sprint issues", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprint: %s (ID: %d)\n", sprint.Name, sprint.ID))
	sb.WriteString(fmt.Sprintf("Goals: %s\n", sprint.Goal))
	sb.WriteString(fmt.Sprintf("Start: %s | End: %s\n", sprint.StartDate, sprint.EndDate))
	sb.WriteString(fmt.Sprintf("Issues: %d\n", len(issues)))

	for _, i := range issues {
		sb.WriteString(fmt.Sprintf("  %s - %s [%s]\n", i.Key, i.Summary, i.Status))
	}
	return mcputil.TextResult(sb.String()), nil
}
