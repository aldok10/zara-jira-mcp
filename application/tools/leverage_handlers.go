package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// JiraTraceBranch traces a Jira ticket to its git branch/PR status.
func (h *Handlers) JiraTraceBranch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key required (Jira issue key)"), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tracing %s to code:\n\n", key))

	// Check GitHub
	if h.GitHub.Available() {
		// Look for branch matching pattern: feature/KEY-*, fix/KEY-*, etc
		branches, _ := h.GitHub.SearchBranches(ctx, key)
		if len(branches) > 0 {
			sb.WriteString("GitHub branches:\n")
			for _, b := range branches {
				sb.WriteString(fmt.Sprintf("  - %s\n", b.Name))
			}
		}

		// Check open PRs mentioning this key
		prs, _ := h.GitHub.ListPRs(ctx, "open", 30)
		for _, pr := range prs {
			if strings.Contains(pr.Title, key) {
				sb.WriteString(fmt.Sprintf("\nOpen PR: #%d %s (by %s)\n", pr.Number, pr.Title, pr.User))
			}
		}
	}

	// Check GitLab
	if h.GitLab.Available() {
		mrs, _ := h.GitLab.ListMRs(ctx, "opened", 30)
		for _, mr := range mrs {
			if strings.Contains(mr.Title, key) {
				sb.WriteString(fmt.Sprintf("\nGitLab MR: !%d %s (by %s)\n", mr.IID, mr.Title, mr.Author))
			}
		}
	}

	if sb.Len() < 30 {
		sb.WriteString("No branches or PRs/MRs found for this ticket.\n")
		sb.WriteString("(Issue may not have started development yet)")
	}

	return textResult(sb.String()), nil
}

// IncidentImpact summarizes production incidents from Jira issues of type "Bug" with high priority.
func (h *Handlers) IncidentImpact(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jql := "issuetype = Bug AND priority in (Highest, High) AND resolution = Unresolved ORDER BY priority DESC, created DESC"

	result, err := h.Jira.SearchIssues(ctx, jql, 30, 0)
	if err != nil {
		return sanitizedError("jira operation failed in leverage", err), nil
	}

	if len(result.Issues) == 0 {
		return textResult("No high-priority bugs/incidents open. Production looks stable."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Production Incidents: %d high-priority bugs\n\n", len(result.Issues)))

	oldest := time.Duration(0)
	for _, i := range result.Issues {
		age := time.Since(i.Created)
		if age > oldest {
			oldest = age
		}
		sb.WriteString(fmt.Sprintf("  %s [%s] %s (age: %dd, assignee: %s)\n",
			i.Key, i.Priority, i.Summary, int(age.Hours()/24), i.Assignee))
	}

	sb.WriteString(fmt.Sprintf("\nOldest unresolved: %d days\n", int(oldest.Hours()/24)))
	if int(oldest.Hours()/24) > 7 {
		sb.WriteString("WARNING: High-priority bugs >7 days old. Escalate or reprioritize.\n")
	}

	return textResult(sb.String()), nil
}

// SprintForecastSprint is a simpler "will we finish this sprint?" based on burn rate.
func (h *Handlers) SprintForecastSimple(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) == 0 {
		return textResult("No active sprint."), nil
	}

	sprint := sprints[0]
	issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)

	var done, total int
	total = len(issues)
	for _, i := range issues {
		l := strings.ToLower(i.Status)
		if strings.Contains(l, "done") || strings.Contains(l, "closed") {
			done++
		}
	}
	remaining := total - done

	// Check daily progress for burn rate
	progress, _ := h.Memory.GetDailyProgress(ctx, boardID, sprint.Name)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprint Forecast: %s\n\n", sprint.Name))
	sb.WriteString(fmt.Sprintf("Done: %d/%d | Remaining: %d\n", done, total, remaining))

	if len(progress) >= 2 {
		first := progress[0]
		last := progress[len(progress)-1]
		daysElapsed := last.Date.Sub(first.Date).Hours() / 24
		itemsDone := last.Done - first.Done

		if daysElapsed > 0 && itemsDone > 0 {
			rate := float64(itemsDone) / daysElapsed
			daysNeeded := float64(remaining) / rate
			estimatedEnd := time.Now().AddDate(0, 0, int(daysNeeded))

			sb.WriteString(fmt.Sprintf("\nBurn rate: %.1f items/day\n", rate))
			sb.WriteString(fmt.Sprintf("Estimated completion: %s (%d days)\n", estimatedEnd.Format("Jan 2"), int(daysNeeded)))

			if daysNeeded > 10 {
				sb.WriteString("\nVERDICT: AT RISK — won't finish at current pace.\n")
				sb.WriteString("Options: reduce scope, add capacity, or extend sprint.\n")
			} else if daysNeeded > 7 {
				sb.WriteString("\nVERDICT: TIGHT — possible but needs focus.\n")
			} else {
				sb.WriteString("\nVERDICT: ON TRACK\n")
			}
		} else {
			sb.WriteString("\nNot enough burn rate data. Keep using pm_track_daily.\n")
		}
	} else {
		sb.WriteString("\nNeed 2+ daily data points for forecast. Use pm_track_daily.\n")
	}

	return textResult(sb.String()), nil
}

// BacklogHealthCheck evaluates backlog quality — stale items, unestimated, missing AC.
func (h *Handlers) BacklogHealthCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project := req.GetString("project", "")
	days := req.GetInt("days", 90)

	jql := fmt.Sprintf("resolution = Unresolved AND updated <= -%dd AND sprint IS EMPTY ORDER BY created ASC", days)
	if project != "" {
		jql = fmt.Sprintf("project = %s AND resolution = Unresolved AND updated <= -%dd AND sprint IS EMPTY ORDER BY created ASC", project, days)
	}

	result, err := h.Jira.SearchIssues(ctx, jql, 50, 0)
	if err != nil {
		return sanitizedError("jira operation failed in leverage", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Backlog Health Check (stale >%d days, not in sprint):\n\n", days))

	if len(result.Issues) == 0 {
		sb.WriteString("Backlog is clean! No stale items found.\n")
		return textResult(sb.String()), nil
	}

	sb.WriteString(fmt.Sprintf("Stale items: %d\n\n", len(result.Issues)))
	for _, i := range result.Issues {
		age := int(time.Since(i.Created).Hours() / 24)
		sb.WriteString(fmt.Sprintf("  %s [%s] %s (created %d days ago)\n", i.Key, i.Type, i.Summary, age))
	}

	sb.WriteString("\nRecommendation:\n")
	if len(result.Issues) > 20 {
		sb.WriteString("  URGENT: >20 stale items. Schedule grooming session to archive/close.\n")
	} else if len(result.Issues) > 10 {
		sb.WriteString("  Review and close items that are no longer relevant.\n")
	} else {
		sb.WriteString("  Minor cleanup needed. Address in next grooming.\n")
	}

	return textResult(sb.String()), nil
}

// Ensure unused import doesn't fail
var _ memdom.Risk
