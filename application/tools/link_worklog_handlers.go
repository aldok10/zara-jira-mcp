package tools

import (
	"context"
	"fmt"
	"strings"

	domain "github.com/aldok10/zara-jira-mcp/domain/jira"
	"github.com/mark3labs/mcp-go/mcp"
)

// LinkIssues creates a link between two issues.
func (h *Handlers) LinkIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	inwardKey, err := req.RequireString("inward_key")
	if err != nil {
		return errorResult("inward_key is required"), nil
	}
	outwardKey, err := req.RequireString("outward_key")
	if err != nil {
		return errorResult("outward_key is required"), nil
	}
	linkType, err := req.RequireString("link_type")
	if err != nil {
		return errorResult("link_type is required"), nil
	}

	if err := h.Jira.LinkIssues(ctx, inwardKey, outwardKey, linkType); err != nil {
		return errorResult("Failed to link issues: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Linked %s -> %s (type: %s)", inwardKey, outwardKey, linkType)), nil
}

// LinkTypes lists available issue link types.
func (h *Handlers) LinkTypes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	types, err := h.Jira.GetLinkTypes(ctx)
	if err != nil {
		return errorResult("Failed to get link types: " + err.Error()), nil
	}
	var sb strings.Builder
	sb.WriteString("Available link types:\n\n")
	for _, lt := range types {
		sb.WriteString(fmt.Sprintf("- **%s** (inward: %q, outward: %q)\n", lt.Name, lt.Inward, lt.Outward))
	}
	return textResult(sb.String()), nil
}

// WorklogAdd logs time on an issue.
func (h *Handlers) WorklogAdd(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key is required"), nil
	}
	timeSpent, err := req.RequireString("time_spent")
	if err != nil {
		return errorResult("time_spent is required"), nil
	}
	comment := req.GetString("comment", "")

	if err := h.Jira.AddWorklog(ctx, key, timeSpent, comment); err != nil {
		return errorResult("Failed to add worklog: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Logged %s on %s", timeSpent, key)), nil
}

// WorklogList lists worklogs for an issue.
func (h *Handlers) WorklogList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key is required"), nil
	}

	worklogs, err := h.Jira.GetWorklogs(ctx, key)
	if err != nil {
		return errorResult("Failed to get worklogs: " + err.Error()), nil
	}
	if len(worklogs) == 0 {
		return textResult(fmt.Sprintf("No worklogs on %s", key)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Worklogs for %s (%d):\n\n", key, len(worklogs)))
	for _, w := range worklogs {
		sb.WriteString(fmt.Sprintf("- %s | %s | %s", w.Author, w.TimeSpent, w.Started))
		if w.Comment != "" {
			sb.WriteString(fmt.Sprintf(" | %s", w.Comment))
		}
		sb.WriteString("\n")
	}
	return textResult(sb.String()), nil
}

// Watch adds a watcher to an issue.
func (h *Handlers) Watch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key is required"), nil
	}
	accountID, err := req.RequireString("account_id")
	if err != nil {
		return errorResult("account_id is required"), nil
	}

	if err := h.Jira.AddWatcher(ctx, key, accountID); err != nil {
		return errorResult("Failed to add watcher: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Added watcher to %s", key)), nil
}

// Watchers lists watchers on an issue.
func (h *Handlers) Watchers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key is required"), nil
	}

	watchers, err := h.Jira.GetWatchers(ctx, key)
	if err != nil {
		return errorResult("Failed to get watchers: " + err.Error()), nil
	}
	if len(watchers) == 0 {
		return textResult(fmt.Sprintf("No watchers on %s", key)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Watchers on %s (%d):\n\n", key, len(watchers)))
	for _, w := range watchers {
		sb.WriteString(fmt.Sprintf("- %s (%s)\n", w.DisplayName, w.AccountID))
	}
	return textResult(sb.String()), nil
}

// LabelsSet sets labels on an issue (replaces all).
func (h *Handlers) LabelsSet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key is required"), nil
	}
	labelsRaw, err := req.RequireString("labels")
	if err != nil {
		return errorResult("labels is required"), nil
	}

	labels := strings.Split(labelsRaw, ",")
	for i := range labels {
		labels[i] = strings.TrimSpace(labels[i])
	}

	input := &domain.UpdateIssueInput{Key: key, Labels: labels}
	if err := h.Jira.UpdateIssue(ctx, input); err != nil {
		return errorResult("Failed to set labels: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Labels set on %s: %s", key, strings.Join(labels, ", "))), nil
}
