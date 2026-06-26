package tools

import (
	"context"
	"fmt"
	"strings"

	domain "github.com/aldok10/zara-jira-mcp/domain/jira"
	"github.com/mark3labs/mcp-go/mcp"
)

// AssignIssue assigns an issue to a user by account ID.
func (h *Handlers) AssignIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter is required"), nil
	}
	assigneeID, err := req.RequireString("assignee_id")
	if err != nil {
		return errorResult("assignee_id parameter is required"), nil
	}

	if err := h.Jira.AssignIssue(ctx, key, assigneeID); err != nil {
		return sanitizedError("failed to assign issue", err), nil
	}
	return textResult(fmt.Sprintf("Issue %s assigned successfully", key)), nil
}

// UnassignIssue removes the assignee from an issue.
func (h *Handlers) UnassignIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter is required"), nil
	}

	if err := h.Jira.AssignIssue(ctx, key, ""); err != nil {
		return sanitizedError("failed to unassign issue", err), nil
	}
	return textResult(fmt.Sprintf("Issue %s unassigned successfully", key)), nil
}

// DeleteIssue deletes a Jira issue.
func (h *Handlers) DeleteIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter is required"), nil
	}

	if err := h.Jira.DeleteIssue(ctx, key); err != nil {
		return sanitizedError("failed to delete issue", err), nil
	}
	return textResult(fmt.Sprintf("Issue %s deleted successfully", key)), nil
}

// CreateSubtask creates a subtask under a parent issue.
func (h *Handlers) CreateSubtask(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	parentKey, err := req.RequireString("parent_key")
	if err != nil {
		return errorResult("parent_key parameter is required"), nil
	}
	summary, err := req.RequireString("summary")
	if err != nil {
		return errorResult("summary parameter is required"), nil
	}

	// Extract project key from parent key (e.g. "PROJ-123" -> "PROJ")
	parts := strings.SplitN(parentKey, "-", 2)
	if len(parts) != 2 {
		return errorResult("invalid parent_key format, expected PROJ-123"), nil
	}

	input := &domain.CreateIssueInput{
		Project:     parts[0],
		Summary:     summary,
		Description: req.GetString("description", ""),
		Priority:    req.GetString("priority", ""),
		Assignee:    req.GetString("assignee_id", ""),
	}

	created, err := h.Jira.CreateSubtask(ctx, parentKey, input)
	if err != nil {
		return sanitizedError("failed to create subtask", err), nil
	}
	return textResult(fmt.Sprintf("Created subtask: %s - %s (parent: %s)", created.Key, created.Summary, parentKey)), nil
}

// FindUser searches for Jira users by name or email.
func (h *Handlers) FindUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return errorResult("query parameter is required"), nil
	}

	users, err := h.Jira.FindUser(ctx, query)
	if err != nil {
		return sanitizedError("user search failed", err), nil
	}

	if len(users) == 0 {
		return textResult("No users found matching: " + query), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d user(s):\n\n", len(users)))
	for _, u := range users {
		sb.WriteString(fmt.Sprintf("- %s | ID: %s | Email: %s\n", u.DisplayName, u.AccountID, u.Email))
	}
	return textResult(sb.String()), nil
}
