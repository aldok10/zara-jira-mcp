package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/github"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

// Handlers exposes GitHub operations as MCP tool handlers.
type Handlers struct {
	client *github.Client
	errH   *mcputil.ErrorHandler
}

func NewHandlers(client *github.Client) *Handlers {
	return &Handlers{
		client: client,
		errH:   mcputil.NewErrorHandler(nil),
	}
}

// ListPRs lists open pull requests with review status, age, and assignees.
func (h *Handlers) ListPRs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitHub not configured. Set GITHUB_TOKEN, GITHUB_OWNER, GITHUB_REPO."), nil
	}
	state := req.GetString("state", "open")
	limit := req.GetInt("limit", 30)

	prs, err := h.client.ListPRs(ctx, state, limit)
	if err != nil {
		return h.errH.Wrap("list PRs", err), nil
	}
	if len(prs) == 0 {
		return mcputil.TextResult("No pull requests found."), nil
	}

	var b strings.Builder
	for _, pr := range prs {
		age := time.Since(pr.CreatedAt).Truncate(time.Hour)
		reviewers := "none"
		if len(pr.Reviewers) > 0 {
			reviewers = strings.Join(pr.Reviewers, ", ")
		}
		draft := ""
		if pr.Draft {
			draft = " [DRAFT]"
		}
		b.WriteString(fmt.Sprintf("#%d%s %s (%s)\n  Author: %s | Age: %s | Reviewers: %s\n",
			pr.Number, draft, pr.Title, pr.State, pr.User, age, reviewers))
	}
	return mcputil.TextResult(b.String()), nil
}

// ListReleases lists recent releases/tags.
func (h *Handlers) ListReleases(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitHub not configured."), nil
	}
	limit := req.GetInt("limit", 10)

	releases, err := h.client.ListReleases(ctx, limit)
	if err != nil {
		return h.errH.Wrap("list releases", err), nil
	}
	if len(releases) == 0 {
		return mcputil.TextResult("No releases found."), nil
	}
	var b strings.Builder
	for _, r := range releases {
		b.WriteString(fmt.Sprintf("%s — %s (%s by %s)\n", r.TagName, r.Name, r.PublishedAt.Format("2006-01-02"), r.Author))
	}
	return mcputil.TextResult(b.String()), nil
}

// GetActivity returns repo activity summary (commits, PRs merged, issues closed).
func (h *Handlers) GetActivity(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitHub not configured."), nil
	}
	days := req.GetInt("days", 7)

	act, err := h.client.GetActivity(ctx, days)
	if err != nil {
		return h.errH.Wrap("get activity", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Last %d days:\n  Commits: %d\n  PRs merged: %d\n  Issues closed: %d",
		days, act.CommitCount, act.PRsMerged, act.IssuesClosed)), nil
}

// SearchBranches finds branches matching a pattern (e.g. issue key).
func (h *Handlers) SearchBranches(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitHub not configured."), nil
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
		b.WriteString(br.Name + "\n")
	}
	return mcputil.TextResult(b.String()), nil
}

// SearchPRsByBranch finds PRs for a branch.
func (h *Handlers) SearchPRsByBranch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitHub not configured."), nil
	}
	branch, err := req.RequireString("branch")
	if err != nil {
		return mcputil.ErrInvalid("branch parameter is required"), nil
	}

	prs, err := h.client.SearchPRsByBranch(ctx, branch)
	if err != nil {
		return h.errH.Wrap("search PRs by branch", err), nil
	}
	if len(prs) == 0 {
		return mcputil.TextResult(fmt.Sprintf("No PRs found for branch %q", branch)), nil
	}
	var b strings.Builder
	for _, pr := range prs {
		target := ""
		if pr.MergedTo != "" {
			target = fmt.Sprintf(" → %s", pr.MergedTo)
		}
		b.WriteString(fmt.Sprintf("#%d %s [%s]%s by %s\n", pr.Number, pr.Title, pr.State, target, pr.Author))
	}
	return mcputil.TextResult(b.String()), nil
}

// CreateIssue creates a GitHub issue.
func (h *Handlers) CreateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitHub not configured."), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return mcputil.ErrInvalid("title parameter is required"), nil
	}
	body := req.GetString("body", "")
	labelsStr := req.GetString("labels", "")
	assigneesStr := req.GetString("assignees", "")
	milestone := req.GetInt("milestone", 0)

	var labels, assignees []string
	if labelsStr != "" {
		labels = strings.Split(labelsStr, ",")
		for i := range labels {
			labels[i] = strings.TrimSpace(labels[i])
		}
	}
	if assigneesStr != "" {
		assignees = strings.Split(assigneesStr, ",")
		for i := range assignees {
			assignees[i] = strings.TrimSpace(assignees[i])
		}
	}

	issue, err := h.client.CreateIssue(ctx, title, body, labels, assignees, milestone)
	if err != nil {
		return h.errH.Wrap("create issue", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Created issue #%d: %s [%s]", issue.Number, issue.Title, issue.State)), nil
}

// ListIssues lists GitHub issues.
func (h *Handlers) ListIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitHub not configured."), nil
	}
	state := req.GetString("state", "open")
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
		b.WriteString(fmt.Sprintf("#%d %s [%s] — %s\n", iss.Number, iss.Title, iss.State, assignee))
	}
	return mcputil.TextResult(b.String()), nil
}

// CreateMilestone creates a GitHub milestone.
func (h *Handlers) CreateMilestone(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitHub not configured."), nil
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
	return mcputil.TextResult(fmt.Sprintf("Created milestone #%d: %s", m.Number, m.Title)), nil
}

// ListMilestones lists GitHub milestones with progress.
func (h *Handlers) ListMilestones(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitHub not configured."), nil
	}
	state := req.GetString("state", "open")

	milestones, err := h.client.ListMilestones(ctx, state)
	if err != nil {
		return h.errH.Wrap("list milestones", err), nil
	}
	if len(milestones) == 0 {
		return mcputil.TextResult("No milestones found."), nil
	}
	var b strings.Builder
	for _, m := range milestones {
		b.WriteString(fmt.Sprintf("#%d %s [%s] — %d open / %d closed",
			m.Number, m.Title, m.State, m.OpenIssues, m.ClosedIssues))
		if m.DueOn != "" {
			b.WriteString(fmt.Sprintf(" (due: %s)", m.DueOn[:10]))
		}
		b.WriteString("\n")
	}
	return mcputil.TextResult(b.String()), nil
}

// GetPRMetrics shows PR aging metrics.
func (h *Handlers) GetPRMetrics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("GitHub not configured."), nil
	}
	staleDays := req.GetInt("stale_days", 7)

	prs, err := h.client.ListPRs(ctx, "open", 100)
	if err != nil {
		return h.errH.Wrap("list PRs for metrics", err), nil
	}
	if len(prs) == 0 {
		return mcputil.TextResult("No open PRs."), nil
	}

	var totalAge time.Duration
	staleCount := 0
	staleThreshold := time.Duration(staleDays) * 24 * time.Hour
	for _, pr := range prs {
		age := time.Since(pr.CreatedAt)
		totalAge += age
		if age > staleThreshold {
			staleCount++
		}
	}

	avgAge := totalAge / time.Duration(len(prs))
	return mcputil.TextResult(fmt.Sprintf(
		"Open PRs: %d\nAvg age: %s\nStale (>%dd): %d\nStale rate: %d%%",
		len(prs), avgAge.Truncate(time.Hour), staleDays, staleCount,
		staleCount*100/len(prs))), nil
}
