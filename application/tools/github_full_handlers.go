package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// GitHubPRs lists open PRs with review status and age.
func (h *Handlers) GitHubPRs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.GitHub == nil || !h.GitHub.Available() {
		return errorResult("GitHub not configured. Set GITHUB_TOKEN, GITHUB_OWNER, GITHUB_REPO."), nil
	}
	state := req.GetString("state", "open")
	limit := req.GetInt("limit", 30)

	prs, err := h.GitHub.ListPRs(ctx, state, limit)
	if err != nil {
		return sanitizedError("GitHub: failed to list PRs", err), nil
	}
	if len(prs) == 0 {
		return textResult("No PRs found."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pull Requests (%s): %d\n\n", state, len(prs)))
	now := time.Now()
	for _, pr := range prs {
		age := now.Sub(pr.CreatedAt)
		draft := ""
		if pr.Draft {
			draft = " [DRAFT]"
		}
		reviewers := "none"
		if len(pr.Reviewers) > 0 {
			reviewers = strings.Join(pr.Reviewers, ", ")
		}
		sb.WriteString(fmt.Sprintf("- #%d %s%s | by %s | age: %dd | reviewers: %s\n",
			pr.Number, pr.Title, draft, pr.User, int(age.Hours()/24), reviewers))
	}
	return textResult(sb.String()), nil
}

// GitHubPRMetrics shows PR metrics: avg merge time, stale PRs count.
func (h *Handlers) GitHubPRMetrics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.GitHub == nil || !h.GitHub.Available() {
		return errorResult("GitHub not configured. Set GITHUB_TOKEN, GITHUB_OWNER, GITHUB_REPO."), nil
	}

	prs, err := h.GitHub.ListPRs(ctx, "open", 100)
	if err != nil {
		return sanitizedError("GitHub: failed to list PRs for metrics", err), nil
	}

	now := time.Now()
	staleDays := req.GetInt("stale_days", 7)
	staleCount := 0
	var totalAge time.Duration
	for _, pr := range prs {
		age := now.Sub(pr.CreatedAt)
		totalAge += age
		if int(age.Hours()/24) > staleDays {
			staleCount++
		}
	}

	var sb strings.Builder
	sb.WriteString("GitHub PR Metrics:\n\n")
	sb.WriteString(fmt.Sprintf("- Open PRs: %d\n", len(prs)))
	if len(prs) > 0 {
		avgAge := totalAge / time.Duration(len(prs))
		sb.WriteString(fmt.Sprintf("- Avg age of open PRs: %.1f days\n", avgAge.Hours()/24))
	}
	sb.WriteString(fmt.Sprintf("- Stale PRs (>%d days): %d\n", staleDays, staleCount))

	return textResult(sb.String()), nil
}

// GitHubReleases lists recent releases/tags.
func (h *Handlers) GitHubReleases(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.GitHub == nil || !h.GitHub.Available() {
		return errorResult("GitHub not configured. Set GITHUB_TOKEN, GITHUB_OWNER, GITHUB_REPO."), nil
	}
	limit := req.GetInt("limit", 10)

	releases, err := h.GitHub.ListReleases(ctx, limit)
	if err != nil {
		return sanitizedError("GitHub: failed to list releases", err), nil
	}
	if len(releases) == 0 {
		return textResult("No releases found."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Recent releases (%d):\n\n", len(releases)))
	for _, r := range releases {
		name := r.Name
		if name == "" {
			name = r.TagName
		}
		sb.WriteString(fmt.Sprintf("- %s (%s) | by %s | %s\n",
			name, r.TagName, r.Author, r.PublishedAt.Format("2006-01-02")))
	}
	return textResult(sb.String()), nil
}

// GitHubActivity shows repo activity summary.
func (h *Handlers) GitHubActivity(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.GitHub == nil || !h.GitHub.Available() {
		return errorResult("GitHub not configured. Set GITHUB_TOKEN, GITHUB_OWNER, GITHUB_REPO."), nil
	}
	days := req.GetInt("days", 7)

	activity, err := h.GitHub.GetActivity(ctx, days)
	if err != nil {
		return sanitizedError("GitHub: failed to get activity", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Repo activity (last %d days):\n\n", days))
	sb.WriteString(fmt.Sprintf("- Commits: %d\n", activity.CommitCount))
	sb.WriteString(fmt.Sprintf("- PRs merged: %d\n", activity.PRsMerged))
	sb.WriteString(fmt.Sprintf("- Issues closed: %d\n", activity.IssuesClosed))
	return textResult(sb.String()), nil
}
