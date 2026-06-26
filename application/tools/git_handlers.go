package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// GitHubCreateIssue creates a GitHub issue.
func (h *Handlers) GitHubCreateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitHub.Available() {
		return errorResult("GitHub not configured. Set GITHUB_TOKEN, GITHUB_OWNER, GITHUB_REPO."), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}

	body := req.GetString("body", "")
	labelsRaw := req.GetString("labels", "")
	assigneesRaw := req.GetString("assignees", "")
	milestone := req.GetInt("milestone", 0)

	var labels, assignees []string
	if labelsRaw != "" {
		labels = strings.Split(labelsRaw, ",")
	}
	if assigneesRaw != "" {
		assignees = strings.Split(assigneesRaw, ",")
	}

	issue, err := h.GitHub.CreateIssue(ctx, title, body, labels, assignees, milestone)
	if err != nil {
		return errorResult("GitHub error: " + err.Error()), nil
	}

	return textResult(fmt.Sprintf("GitHub issue created: #%d %s", issue.Number, issue.Title)), nil
}

// GitHubListIssues lists GitHub issues.
func (h *Handlers) GitHubListIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitHub.Available() {
		return errorResult("GitHub not configured."), nil
	}

	state := req.GetString("state", "open")
	labels := req.GetString("labels", "")
	limit := req.GetInt("limit", 20)

	issues, err := h.GitHub.ListIssues(ctx, state, labels, limit)
	if err != nil {
		return errorResult("GitHub error: " + err.Error()), nil
	}

	if len(issues) == 0 {
		return textResult("No issues found."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("GitHub Issues (%d):\n\n", len(issues)))
	for _, i := range issues {
		labels := ""
		if len(i.Labels) > 0 {
			labels = " [" + strings.Join(i.Labels, ",") + "]"
		}
		sb.WriteString(fmt.Sprintf("#%d %s%s (assignee: %s)\n", i.Number, i.Title, labels, i.Assignee))
	}
	return textResult(sb.String()), nil
}

// GitHubCreateMilestone creates a GitHub milestone.
func (h *Handlers) GitHubCreateMilestone(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitHub.Available() {
		return errorResult("GitHub not configured."), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}

	m, err := h.GitHub.CreateMilestone(ctx, title, req.GetString("description", ""), req.GetString("due_date", ""))
	if err != nil {
		return errorResult("GitHub error: " + err.Error()), nil
	}

	return textResult(fmt.Sprintf("Milestone created: #%d %s", m.Number, m.Title)), nil
}

// GitHubListMilestones lists milestones.
func (h *Handlers) GitHubListMilestones(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitHub.Available() {
		return errorResult("GitHub not configured."), nil
	}

	milestones, err := h.GitHub.ListMilestones(ctx, req.GetString("state", "open"))
	if err != nil {
		return errorResult("GitHub error: " + err.Error()), nil
	}

	if len(milestones) == 0 {
		return textResult("No milestones."), nil
	}

	var sb strings.Builder
	for _, m := range milestones {
		sb.WriteString(fmt.Sprintf("#%d %s [%s] (open: %d, closed: %d) due: %s\n",
			m.Number, m.Title, m.State, m.OpenIssues, m.ClosedIssues, m.DueOn))
	}
	return textResult(sb.String()), nil
}

// GitHubReadFile reads a file from the repo.
func (h *Handlers) GitHubReadFile(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitHub.Available() {
		return errorResult("GitHub not configured."), nil
	}
	path, err := req.RequireString("path")
	if err != nil {
		return errorResult("path required"), nil
	}

	content, err := h.GitHub.GetFileContent(ctx, path, req.GetString("ref", ""))
	if err != nil {
		return errorResult("GitHub error: " + err.Error()), nil
	}

	return textResult(content), nil
}

// GitHubListFiles lists files in a directory.
func (h *Handlers) GitHubListFiles(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitHub.Available() {
		return errorResult("GitHub not configured."), nil
	}

	files, err := h.GitHub.ListFiles(ctx, req.GetString("path", ""), req.GetString("ref", ""))
	if err != nil {
		return errorResult("GitHub error: " + err.Error()), nil
	}

	return textResult(strings.Join(files, "\n")), nil
}

// GitLabCreateIssue creates a GitLab issue.
func (h *Handlers) GitLabCreateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitLab.Available() {
		return errorResult("GitLab not configured. Set GITLAB_TOKEN, GITLAB_PROJECT_ID."), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}

	var labels []string
	if raw := req.GetString("labels", ""); raw != "" {
		labels = strings.Split(raw, ",")
	}

	issue, err := h.GitLab.CreateIssue(ctx, title, req.GetString("description", ""), labels,
		req.GetInt("assignee_id", 0), req.GetInt("milestone_id", 0))
	if err != nil {
		return errorResult("GitLab error: " + err.Error()), nil
	}

	return textResult(fmt.Sprintf("GitLab issue created: #%d %s\n%s", issue.IID, issue.Title, issue.WebURL)), nil
}

// GitLabListIssues lists GitLab issues.
func (h *Handlers) GitLabListIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitLab.Available() {
		return errorResult("GitLab not configured."), nil
	}

	issues, err := h.GitLab.ListIssues(ctx, req.GetString("state", "opened"), req.GetString("labels", ""), req.GetInt("limit", 20))
	if err != nil {
		return errorResult("GitLab error: " + err.Error()), nil
	}

	if len(issues) == 0 {
		return textResult("No issues found."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("GitLab Issues (%d):\n\n", len(issues)))
	for _, i := range issues {
		labels := ""
		if len(i.Labels) > 0 {
			labels = " [" + strings.Join(i.Labels, ",") + "]"
		}
		sb.WriteString(fmt.Sprintf("#%d %s%s (assignee: %s)\n", i.IID, i.Title, labels, i.Assignee))
	}
	return textResult(sb.String()), nil
}

// GitLabCreateMilestone creates a GitLab milestone.
func (h *Handlers) GitLabCreateMilestone(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitLab.Available() {
		return errorResult("GitLab not configured."), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}

	m, err := h.GitLab.CreateMilestone(ctx, title, req.GetString("description", ""), req.GetString("due_date", ""))
	if err != nil {
		return errorResult("GitLab error: " + err.Error()), nil
	}

	return textResult(fmt.Sprintf("GitLab milestone created: #%d %s", m.IID, m.Title)), nil
}

// GitLabListMilestones lists GitLab milestones.
func (h *Handlers) GitLabListMilestones(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitLab.Available() {
		return errorResult("GitLab not configured."), nil
	}

	milestones, err := h.GitLab.ListMilestones(ctx, req.GetString("state", "active"))
	if err != nil {
		return errorResult("GitLab error: " + err.Error()), nil
	}

	if len(milestones) == 0 {
		return textResult("No milestones."), nil
	}

	var sb strings.Builder
	for _, m := range milestones {
		sb.WriteString(fmt.Sprintf("#%d %s [%s] due: %s\n", m.IID, m.Title, m.State, m.DueDate))
	}
	return textResult(sb.String()), nil
}

// GitLabListMRs lists merge requests.
func (h *Handlers) GitLabListMRs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitLab.Available() {
		return errorResult("GitLab not configured."), nil
	}

	mrs, err := h.GitLab.ListMRs(ctx, req.GetString("state", "opened"), req.GetInt("limit", 20))
	if err != nil {
		return errorResult("GitLab error: " + err.Error()), nil
	}

	if len(mrs) == 0 {
		return textResult("No merge requests."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("GitLab MRs (%d):\n\n", len(mrs)))
	for _, mr := range mrs {
		draft := ""
		if mr.Draft {
			draft = " [DRAFT]"
		}
		sb.WriteString(fmt.Sprintf("!%d %s%s (by: %s, state: %s)\n", mr.IID, mr.Title, draft, mr.Author, mr.State))
	}
	return textResult(sb.String()), nil
}

// GitLabReadFile reads a file from GitLab repo.
func (h *Handlers) GitLabReadFile(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitLab.Available() {
		return errorResult("GitLab not configured."), nil
	}
	path, err := req.RequireString("path")
	if err != nil {
		return errorResult("path required"), nil
	}

	content, err := h.GitLab.GetFileContent(ctx, path, req.GetString("ref", ""))
	if err != nil {
		return errorResult("GitLab error: " + err.Error()), nil
	}

	return textResult(content), nil
}

// GitLabListFiles lists files in GitLab repo.
func (h *Handlers) GitLabListFiles(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.GitLab.Available() {
		return errorResult("GitLab not configured."), nil
	}

	files, err := h.GitLab.ListFiles(ctx, req.GetString("path", ""), req.GetString("ref", ""))
	if err != nil {
		return errorResult("GitLab error: " + err.Error()), nil
	}

	return textResult(strings.Join(files, "\n")), nil
}
