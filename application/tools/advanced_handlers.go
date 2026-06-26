package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// SprintHealthScore computes a numeric health score for the current sprint.
func (h *Handlers) SprintHealthScore(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil {
		return errorResult("Failed to get sprints: " + err.Error()), nil
	}
	if len(sprints) == 0 {
		return textResult("No active sprint."), nil
	}

	sprint := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprint.ID)
	if err != nil {
		return errorResult("Failed to get sprint issues: " + err.Error()), nil
	}

	total := len(issues)
	if total == 0 {
		return textResult("Sprint has no issues."), nil
	}

	var done, inProgress, blocked int
	for _, issue := range issues {
		lower := strings.ToLower(issue.Status)
		switch {
		case strings.Contains(lower, "done") || strings.Contains(lower, "closed") || strings.Contains(lower, "resolved"):
			done++
		case strings.Contains(lower, "progress") || strings.Contains(lower, "review"):
			inProgress++
		case strings.Contains(lower, "block"):
			blocked++
		}
	}

	// Velocity score (0-25): completion rate
	completionRate := float64(done) / float64(total)
	velocityScore := int(completionRate * 25)

	// Blocker score (0-25): fewer blockers = higher score
	blockerRatio := float64(blocked) / float64(total)
	blockerScore := 25 - int(blockerRatio*25*4) // amplify penalty
	if blockerScore < 0 {
		blockerScore = 0
	}

	// Scope score (0-25): based on historical carryover
	scopeScore := 20 // default decent
	prevSnapshot, _ := h.Memory.GetLatestSnapshot(ctx, boardID)
	if prevSnapshot != nil && prevSnapshot.TotalIssues > 0 {
		scopeChange := float64(total-prevSnapshot.TotalIssues) / float64(prevSnapshot.TotalIssues)
		if scopeChange > 0.2 {
			scopeScore = 10 // scope creep detected
		} else if scopeChange <= 0 {
			scopeScore = 25 // no creep
		}
	}

	// Team score (0-25): even distribution, no one overloaded
	assigneeCounts := map[string]int{}
	for _, issue := range issues {
		if issue.Assignee != "" {
			assigneeCounts[issue.Assignee]++
		}
	}
	teamScore := 20
	if len(assigneeCounts) > 0 {
		avg := total / len(assigneeCounts)
		maxLoad := 0
		for _, count := range assigneeCounts {
			if count > maxLoad {
				maxLoad = count
			}
		}
		if avg > 0 && maxLoad > avg*2 {
			teamScore = 10 // someone is overloaded
		} else {
			teamScore = 25
		}
	}

	overall := velocityScore + blockerScore + scopeScore + teamScore

	details := map[string]any{
		"total_issues": total,
		"done":         done,
		"in_progress":  inProgress,
		"blocked":      blocked,
		"completion":   fmt.Sprintf("%.0f%%", completionRate*100),
		"assignees":    len(assigneeCounts),
	}
	detailsJSON, _ := json.Marshal(details)

	// Save score to memory
	score := &memdom.HealthScore{
		SprintName:    sprint.Name,
		BoardID:       boardID,
		ComputedAt:    time.Now(),
		OverallScore:  overall,
		VelocityScore: velocityScore,
		BlockerScore:  blockerScore,
		ScopeScore:    scopeScore,
		TeamScore:     teamScore,
		Details:       string(detailsJSON),
	}
	_ = h.Memory.SaveHealthScore(ctx, score)

	// Determine status
	status := "HEALTHY"
	if overall < 50 {
		status = "AT RISK"
	} else if overall < 70 {
		status = "WATCH"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprint Health: %s (%d/100)\n", status, overall))
	sb.WriteString(fmt.Sprintf("Sprint: %s\n\n", sprint.Name))
	sb.WriteString(fmt.Sprintf("Velocity:  %d/25 (completion: %.0f%%)\n", velocityScore, completionRate*100))
	sb.WriteString(fmt.Sprintf("Blockers:  %d/25 (blocked: %d/%d)\n", blockerScore, blocked, total))
	sb.WriteString(fmt.Sprintf("Scope:     %d/25 %s\n", scopeScore, scopeNote(scopeScore)))
	sb.WriteString(fmt.Sprintf("Team:      %d/25 %s\n", teamScore, teamNote(teamScore)))
	sb.WriteString(fmt.Sprintf("\nTotal: %d | Done: %d | In Progress: %d | Blocked: %d\n", total, done, inProgress, blocked))

	return textResult(sb.String()), nil
}

func scopeNote(score int) string {
	if score <= 10 {
		return "(scope creep detected)"
	}
	return "(stable scope)"
}

func teamNote(score int) string {
	if score <= 10 {
		return "(workload imbalance)"
	}
	return "(well distributed)"
}

// RecordDependency records a dependency between issues.
func (h *Handlers) RecordDependency(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	from, err := req.RequireString("from_issue")
	if err != nil {
		return errorResult("from_issue required (issue that is blocked)"), nil
	}
	to, err := req.RequireString("to_issue")
	if err != nil {
		return errorResult("to_issue required (what it depends on)"), nil
	}

	d := &memdom.Dependency{
		FromIssueKey:   from,
		ToIssueKey:     to,
		DependencyType: req.GetString("type", "blocks"),
		Description:    req.GetString("description", ""),
		Status:         "open",
		CreatedAt:      time.Now(),
	}

	if err := h.Memory.SaveDependency(ctx, d); err != nil {
		return errorResult("Failed to save dependency: " + err.Error()), nil
	}

	return textResult(fmt.Sprintf("Dependency recorded: %s -[%s]-> %s\n%s",
		from, d.DependencyType, to, d.Description)), nil
}

// ResolveDependency marks a dependency as resolved.
func (h *Handlers) ResolveDependency(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := req.RequireInt("dependency_id")
	if err != nil {
		return errorResult("dependency_id required"), nil
	}

	if err := h.Memory.ResolveDependency(ctx, int64(id)); err != nil {
		return errorResult("Failed to resolve: " + err.Error()), nil
	}

	return textResult(fmt.Sprintf("Dependency #%d resolved.", id)), nil
}

// GetDependencies shows dependency map.
func (h *Handlers) GetDependencies(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	issueKey := req.GetString("issue_key", "")

	var deps []memdom.Dependency
	var err error

	if issueKey != "" {
		deps, err = h.Memory.GetDependenciesForIssue(ctx, issueKey)
	} else {
		deps, err = h.Memory.GetOpenDependencies(ctx)
	}
	if err != nil {
		return errorResult("Failed to get dependencies: " + err.Error()), nil
	}

	if len(deps) == 0 {
		return textResult("No open dependencies."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Dependencies (%d):\n\n", len(deps)))
	for _, d := range deps {
		days := int(time.Since(d.CreatedAt).Hours() / 24)
		status := "OPEN"
		if d.ResolvedAt != nil {
			status = "RESOLVED"
		}
		sb.WriteString(fmt.Sprintf("#%d [%s] %s -[%s]-> %s (%d days)\n",
			d.ID, status, d.FromIssueKey, d.DependencyType, d.ToIssueKey, days))
		if d.Description != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", d.Description))
		}
	}

	return textResult(sb.String()), nil
}

// RecordMeeting records meeting notes.
func (h *Handlers) RecordMeeting(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	meetingType, err := req.RequireString("meeting_type")
	if err != nil {
		return errorResult("meeting_type required (standup, planning, retro, grooming, adhoc)"), nil
	}

	m := &memdom.MeetingNote{
		MeetingType: meetingType,
		Date:        time.Now(),
		Attendees:   req.GetString("attendees", ""),
		Notes:       req.GetString("notes", ""),
		Decisions:   req.GetString("decisions", ""),
		ActionItems: req.GetString("action_items", ""),
		SprintName:  req.GetString("sprint_name", ""),
	}

	if err := h.Memory.SaveMeetingNote(ctx, m); err != nil {
		return errorResult("Failed to save meeting: " + err.Error()), nil
	}

	return textResult(fmt.Sprintf("Meeting recorded: %s (%s)\nDecisions: %s\nActions: %s",
		meetingType, m.Date.Format("2006-01-02"), m.Decisions, m.ActionItems)), nil
}

// GetMeetings retrieves meeting notes history.
func (h *Handlers) GetMeetings(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	meetingType := req.GetString("meeting_type", "")
	limit := req.GetInt("limit", 10)

	notes, err := h.Memory.GetMeetingNotes(ctx, meetingType, limit)
	if err != nil {
		return errorResult("Failed to get meetings: " + err.Error()), nil
	}

	if len(notes) == 0 {
		return textResult("No meeting notes recorded yet."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Meeting Notes (%d):\n\n", len(notes)))
	for _, m := range notes {
		sb.WriteString(fmt.Sprintf("[%s] %s — Sprint: %s\n", m.Date.Format("2006-01-02"), m.MeetingType, m.SprintName))
		if m.Notes != "" {
			sb.WriteString(fmt.Sprintf("  Notes: %s\n", m.Notes))
		}
		if m.Decisions != "" {
			sb.WriteString(fmt.Sprintf("  Decisions: %s\n", m.Decisions))
		}
		if m.ActionItems != "" {
			sb.WriteString(fmt.Sprintf("  Actions: %s\n", m.ActionItems))
		}
		sb.WriteString("\n")
	}

	return textResult(sb.String()), nil
}

// CapacityPlan generates capacity planning based on velocity history and team size.
func (h *Handlers) CapacityPlan(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	teamSize := req.GetInt("team_size", 0)
	sprintDays := req.GetInt("sprint_days", 10)
	plannedLeave := req.GetInt("planned_leave_days", 0)

	// Get velocity history
	snaps, err := h.Memory.GetSprintSnapshots(ctx, boardID, 6)
	if err != nil || len(snaps) == 0 {
		return textResult("Not enough sprint history for capacity planning. Capture at least 3 sprint snapshots first."), nil
	}

	var totalVelocity int
	for _, s := range snaps {
		totalVelocity += s.Velocity
	}
	avgVelocity := totalVelocity / len(snaps)

	// Calculate available capacity
	var sb strings.Builder
	sb.WriteString("Capacity Planning\n\n")
	sb.WriteString(fmt.Sprintf("Historical average velocity: %d points/sprint (%d sprints)\n", avgVelocity, len(snaps)))
	sb.WriteString(fmt.Sprintf("Sprint duration: %d days\n", sprintDays))

	if teamSize > 0 {
		totalDays := teamSize * sprintDays
		availableDays := totalDays - plannedLeave
		capacityRatio := float64(availableDays) / float64(totalDays)
		recommendedCapacity := int(float64(avgVelocity) * capacityRatio)

		sb.WriteString(fmt.Sprintf("\nTeam size: %d\n", teamSize))
		sb.WriteString(fmt.Sprintf("Planned leave: %d days\n", plannedLeave))
		sb.WriteString(fmt.Sprintf("Available capacity: %d/%d days (%.0f%%)\n", availableDays, totalDays, capacityRatio*100))
		sb.WriteString(fmt.Sprintf("\nRecommended commitment: %d story points\n", recommendedCapacity))
		sb.WriteString(fmt.Sprintf("Conservative (80%%): %d points\n", int(float64(recommendedCapacity)*0.8)))
		sb.WriteString(fmt.Sprintf("Stretch (120%%): %d points\n", int(float64(recommendedCapacity)*1.2)))
	} else {
		sb.WriteString(fmt.Sprintf("\nRecommended commitment: %d story points (avg velocity)\n", avgVelocity))
		sb.WriteString(fmt.Sprintf("Conservative (80%%): %d points\n", int(float64(avgVelocity)*0.8)))
		sb.WriteString("(Provide team_size + planned_leave_days for more precise planning)\n")
	}

	// Trend-based adjustment
	if len(snaps) >= 3 {
		recent3 := (snaps[0].Velocity + snaps[1].Velocity + snaps[2].Velocity) / 3
		if recent3 > avgVelocity {
			sb.WriteString(fmt.Sprintf("\nTrend: IMPROVING (recent 3-sprint avg: %d vs overall: %d)\n", recent3, avgVelocity))
		} else if recent3 < avgVelocity {
			sb.WriteString(fmt.Sprintf("\nTrend: DECLINING (recent 3-sprint avg: %d vs overall: %d) — investigate\n", recent3, avgVelocity))
		}
	}

	return textResult(sb.String()), nil
}

// AutoDetectRisks proactively scans Jira for risk signals and records them.
func (h *Handlers) AutoDetectRisks(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	var findings []string

	// 1. Check for stale tickets (>7 days no update in active sprint)
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) > 0 {
		issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
		staleCount := 0
		for _, issue := range issues {
			lower := strings.ToLower(issue.Status)
			if !strings.Contains(lower, "done") && !strings.Contains(lower, "closed") {
				daysSinceUpdate := int(time.Since(issue.Updated).Hours() / 24)
				if daysSinceUpdate >= 7 {
					staleCount++
				}
			}
		}
		if staleCount > 0 {
			findings = append(findings, fmt.Sprintf("STALE: %d issues in active sprint with no update in 7+ days", staleCount))
			_ = h.Memory.SaveRisk(ctx, &memdom.Risk{
				Title:        fmt.Sprintf("%d stale issues in sprint (auto-detected)", staleCount),
				Severity:     "medium",
				Status:       "open",
				IdentifiedAt: time.Now(),
				SprintName:   sprints[0].Name,
				Mitigation:   "Review stale issues in standup, reassign or remove from sprint",
			})
		}

		// 2. Check workload imbalance
		assigneeCounts := map[string]int{}
		for _, issue := range issues {
			if issue.Assignee != "" {
				assigneeCounts[issue.Assignee]++
			}
		}
		if len(assigneeCounts) > 1 {
			var maxPerson string
			var maxCount int
			totalAssigned := 0
			for person, count := range assigneeCounts {
				totalAssigned += count
				if count > maxCount {
					maxCount = count
					maxPerson = person
				}
			}
			avg := totalAssigned / len(assigneeCounts)
			if maxCount > avg*2 && maxCount > 5 {
				findings = append(findings, fmt.Sprintf("OVERLOAD: %s has %d issues (avg: %d)", maxPerson, maxCount, avg))
				_ = h.Memory.SaveRisk(ctx, &memdom.Risk{
					Title:        fmt.Sprintf("Workload imbalance: %s has %dx average load (auto-detected)", maxPerson, maxCount/avg),
					Severity:     "high",
					Status:       "open",
					Owner:        maxPerson,
					IdentifiedAt: time.Now(),
					SprintName:   sprints[0].Name,
					Mitigation:   "Redistribute work or negotiate scope reduction",
				})
			}
		}

		// 3. Check blocked count
		blockedCount := 0
		for _, issue := range issues {
			if strings.Contains(strings.ToLower(issue.Status), "block") {
				blockedCount++
			}
		}
		if blockedCount >= 3 {
			findings = append(findings, fmt.Sprintf("BLOCKED: %d issues currently blocked", blockedCount))
			_ = h.Memory.SaveRisk(ctx, &memdom.Risk{
				Title:        fmt.Sprintf("%d blocked issues (auto-detected)", blockedCount),
				Severity:     "high",
				Status:       "open",
				IdentifiedAt: time.Now(),
				SprintName:   sprints[0].Name,
				Mitigation:   "Escalate blockers immediately, consider sprint goal revision",
			})
		}
	}

	// 4. Check for long-standing blockers in memory
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	for _, b := range blockers {
		days := int(time.Since(b.BlockedSince).Hours() / 24)
		if days > 5 {
			findings = append(findings, fmt.Sprintf("CHRONIC BLOCKER: #%d stuck for %d days — %s", b.ID, days, b.Description))
		}
	}

	// 5. Check pending action items overdue
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	overdueActions := 0
	for _, a := range actions {
		if a.DueDate != nil && a.DueDate.Before(time.Now()) {
			overdueActions++
		}
	}
	if overdueActions > 0 {
		findings = append(findings, fmt.Sprintf("OVERDUE ACTIONS: %d retro action items past due date", overdueActions))
	}

	if len(findings) == 0 {
		return textResult("Auto-scan complete: No new risks detected. Team is looking healthy!"), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Auto Risk Detection: %d findings\n\n", len(findings)))
	for i, f := range findings {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, f))
	}
	sb.WriteString("\n(Detected risks have been auto-recorded to the risk register)")

	return textResult(sb.String()), nil
}

// HealthHistory shows health score trends over time.
func (h *Handlers) HealthHistory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	scores, err := h.Memory.GetHealthScores(ctx, boardID, 10)
	if err != nil {
		return errorResult("Failed to get health scores: " + err.Error()), nil
	}

	if len(scores) == 0 {
		return textResult("No health scores recorded yet. Use pm_sprint_health to compute one."), nil
	}

	var sb strings.Builder
	sb.WriteString("Health Score History:\n\n")
	sb.WriteString("Date       | Sprint | Overall | Velocity | Blockers | Scope | Team\n")
	sb.WriteString("-----------|--------|---------|----------|----------|-------|-----\n")
	for _, s := range scores {
		sb.WriteString(fmt.Sprintf("%s | %s | %d | %d | %d | %d | %d\n",
			s.ComputedAt.Format("2006-01-02"), s.SprintName, s.OverallScore,
			s.VelocityScore, s.BlockerScore, s.ScopeScore, s.TeamScore))
	}

	return textResult(sb.String()), nil
}
