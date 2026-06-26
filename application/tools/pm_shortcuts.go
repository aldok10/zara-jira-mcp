package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	domain "github.com/aldok10/zara-jira-mcp/domain/jira"
	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// PMQuickStatus is THE one tool a PM needs to start their day.
func (h *Handlers) PMQuickStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)
	if boardID == 0 {
		boards, _ := h.Jira.GetBoards(ctx)
		if len(boards) > 0 {
			boardID = boards[0].ID
		}
	}

	var sb strings.Builder
	sb.WriteString("Project Status:\n\n")

	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			sprint := sprints[0]
			issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)
			var done, blocked, total int
			total = len(issues)
			for _, i := range issues {
				l := strings.ToLower(i.Status)
				if strings.Contains(l, "done") || strings.Contains(l, "closed") {
					done++
				} else if strings.Contains(l, "block") {
					blocked++
				}
			}
			pct := 0.0
			if total > 0 {
				pct = float64(done) / float64(total) * 100
			}
			sb.WriteString(fmt.Sprintf("SPRINT: %s — %d/%d done (%.0f%%), %d blocked\n", sprint.Name, done, total, pct, blocked))
			if sprint.Goal != "" {
				sb.WriteString(fmt.Sprintf("  Goal: %s\n", sprint.Goal))
			}
			sb.WriteString("\n")
		}
	}

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		sb.WriteString(fmt.Sprintf("BLOCKERS: %d active\n", len(blockers)))
		for _, b := range blockers {
			sb.WriteString(fmt.Sprintf("  - [%dd] %s\n", int(time.Since(b.BlockedSince).Hours()/24), b.Description))
		}
		sb.WriteString("\n")
	}

	risks, _ := h.Memory.GetOpenRisks(ctx)
	high := 0
	for _, r := range risks {
		if r.Severity == "critical" || r.Severity == "high" {
			high++
		}
	}
	if high > 0 {
		sb.WriteString(fmt.Sprintf("RISKS: %d high/critical\n\n", high))
	}

	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 0 {
		sb.WriteString(fmt.Sprintf("PENDING RETRO ACTIONS: %d\n\n", len(actions)))
	}

	return textResult(sb.String()), nil
}

// PMCreate creates work on Jira, GitHub, or GitLab with minimal params.
func (h *Handlers) PMCreate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}
	desc := req.GetString("description", "")
	labels := req.GetString("labels", "")
	assignee := req.GetString("assignee", "")
	platform := req.GetString("platform", "jira")

	switch platform {
	case "github":
		if !h.GitHub.Available() {
			return errorResult("GitHub not configured"), nil
		}
		var ll, aa []string
		if labels != "" {
			ll = strings.Split(labels, ",")
		}
		if assignee != "" {
			aa = []string{assignee}
		}
		iss, err := h.GitHub.CreateIssue(ctx, title, desc, ll, aa, 0)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return textResult(fmt.Sprintf("GitHub #%d: %s", iss.Number, iss.Title)), nil

	case "gitlab":
		if !h.GitLab.Available() {
			return errorResult("GitLab not configured"), nil
		}
		var ll []string
		if labels != "" {
			ll = strings.Split(labels, ",")
		}
		iss, err := h.GitLab.CreateIssue(ctx, title, desc, ll, 0, 0)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return textResult(fmt.Sprintf("GitLab #%d: %s\n%s", iss.IID, iss.Title, iss.WebURL)), nil

	default:
		project := req.GetString("project", "")
		if project == "" {
			return errorResult("project required for Jira"), nil
		}
		input := &domain.CreateIssueInput{
			Project: project, Summary: title, Description: desc,
			Priority: req.GetString("priority", "Medium"),
			IssueType: req.GetString("type", "Task"), Assignee: assignee,
		}
		if labels != "" {
			input.Labels = strings.Split(labels, ",")
		}
		created, err := h.Jira.CreateIssue(ctx, input)
		if err != nil {
			return errorResult(err.Error()), nil
		}
		return textResult(fmt.Sprintf("Jira %s: %s", created.Key, created.Summary)), nil
	}
}

// PMDecide records a decision quickly.
func (h *Handlers) PMDecide(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	what, err := req.RequireString("what")
	if err != nil {
		return errorResult("what required"), nil
	}
	d := &memdom.Decision{
		Title: what, Decision: what, Rationale: req.GetString("why", ""),
		MadeBy: req.GetString("who", "team"), MadeAt: time.Now(), Tags: "quick",
	}
	if err := h.Memory.SaveDecision(ctx, d); err != nil {
		return errorResult(err.Error()), nil
	}
	return textResult(fmt.Sprintf("Decision recorded: %s", what)), nil
}

// PMRisk records a risk quickly.
func (h *Handlers) PMRisk(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	what, err := req.RequireString("what")
	if err != nil {
		return errorResult("what required"), nil
	}
	r := &memdom.Risk{
		Title: what, Severity: req.GetString("severity", "medium"),
		Status: "open", Owner: req.GetString("owner", ""), IdentifiedAt: time.Now(),
	}
	if err := h.Memory.SaveRisk(ctx, r); err != nil {
		return errorResult(err.Error()), nil
	}
	return textResult(fmt.Sprintf("Risk recorded: [%s] %s", r.Severity, what)), nil
}

// PMNext suggests what the PM should do next.
func (h *Handlers) PMNext(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var suggestions []string

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	for _, b := range blockers {
		if time.Since(b.BlockedSince).Hours() > 72 {
			suggestions = append(suggestions, fmt.Sprintf("Escalate: '%s' blocked %d days", b.Description, int(time.Since(b.BlockedSince).Hours()/24)))
		}
	}

	risks, _ := h.Memory.GetOpenRisks(ctx)
	for _, r := range risks {
		if r.Severity == "critical" && r.Mitigation == "" {
			suggestions = append(suggestions, fmt.Sprintf("Mitigate: risk '%s' has no plan", r.Title))
		}
	}

	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 5 {
		suggestions = append(suggestions, fmt.Sprintf("Clean up: %d retro actions piling up", len(actions)))
	}

	if len(suggestions) == 0 {
		return textResult("All clear. Focus on sprint goal delivery or check in with team members."), nil
	}

	var sb strings.Builder
	sb.WriteString("Next actions:\n")
	for i, s := range suggestions {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, s))
	}
	return textResult(sb.String()), nil
}

// Ensure unused imports don't fail
var _ = domain.CreateIssueInput{}
var _ memdom.Risk
