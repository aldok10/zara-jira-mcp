package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/gitlab"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

// Handlers exposes GitLab operations as MCP tool handlers.
type Handlers struct {
	client *gitlab.Client
	errH   *mcputil.ErrorHandler
}

func NewHandlers(client *gitlab.Client) *Handlers {
	return &Handlers{
		client: client,
		errH:   mcputil.NewErrorHandler(nil),
	}
}

// ListMRs lists merge requests.
func (h *Handlers) ListMRs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitLab not configured. Set GITLAB_TOKEN and GITLAB_PROJECT_ID."), nil
	}
	state := req.GetString("state", "opened")
	limit := req.GetInt("limit", 20)

	mrs, err := h.client.ListMRs(ctx, state, limit)
	if err != nil {
		return h.errH.Wrap("list MRs", err), nil
	}
	if len(mrs) == 0 {
		return mcputil.TextResult("No merge requests found."), nil
	}
	var b strings.Builder
	for _, mr := range mrs {
		draft := ""
		if mr.Draft {
			draft = " [DRAFT]"
		}
		age := time.Since(mr.CreatedAt).Truncate(time.Hour)
		b.WriteString(fmt.Sprintf("!%d%s %s (%s)\n  Author: %s | Age: %s\n",
			mr.IID, draft, mr.Title, mr.State, mr.Author, age))
	}
	return mcputil.TextResult(b.String()), nil
}

// ListIssues lists GitLab issues.
func (h *Handlers) ListIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitLab not configured."), nil
	}
	state := req.GetString("state", "opened")
	labels := req.GetString("labels", "")
	limit := req.GetInt("limit", 20)

	issues, err := h.client.ListIssues(ctx, state, labels, limit)
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
		b.WriteString(fmt.Sprintf("#%d %s [%s] — %s\n", iss.IID, iss.Title, iss.State, assignee))
	}
	return mcputil.TextResult(b.String()), nil
}

// CreateIssue creates a GitLab issue.
func (h *Handlers) CreateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitLab not configured."), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return mcputil.ErrInvalid("title parameter is required"), nil
	}
	description := req.GetString("description", "")
	labelsStr := req.GetString("labels", "")
	assigneeID := req.GetInt("assignee_id", 0)
	milestoneID := req.GetInt("milestone_id", 0)

	var labels []string
	if labelsStr != "" {
		labels = strings.Split(labelsStr, ",")
		for i := range labels {
			labels[i] = strings.TrimSpace(labels[i])
		}
	}

	issue, err := h.client.CreateIssue(ctx, title, description, labels, assigneeID, milestoneID)
	if err != nil {
		return h.errH.Wrap("create issue", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Created issue #%d: %s [%s]", issue.IID, issue.Title, issue.State)), nil
}

// ListMilestones lists GitLab milestones.
func (h *Handlers) ListMilestones(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitLab not configured."), nil
	}
	state := req.GetString("state", "active")

	milestones, err := h.client.ListMilestones(ctx, state)
	if err != nil {
		return h.errH.Wrap("list milestones", err), nil
	}
	if len(milestones) == 0 {
		return mcputil.TextResult("No milestones found."), nil
	}
	var b strings.Builder
	for _, m := range milestones {
		b.WriteString(fmt.Sprintf("#%d %s [%s]", m.IID, m.Title, m.State))
		if m.DueDate != "" {
			b.WriteString(fmt.Sprintf(" (due: %s)", m.DueDate[:10]))
		}
		b.WriteString("\n")
	}
	return mcputil.TextResult(b.String()), nil
}

// CreateMilestone creates a GitLab milestone.
func (h *Handlers) CreateMilestone(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitLab not configured."), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return mcputil.ErrInvalid("title parameter is required"), nil
	}
	description := req.GetString("description", "")
	dueDate := req.GetString("due_date", "")

	m, err := h.client.CreateMilestone(ctx, title, description, dueDate)
	if err != nil {
		return h.errH.Wrap("create milestone", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Created milestone #%d: %s", m.IID, m.Title)), nil
}

// SearchBranches finds branches matching a pattern.
func (h *Handlers) SearchBranches(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitLab not configured."), nil
	}
	pattern, err := req.RequireString("pattern")
	if err != nil {
		return mcputil.ErrInvalid("pattern parameter is required"), nil
	}

	branches, err := h.client.SearchBranches(ctx, pattern)
	if err != nil {
		return h.errH.Wrap("search branches", err), nil
	}
	if len(branches) == 0 {
		return mcputil.TextResult(fmt.Sprintf("No branches matching %q", pattern)), nil
	}
	var b strings.Builder
	for _, br := range branches {
		merged := ""
		if br.Merged {
			merged = " [merged]"
		}
		b.WriteString(br.Name + merged + "\n")
	}
	return mcputil.TextResult(b.String()), nil
}

// SearchMRsByBranch finds merge requests for a source branch.
func (h *Handlers) SearchMRsByBranch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitLab not configured."), nil
	}
	branch, err := req.RequireString("branch")
	if err != nil {
		return mcputil.ErrInvalid("branch parameter is required"), nil
	}

	mrs, err := h.client.SearchMRsByBranch(ctx, branch)
	if err != nil {
		return h.errH.Wrap("search MRs by branch", err), nil
	}
	if len(mrs) == 0 {
		return mcputil.TextResult(fmt.Sprintf("No MRs found for branch %q", branch)), nil
	}
	var b strings.Builder
	for _, mr := range mrs {
		b.WriteString(fmt.Sprintf("!%d %s [%s] by %s → %s\n",
			mr.IID, mr.Title, mr.State, mr.Author, mr.TargetBranch))
	}
	return mcputil.TextResult(b.String()), nil
}

// GetFileContent reads a file from the repository.
func (h *Handlers) GetFileContent(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitLab not configured."), nil
	}
	filePath, err := req.RequireString("path")
	if err != nil {
		return mcputil.ErrInvalid("path parameter is required"), nil
	}
	ref := req.GetString("ref", "main")

	content, err := h.client.GetFileContent(ctx, filePath, ref)
	if err != nil {
		return h.errH.Wrap("get file content", err), nil
	}
	return mcputil.TextResult(content), nil
}

// ListFiles lists directory contents in the repository.
func (h *Handlers) ListFiles(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitLab not configured."), nil
	}
	path := req.GetString("path", "")
	ref := req.GetString("ref", "main")

	files, err := h.client.ListFiles(ctx, path, ref)
	if err != nil {
		return h.errH.Wrap("list files", err), nil
	}
	if len(files) == 0 {
		return mcputil.TextResult("No files found."), nil
	}
	return mcputil.TextResult(strings.Join(files, "\n")), nil
}
