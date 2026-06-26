package tools

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// TrackDailyProgress captures daily sprint state for burndown tracking.
func (h *Handlers) TrackDailyProgress(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil {
		return sanitizedError("failed to get sprints", err), nil
	}
	if len(sprints) == 0 {
		return textResult("No active sprint found."), nil
	}

	sprint := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprint.ID)
	if err != nil {
		return sanitizedError("failed to get sprint issues", err), nil
	}

	var done, inProgress, todo, blocked int
	for _, issue := range issues {
		lower := strings.ToLower(issue.Status)
		switch {
		case strings.Contains(lower, "done") || strings.Contains(lower, "closed") || strings.Contains(lower, "resolved"):
			done++
		case strings.Contains(lower, "progress") || strings.Contains(lower, "review"):
			inProgress++
		case strings.Contains(lower, "block"):
			blocked++
		default:
			todo++
		}
	}

	p := &memdom.DailyProgress{
		SprintName:  sprint.Name,
		BoardID:     boardID,
		Date:        time.Now(),
		TotalIssues: len(issues),
		Done:        done,
		InProgress:  inProgress,
		Todo:        todo,
		Blocked:     blocked,
		PointsDone:  req.GetInt("points_done", 0),
		PointsTotal: req.GetInt("points_total", 0),
	}

	if err := h.Memory.SaveDailyProgress(ctx, p); err != nil {
		return sanitizedError("failed to save daily progress", err), nil
	}

	return textResult(fmt.Sprintf("Daily progress recorded for %s:\nDone: %d | In Progress: %d | Todo: %d | Blocked: %d | Total: %d",
		sprint.Name, done, inProgress, todo, blocked, len(issues))), nil
}

// GetBurndown shows daily progress over time for a sprint.
func (h *Handlers) GetBurndown(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprintName := req.GetString("sprint_name", "")
	if sprintName == "" {
		sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
		if err != nil || len(sprints) == 0 {
			return textResult("No active sprint found. Provide sprint_name explicitly."), nil
		}
		sprintName = sprints[0].Name
	}

	progress, err := h.Memory.GetDailyProgress(ctx, boardID, sprintName)
	if err != nil {
		return sanitizedError("failed to get daily progress", err), nil
	}
	if len(progress) == 0 {
		return textResult("No daily progress recorded yet. Use pm_track_daily to capture data."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Burndown: %s\n\n", sprintName))
	sb.WriteString("Date       | Done | InProg | Todo | Blocked | Total\n")
	sb.WriteString("-----------|------|--------|------|---------|------\n")
	for _, p := range progress {
		sb.WriteString(fmt.Sprintf("%s | %d | %d | %d | %d | %d\n",
			p.Date.Format("2006-01-02"), p.Done, p.InProgress, p.Todo, p.Blocked, p.TotalIssues))
	}

	return textResult(sb.String()), nil
}

// SetSprintGoal sets an explicit sprint goal with key results.
func (h *Handlers) SetSprintGoal(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	goal, err := req.RequireString("goal")
	if err != nil {
		return errorResult("goal required"), nil
	}

	sprintName := req.GetString("sprint_name", "")
	if sprintName == "" {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			sprintName = sprints[0].Name
		}
	}

	g := &memdom.SprintGoal{
		SprintName: sprintName,
		BoardID:    boardID,
		Goal:       goal,
		KeyResults: req.GetString("key_results", ""),
		Status:     "active",
		CreatedAt:  time.Now(),
	}

	if err := h.Memory.SaveSprintGoal(ctx, g); err != nil {
		return sanitizedError("failed to save sprint goal", err), nil
	}

	return textResult(fmt.Sprintf("Sprint goal set for %s:\n%s\nKey Results: %s",
		sprintName, goal, g.KeyResults)), nil
}

// CloseSprintGoal closes a sprint goal with outcome.
func (h *Handlers) CloseSprintGoal(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	goalID, err := req.RequireInt("goal_id")
	if err != nil {
		return errorResult("goal_id required"), nil
	}
	status, err := req.RequireString("status")
	if err != nil {
		return errorResult("status required (achieved, partially_achieved, missed)"), nil
	}

	now := time.Now()
	g := &memdom.SprintGoal{
		ID:       int64(goalID),
		Status:   status,
		Outcome:  req.GetString("outcome", ""),
		ClosedAt: &now,
	}

	if err := h.Memory.UpdateSprintGoal(ctx, g); err != nil {
		return sanitizedError("failed to close sprint goal", err), nil
	}

	return textResult(fmt.Sprintf("Goal #%d closed as: %s", goalID, status)), nil
}

// GetSprintGoals shows active or historical sprint goals.
func (h *Handlers) GetSprintGoals(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	showHistory := req.GetBool("show_history", false)

	var goals []memdom.SprintGoal
	if showHistory {
		goals, err = h.Memory.GetGoalHistory(ctx, boardID, 10)
	} else {
		goals, err = h.Memory.GetActiveGoals(ctx, boardID)
	}
	if err != nil {
		return sanitizedError("failed to get sprint goals", err), nil
	}
	if len(goals) == 0 {
		return textResult("No goals found."), nil
	}

	var sb strings.Builder

	// BLUF: Report goal count and active status first
	sb.WriteString(fmt.Sprintf("Found %d sprint goal(s)\n\n", len(goals)))
	for _, g := range goals {
		sb.WriteString(fmt.Sprintf("#%d [%s] %s\n  Goal: %s\n", g.ID, g.Status, g.SprintName, g.Goal))
		if g.KeyResults != "" {
			sb.WriteString(fmt.Sprintf("  Key Results: %s\n", g.KeyResults))
		}
		if g.Outcome != "" {
			sb.WriteString(fmt.Sprintf("  Outcome: %s\n", g.Outcome))
		}
		sb.WriteString("\n")
	}

	return textResult(sb.String()), nil
}

// ManageDoD handles Definition of Done list/add/remove.
func (h *Handlers) ManageDoD(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	action := req.GetString("action", "list")
	project := req.GetString("project", "*")

	switch action {
	case "add":
		item := req.GetString("item", "")
		if item == "" {
			return errorResult("item required for add action"), nil
		}
		d := &memdom.DoDItem{
			Project:  project,
			Item:     item,
			Category: req.GetString("category", "general"),
			OrderNum: req.GetInt("order", 0),
			Active:   true,
		}
		if err := h.Memory.SaveDoDItem(ctx, d); err != nil {
			return sanitizedError("failed to save DoD item", err), nil
		}
		return textResult(fmt.Sprintf("DoD item added: [%s] %s", d.Category, item)), nil

	case "remove":
		id := req.GetInt("item_id", 0)
		if id == 0 {
			return errorResult("item_id required for remove action"), nil
		}
		if err := h.Memory.DeleteDoDItem(ctx, int64(id)); err != nil {
			return sanitizedError("failed to remove DoD item", err), nil
		}
		return textResult(fmt.Sprintf("DoD item #%d removed.", id)), nil

	default: // list
		items, err := h.Memory.GetDoD(ctx, project)
		if err != nil {
			return sanitizedError("failed to get DoD items", err), nil
		}
		if len(items) == 0 {
			return textResult("No Definition of Done items. Use action=add to create one."), nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Definition of Done (project: %s):\n\n", project))
		for _, item := range items {
			sb.WriteString(fmt.Sprintf("  #%d [%s] %s\n", item.ID, item.Category, item.Item))
		}
		return textResult(sb.String()), nil
	}
}

// PMDashboard shows everything in one view.
func (h *Handlers) PMDashboard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	var sb strings.Builder

	// Compute overall status signal first (BLUF)
	overallSignal := "On Track"
	if boardID > 0 {
		scores, _ := h.Memory.GetHealthScores(ctx, boardID, 1)
		if len(scores) > 0 && scores[0].OverallScore < 70 {
			overallSignal = "⚠ Watch"
		}
		if len(scores) > 0 && scores[0].OverallScore < 50 {
			overallSignal = "🔴 At Risk"
		}
		blockers, _ := h.Memory.GetActiveBlockers(ctx)
		if len(blockers) > 3 {
			overallSignal = "🔴 At Risk"
		}
	}
	sb.WriteString(fmt.Sprintf("Overall: %s\n\n", overallSignal))

	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			sprint := sprints[0]
			issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)
			var done, inProg, blocked int
			for _, issue := range issues {
				lower := strings.ToLower(issue.Status)
				switch {
				case strings.Contains(lower, "done") || strings.Contains(lower, "closed"):
					done++
				case strings.Contains(lower, "progress") || strings.Contains(lower, "review"):
					inProg++
				case strings.Contains(lower, "block"):
					blocked++
				}
			}
			completion := 0.0
			if len(issues) > 0 {
				completion = float64(done) / float64(len(issues)) * 100
			}
			sb.WriteString(fmt.Sprintf("Sprint: %s (Goal: %s)\n", sprint.Name, sprint.Goal))
			sb.WriteString(fmt.Sprintf("Progress: %d/%d (%.0f%%) | In Progress: %d | Blocked: %d\n\n", done, len(issues), completion, inProg, blocked))
		}

		scores, _ := h.Memory.GetHealthScores(ctx, boardID, 1)
		if len(scores) > 0 {
			s := scores[0]
			status := "HEALTHY"
			if s.OverallScore < 50 {
				status = "AT RISK"
			} else if s.OverallScore < 70 {
				status = "WATCH"
			}
			sb.WriteString(fmt.Sprintf("Health: %d/100 (%s)\n\n", s.OverallScore, status))
		}
	}

	risks, _ := h.Memory.GetOpenRisks(ctx)
	if len(risks) > 0 {
		sb.WriteString(fmt.Sprintf("RISKS: %d open\n", len(risks)))
		for _, r := range risks {
			if r.Severity == "critical" || r.Severity == "high" {
				sb.WriteString(fmt.Sprintf("  [%s] %s\n", strings.ToUpper(r.Severity), r.Title))
			}
		}
		sb.WriteString("\n")
	}

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		sb.WriteString(fmt.Sprintf("BLOCKERS: %d active\n", len(blockers)))
		for _, b := range blockers {
			days := int(time.Since(b.BlockedSince).Hours() / 24)
			sb.WriteString(fmt.Sprintf("  [%d days] %s\n", days, b.Description))
		}
		sb.WriteString("\n")
	}

	deps, _ := h.Memory.GetOpenDependencies(ctx)
	if len(deps) > 0 {
		sb.WriteString(fmt.Sprintf("DEPENDENCIES: %d open\n", len(deps)))
		for _, d := range deps {
			sb.WriteString(fmt.Sprintf("  %s -> %s\n", d.FromIssueKey, d.ToIssueKey))
		}
		sb.WriteString("\n")
	}

	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 0 {
		sb.WriteString(fmt.Sprintf("PENDING ACTIONS: %d\n", len(actions)))
		for _, a := range actions {
			sb.WriteString(fmt.Sprintf("  - %s (%s)\n", a.Description, a.Owner))
		}
		sb.WriteString("\n")
	}

	return textResult(sb.String()), nil
}

// GenerateReleaseNotes creates release notes from completed sprint issues.
func (h *Handlers) GenerateReleaseNotes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return textResult("No active sprint found."), nil
	}

	sprint := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprint.ID)
	if err != nil {
		return sanitizedError("failed to get sprint issues", err), nil
	}

	var features, bugs, tasks []string
	for _, issue := range issues {
		lower := strings.ToLower(issue.Status)
		if !strings.Contains(lower, "done") && !strings.Contains(lower, "closed") && !strings.Contains(lower, "resolved") {
			continue
		}
		entry := fmt.Sprintf("- **%s** %s", issue.Key, issue.Summary)
		switch strings.ToLower(issue.Type) {
		case "bug":
			bugs = append(bugs, entry)
		case "story", "feature":
			features = append(features, entry)
		default:
			tasks = append(tasks, entry)
		}
	}

	var sb strings.Builder

	// BLUF: Start with sprint name and goal
	sb.WriteString(fmt.Sprintf("# Release Notes: %s\n\n", sprint.Name))
	if sprint.Goal != "" {
		sb.WriteString(fmt.Sprintf("**Goal:** %s\n\n", sprint.Goal))
	}

	// BLUF: List completed items by category
	if len(features) > 0 {
		sb.WriteString("## Features\n" + strings.Join(features, "\n") + "\n\n")
	}
	if len(bugs) > 0 {
		sb.WriteString("## Bug Fixes\n" + strings.Join(bugs, "\n") + "\n\n")
	}
	if len(tasks) > 0 {
		sb.WriteString("## Tasks\n" + strings.Join(tasks, "\n") + "\n\n")
	}

	total := len(features) + len(bugs) + len(tasks)
	sb.WriteString(fmt.Sprintf("---\n*%d items delivered*\n", total))

	if req.GetBool("send_to_lark", false) {
		if err := h.Lark.SendMarkdown(ctx, "Release: "+sprint.Name, sb.String()); err != nil {
			slog.Warn("Lark release notes send failed", "detail", err.Error())
			return textResult(sb.String() + "\n(Lark send failed — check server logs)"), nil
		}
		return textResult(sb.String() + "\n(Sent to Lark)"), nil
	}

	return textResult(sb.String()), nil
}

// Escalate checks for critical conditions and sends alerts to Lark.
func (h *Handlers) Escalate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	var escalated []string

	// Critical risks open > 3 days
	risks, _ := h.Memory.GetOpenRisks(ctx)
	for _, r := range risks {
		if (r.Severity == "critical" || r.Severity == "high") && time.Since(r.IdentifiedAt).Hours() > 72 {
			title := fmt.Sprintf("RISK: [%s] %s (open %d days)", r.Severity, r.Title, int(time.Since(r.IdentifiedAt).Hours()/24))
			_ = h.Lark.SendMarkdown(ctx, "Risk Escalation", title)
			_ = h.Memory.SaveEscalation(ctx, &memdom.Escalation{
				Type:        "risk",
				ReferenceID: r.ID,
				Title:       title,
				Severity:    r.Severity,
				EscalatedAt: time.Now(),
				Channel:     "lark",
			})
			escalated = append(escalated, title)
		}
	}

	// Blockers open > 3 days
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	for _, b := range blockers {
		if time.Since(b.BlockedSince).Hours() > 72 {
			title := fmt.Sprintf("BLOCKER: %s (blocked %d days)", b.Description, int(time.Since(b.BlockedSince).Hours()/24))
			_ = h.Lark.SendMarkdown(ctx, "Blocker Escalation", title)
			_ = h.Memory.SaveEscalation(ctx, &memdom.Escalation{
				Type:        "blocker",
				ReferenceID: b.ID,
				Title:       title,
				Severity:    "high",
				EscalatedAt: time.Now(),
				Channel:     "lark",
			})
			escalated = append(escalated, title)
		}
	}

	// Sprint health < 50
	scores, _ := h.Memory.GetHealthScores(ctx, boardID, 1)
	if len(scores) > 0 && scores[0].OverallScore < 50 {
		title := fmt.Sprintf("SPRINT AT RISK: Health score %d/100", scores[0].OverallScore)
		_ = h.Lark.SendMarkdown(ctx, "Sprint Risk", title)
		_ = h.Memory.SaveEscalation(ctx, &memdom.Escalation{
			Type:        "sprint_health",
			Title:       title,
			Severity:    "critical",
			EscalatedAt: time.Now(),
			Channel:     "lark",
		})
		escalated = append(escalated, title)
	}

	if len(escalated) == 0 {
		return textResult("No escalation needed. All clear."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Escalated %d items to Lark:\n\n", len(escalated)))
	for _, e := range escalated {
		sb.WriteString(fmt.Sprintf("- %s\n", e))
	}
	return textResult(sb.String()), nil
}
// GetEscalations shows recent escalations.
func (h *Handlers) GetEscalations(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	limit := req.GetInt("limit", 10)

	escalations, err := h.Memory.GetRecentEscalations(ctx, limit)
	if err != nil {
		return sanitizedError("failed to get escalations", err), nil
	}
	if len(escalations) == 0 {
		return textResult("No escalations recorded."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Recent Escalations (%d):\n\n", len(escalations)))
	for _, e := range escalations {
		ack := "pending"
		if e.Acknowledged {
			ack = "ack"
		}
		sb.WriteString(fmt.Sprintf("#%d [%s] [%s] %s (%s via %s)\n",
			e.ID, e.Severity, ack, e.Title, e.Type, e.Channel))
	}

	return textResult(sb.String()), nil
}
