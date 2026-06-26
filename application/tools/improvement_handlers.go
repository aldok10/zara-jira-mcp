package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// PMImprovementDashboard shows meta-metrics: is the team getting better at getting better?
func (h *Handlers) PMImprovementDashboard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	var sb strings.Builder
	sb.WriteString("## Improvement Dashboard\n\n")

	// 1. Velocity trend
	var velocityTrend string
	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 6)
		if len(snaps) >= 2 {
			recent := snaps[0].Velocity
			older := snaps[len(snaps)-1].Velocity
			switch {
			case recent > older:
				velocityTrend = "IMPROVING"
			case recent < older:
				velocityTrend = "DECLINING"
			default:
				velocityTrend = "STABLE"
			}
			sb.WriteString(fmt.Sprintf("**Velocity Trend:** %s (%d -> %d over %d sprints)\n", velocityTrend, older, recent, len(snaps)))
		} else {
			sb.WriteString("**Velocity Trend:** Not enough data (need 2+ snapshots)\n")
		}
	} else {
		sb.WriteString("**Velocity Trend:** No board_id provided\n")
	}

	// 2. Action item completion rate
	pending, _ := h.Memory.GetPendingActionItems(ctx)
	allBlockers, _ := h.Memory.GetBlockerHistory(ctx, 50)
	retros, _ := h.Memory.GetRetrospectives(ctx, 10)

	totalActions := len(pending)
	for _, r := range retros {
		if r.ActionItems != "" {
			totalActions += len(strings.Split(r.ActionItems, "\n"))
		}
	}
	completedActions := totalActions - len(pending)
	completionRate := 0.0
	if totalActions > 0 {
		completionRate = float64(completedActions) / float64(totalActions) * 100
	}
	sb.WriteString(fmt.Sprintf("**Action Item Completion:** %.0f%% (%d/%d)\n", completionRate, completedActions, totalActions))

	// 3. Sprint goal hit rate
	goalDecisions, _ := h.Memory.SearchDecisions(ctx, "sprint_goal_result")
	hits := 0
	for _, d := range goalDecisions {
		if strings.Contains(strings.ToLower(d.Decision), "hit") || strings.Contains(strings.ToLower(d.Decision), "achieved") {
			hits++
		}
	}
	if len(goalDecisions) > 0 {
		sb.WriteString(fmt.Sprintf("**Sprint Goal Hit Rate:** %d/%d (%.0f%%)\n", hits, len(goalDecisions), float64(hits)/float64(len(goalDecisions))*100))
	} else {
		sb.WriteString("**Sprint Goal Hit Rate:** No data (use pm_sprint_goal_track to record)\n")
	}

	// 4. Retro cadence
	if len(retros) >= 2 {
		interval := retros[0].Date.Sub(retros[1].Date).Hours() / 24
		sb.WriteString(fmt.Sprintf("**Retro Cadence:** %.0f days between last 2 retros\n", interval))
	} else {
		sb.WriteString("**Retro Cadence:** Not enough retros recorded\n")
	}

	// 5. Blocker resolution speed
	var totalResolutionDays int
	resolvedCount := 0
	for _, b := range allBlockers {
		if b.ResolvedAt != nil {
			days := int(b.ResolvedAt.Sub(b.BlockedSince).Hours() / 24)
			totalResolutionDays += days
			resolvedCount++
		}
	}
	if resolvedCount > 0 {
		avg := totalResolutionDays / resolvedCount
		sb.WriteString(fmt.Sprintf("**Avg Blocker Resolution:** %d days (%d resolved)\n", avg, resolvedCount))
	}

	sb.WriteString("\n---\n")
	sb.WriteString("Tip: Record sprint goals with pm_sprint_goal_track, snapshot sprints with pm_snapshot_sprint, and run retros with pm_record_retro to fill gaps.\n")

	return textResult(sb.String()), nil
}

// PMBusFactor detects single points of failure in the team.
func (h *Handlers) PMBusFactor(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return textResult("No active sprint found."), nil
	}

	issues, err := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
	if err != nil {
		return sanitizedError("failed to get sprint issues", err), nil
	}

	if len(issues) == 0 {
		return textResult("No issues in current sprint."), nil
	}

	// Count by assignee
	assigneeCounts := map[string]int{}
	// Count by assignee+label/component
	assigneeLabels := map[string]map[string]int{}
	for _, issue := range issues {
		assignee := issue.Assignee
		if assignee == "" {
			assignee = "Unassigned"
		}
		assigneeCounts[assignee]++
		if assigneeLabels[assignee] == nil {
			assigneeLabels[assignee] = map[string]int{}
		}
		for _, label := range issue.Labels {
			assigneeLabels[assignee][label]++
		}
	}

	total := len(issues)
	var sb strings.Builder
	sb.WriteString("## Bus Factor Report\n\n")
	sb.WriteString(fmt.Sprintf("Sprint: %s (%d issues)\n\n", sprints[0].Name, total))

	// Check concentration
	sb.WriteString("### Workload Distribution\n")
	var risks []string
	for person, count := range assigneeCounts {
		pct := float64(count) / float64(total) * 100
		indicator := ""
		if pct > 50 {
			indicator = " [RISK: >50%]"
			risks = append(risks, fmt.Sprintf("%s owns %.0f%% of sprint items", person, pct))
		}
		sb.WriteString(fmt.Sprintf("- %s: %d issues (%.0f%%)%s\n", person, count, pct, indicator))
	}

	// Check label/area concentration
	sb.WriteString("\n### Area Concentration\n")
	labelOwnership := map[string]map[string]int{}
	for person, labels := range assigneeLabels {
		for label, count := range labels {
			if labelOwnership[label] == nil {
				labelOwnership[label] = map[string]int{}
			}
			labelOwnership[label][person] = count
		}
	}

	for label, owners := range labelOwnership {
		totalInLabel := 0
		for _, c := range owners {
			totalInLabel += c
		}
		for person, count := range owners {
			if totalInLabel >= 3 && float64(count)/float64(totalInLabel) > 0.5 {
				risks = append(risks, fmt.Sprintf("%s owns %d/%d items in '%s'", person, count, totalInLabel, label))
				sb.WriteString(fmt.Sprintf("- [RISK] %s: %d/%d in label '%s'\n", person, count, totalInLabel, label))
			}
		}
	}

	if len(risks) == 0 {
		sb.WriteString("\nNo bus factor risks detected. Workload is distributed.\n")
	} else {
		sb.WriteString(fmt.Sprintf("\n### Risks Found: %d\n", len(risks)))
		for _, r := range risks {
			sb.WriteString(fmt.Sprintf("- %s\n", r))
		}
		sb.WriteString("\nRecommendation: Pair up on concentrated areas. Cross-train before the bus arrives.\n")
	}

	return textResult(sb.String()), nil
}

// PMAsyncStandup generates async standup from Jira transitions and blockers.
func (h *Handlers) PMAsyncStandup(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return textResult("No active sprint found."), nil
	}

	issues, err := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
	if err != nil {
		return sanitizedError("failed to get sprint issues", err), nil
	}

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	// Group by assignee
	type memberUpdate struct {
		moved   []string
		planned []string
	}
	members := map[string]*memberUpdate{}

	for _, issue := range issues {
		assignee := issue.Assignee
		if assignee == "" {
			assignee = "Unassigned"
		}
		if members[assignee] == nil {
			members[assignee] = &memberUpdate{}
		}

		lower := strings.ToLower(issue.Status)
		if issue.Updated.After(yesterday) {
			members[assignee].moved = append(members[assignee].moved, fmt.Sprintf("%s: %s [%s]", issue.Key, issue.Summary, issue.Status))
		} else if strings.Contains(lower, "progress") || strings.Contains(lower, "review") {
			members[assignee].planned = append(members[assignee].planned, fmt.Sprintf("%s: %s", issue.Key, issue.Summary))
		}
	}

	// Get active blockers
	blockers, _ := h.Memory.GetActiveBlockers(ctx)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Async Standup - %s\n", now.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("Sprint: %s\n\n", sprints[0].Name))

	for person, update := range members {
		if len(update.moved) == 0 && len(update.planned) == 0 {
			continue
		}
		sb.WriteString(fmt.Sprintf("### %s\n", person))
		if len(update.moved) > 0 {
			sb.WriteString("**Yesterday:**\n")
			for _, m := range update.moved {
				sb.WriteString(fmt.Sprintf("- %s\n", m))
			}
		}
		if len(update.planned) > 0 {
			sb.WriteString("**Today:**\n")
			for _, p := range update.planned {
				sb.WriteString(fmt.Sprintf("- %s\n", p))
			}
		}
		sb.WriteString("\n")
	}

	if len(blockers) > 0 {
		sb.WriteString("### Blockers\n")
		for _, b := range blockers {
			days := int(now.Sub(b.BlockedSince).Hours() / 24)
			sb.WriteString(fmt.Sprintf("- %s (owner: %s, %d days) %s\n", b.Description, b.Owner, days, b.IssueKey))
		}
	}

	return textResult(sb.String()), nil
}

// PMSprintGoalTrack records sprint goal outcome and shows hit rate history.
func (h *Handlers) PMSprintGoalTrack(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintName, err := req.RequireString("sprint_name")
	if err != nil {
		return errorResult("sprint_name required"), nil
	}

	hit := req.GetBool("hit", false)
	notes := req.GetString("notes", "")

	outcome := "missed"
	if hit {
		outcome = "hit"
	}

	decision := fmt.Sprintf("Sprint goal %s: %s", outcome, sprintName)
	if notes != "" {
		decision += " - " + notes
	}

	// Store as a decision with sprint_goal_result tag
	err = h.Memory.SaveDecision(ctx, &memdom.Decision{
		Title:    fmt.Sprintf("Sprint Goal Result: %s", sprintName),
		Decision: decision,
		Context:  fmt.Sprintf("Sprint: %s, Hit: %v", sprintName, hit),
		MadeBy:   "pm_sprint_goal_track",
		MadeAt:   time.Now(),
		Tags:     "sprint_goal_result",
	})
	if err != nil {
		return sanitizedError("failed to save improvement data", err), nil
	}

	// Show history
	goalDecisions, _ := h.Memory.SearchDecisions(ctx, "sprint_goal_result")
	totalGoals := len(goalDecisions)
	hits := 0
	for _, d := range goalDecisions {
		if strings.Contains(strings.ToLower(d.Decision), "hit") || strings.Contains(strings.ToLower(d.Decision), "achieved") {
			hits++
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Recorded: %s goal %s\n\n", sprintName, outcome))
	sb.WriteString(fmt.Sprintf("**Goal Hit Rate: %d/%d sprints (%.0f%%)**\n", hits, totalGoals, float64(hits)/float64(totalGoals)*100))

	if totalGoals >= 3 {
		if float64(hits)/float64(totalGoals) < 0.5 {
			sb.WriteString("\nWarning: Hit rate below 50%. Consider setting smaller, more achievable goals.\n")
		} else if float64(hits)/float64(totalGoals) >= 0.8 {
			sb.WriteString("\nStrong performance. Team is consistent at delivering on commitments.\n")
		}
	}

	return textResult(sb.String()), nil
}
