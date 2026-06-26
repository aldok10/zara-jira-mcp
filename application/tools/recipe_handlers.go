package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// RecipeStartWork assigns, transitions to In Progress, and suggests a branch name.
func (h *Handlers) RecipeStartWork(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key required"), nil
	}

	issue, err := h.Jira.GetIssue(ctx, key)
	if err != nil {
		return sanitizedError("Failed to get issue", err), nil
	}

	transitions, _ := h.Jira.GetTransitions(ctx, key)
	var progressID string
	for _, t := range transitions {
		lower := strings.ToLower(t.Name)
		if strings.Contains(lower, "progress") || strings.Contains(lower, "start") {
			progressID = t.ID
			break
		}
	}

	var actions []string

	assigneeID := req.GetString("assignee_id", "")
	if assigneeID != "" {
		if err := h.Jira.AssignIssue(ctx, key, assigneeID); err == nil {
			actions = append(actions, "Assigned to you")
		}
	}

	if progressID != "" {
		if err := h.Jira.TransitionIssue(ctx, key, progressID); err == nil {
			actions = append(actions, "Transitioned to In Progress")
		}
	} else {
		actions = append(actions, "No 'In Progress' transition found (may already be in progress)")
	}

	slug := strings.ToLower(issue.Summary)
	slug = strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			return r
		}
		if r == ' ' || r == '-' || r == '_' {
			return '-'
		}
		return -1
	}, slug)
	if len(slug) > 50 {
		slug = slug[:50]
	}
	branchName := fmt.Sprintf("feature/%s-%s", strings.ToLower(key), strings.Trim(slug, "-"))

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Started: %s - %s\n\n", key, issue.Summary))
	sb.WriteString("Actions:\n")
	for _, a := range actions {
		sb.WriteString(fmt.Sprintf("  - %s\n", a))
	}
	sb.WriteString(fmt.Sprintf("\nSuggested branch: %s\n", branchName))
	sb.WriteString(fmt.Sprintf("git checkout -b %s\n", branchName))

	return textResult(sb.String()), nil
}

// RecipeDone transitions to Done, optionally logs time and adds a comment.
func (h *Handlers) RecipeDone(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key required"), nil
	}

	var actions []string

	transitions, _ := h.Jira.GetTransitions(ctx, key)
	var doneID string
	for _, t := range transitions {
		lower := strings.ToLower(t.Name)
		if strings.Contains(lower, "done") || strings.Contains(lower, "close") || strings.Contains(lower, "resolve") {
			doneID = t.ID
			break
		}
	}

	if doneID != "" {
		if err := h.Jira.TransitionIssue(ctx, key, doneID); err == nil {
			actions = append(actions, "Transitioned to Done")
		} else {
			actions = append(actions, "Transition failed: "+err.Error())
		}
	} else {
		actions = append(actions, "No 'Done' transition available")
	}

	timeSpent := req.GetString("time_spent", "")
	if timeSpent != "" {
		if err := h.Jira.AddWorklog(ctx, key, timeSpent, ""); err == nil {
			actions = append(actions, fmt.Sprintf("Logged %s", timeSpent))
		}
	}

	comment := req.GetString("comment", "")
	if comment != "" {
		if err := h.Jira.AddComment(ctx, key, comment); err == nil {
			actions = append(actions, "Comment added")
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Completed: %s\n\nActions:\n", key))
	for _, a := range actions {
		sb.WriteString(fmt.Sprintf("  - %s\n", a))
	}
	return textResult(sb.String()), nil
}

// RecipeBlock flags an issue as blocked: saves to memory and comments on the issue.
func (h *Handlers) RecipeBlock(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key required"), nil
	}
	reason, err := req.RequireString("reason")
	if err != nil {
		return errorResult("reason required"), nil
	}

	var actions []string

	b := &memory.Blocker{
		IssueKey:     key,
		Description:  reason,
		BlockedSince: time.Now(),
		Owner:        req.GetString("owner", ""),
	}
	if err := h.Memory.SaveBlocker(ctx, b); err == nil {
		actions = append(actions, "Blocker recorded in memory")
	}

	comment := fmt.Sprintf("BLOCKED: %s", reason)
	if err := h.Jira.AddComment(ctx, key, comment); err == nil {
		actions = append(actions, "Blocker comment added to issue")
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Blocked: %s\nReason: %s\n\nActions:\n", key, reason))
	for _, a := range actions {
		sb.WriteString(fmt.Sprintf("  - %s\n", a))
	}
	return textResult(sb.String()), nil
}
