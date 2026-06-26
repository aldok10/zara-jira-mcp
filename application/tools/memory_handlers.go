package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// SnapshotSprint captures current sprint state into memory.
func (h *Handlers) SnapshotSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		switch strings.ToLower(issue.Status) {
		case "done", "closed", "resolved":
			done++
		case "in progress", "in review", "in development":
			inProgress++
		default:
			if strings.Contains(strings.ToLower(issue.Status), "block") {
				blocked++
			} else {
				todo++
			}
		}
	}

	total := len(issues)
	var completionRate float64
	if total > 0 {
		completionRate = float64(done) / float64(total) * 100
	}

	velocity := req.GetInt("velocity", 0)
	carryover := req.GetInt("carryover", 0)
	notes := req.GetString("notes", "")

	snap := &memdom.SprintSnapshot{
		SprintName:     sprint.Name,
		BoardID:        boardID,
		SnapshotDate:   time.Now(),
		TotalIssues:    total,
		Done:           done,
		InProgress:     inProgress,
		Todo:           todo,
		Blocked:        blocked,
		Carryover:      carryover,
		Velocity:       velocity,
		CompletionRate: completionRate,
		Notes:          notes,
	}

	if err := h.Memory.SaveSprintSnapshot(ctx, snap); err != nil {
		return sanitizedError("failed to save sprint snapshot", err), nil
	}

	return textResult(fmt.Sprintf("Sprint snapshot saved: %s\nDone: %d | In Progress: %d | Todo: %d | Blocked: %d\nCompletion: %.0f%%",
		sprint.Name, done, inProgress, todo, blocked, completionRate)), nil
}

// RecordRisk records a new project risk.
func (h *Handlers) RecordRisk(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}

	r := &memdom.Risk{
		Title:        title,
		Description:  req.GetString("description", ""),
		Severity:     req.GetString("severity", "medium"),
		Status:       "open",
		Owner:        req.GetString("owner", ""),
		Mitigation:   req.GetString("mitigation", ""),
		IdentifiedAt: time.Now(),
		SprintName:   req.GetString("sprint_name", ""),
	}

	if err := h.Memory.SaveRisk(ctx, r); err != nil {
		return sanitizedError("failed to save risk", err), nil
	}

	return textResult(fmt.Sprintf("Risk recorded: [%s] %s\nOwner: %s | Mitigation: %s",
		r.Severity, r.Title, r.Owner, r.Mitigation)), nil
}

// UpdateRisk updates an existing risk status.
func (h *Handlers) UpdateRisk(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := req.RequireInt("risk_id")
	if err != nil {
		return errorResult("risk_id required"), nil
	}

	status := req.GetString("status", "")
	if status == "" {
		return errorResult("status required (open, mitigating, resolved, accepted)"), nil
	}

	r := &memdom.Risk{ID: int64(id), Status: status}
	if status == "resolved" {
		now := time.Now()
		r.ResolvedAt = &now
	}

	r.Mitigation = req.GetString("mitigation", "")
	r.Owner = req.GetString("owner", "")
	r.Severity = req.GetString("severity", "")
	r.Title = req.GetString("title", "")
	r.Description = req.GetString("description", "")

	if err := h.Memory.UpdateRisk(ctx, r); err != nil {
		return sanitizedError("failed to update risk", err), nil
	}

	return textResult(fmt.Sprintf("Risk #%d updated to status: %s", id, status)), nil
}

// GetRiskDashboard shows all open risks sorted by severity.
func (h *Handlers) GetRiskDashboard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	risks, err := h.Memory.GetOpenRisks(ctx)
	if err != nil {
		return sanitizedError("failed to get risks", err), nil
	}

	if len(risks) == 0 {
		return textResult("No open risks. Looking good!"), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Open Risks (%d):\n\n", len(risks)))
	for _, r := range risks {
		days := int(time.Since(r.IdentifiedAt).Hours() / 24)
		sb.WriteString(fmt.Sprintf("[%s] #%d %s\n  Owner: %s | Days open: %d | Sprint: %s\n  Mitigation: %s\n\n",
			strings.ToUpper(r.Severity), r.ID, r.Title, r.Owner, days, r.SprintName, r.Mitigation))
	}

	return textResult(sb.String()), nil
}

// RecordDecision saves a project decision to memory.
func (h *Handlers) RecordDecision(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}
	decision, err := req.RequireString("decision")
	if err != nil {
		return errorResult("decision required (what was decided)"), nil
	}

	d := &memdom.Decision{
		Title:     title,
		Context:   req.GetString("context", ""),
		Decision:  decision,
		Rationale: req.GetString("rationale", ""),
		MadeBy:    req.GetString("made_by", ""),
		MadeAt:    time.Now(),
		Tags:      req.GetString("tags", ""),
	}

	if err := h.Memory.SaveDecision(ctx, d); err != nil {
		return sanitizedError("failed to save decision", err), nil
	}

	return textResult(fmt.Sprintf("Decision recorded: %s\nDecision: %s\nRationale: %s",
		d.Title, d.Decision, d.Rationale)), nil
}

// SearchDecisions searches decision log.
func (h *Handlers) SearchDecisions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := req.GetString("query", "")
	limit := req.GetInt("limit", 10)

	var decisions []memdom.Decision
	var err error

	if query != "" {
		decisions, err = h.Memory.SearchDecisions(ctx, query)
	} else {
		decisions, err = h.Memory.GetDecisions(ctx, limit)
	}
	if err != nil {
		return sanitizedError("failed to get decisions", err), nil
	}

	if len(decisions) == 0 {
		return textResult("No decisions found."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Decisions (%d):\n\n", len(decisions)))
	for _, d := range decisions {
		sb.WriteString(fmt.Sprintf("#%d [%s] %s\n  Decision: %s\n  Rationale: %s\n  By: %s | Tags: %s\n\n",
			d.ID, d.MadeAt.Format("2006-01-02"), d.Title, d.Decision, d.Rationale, d.MadeBy, d.Tags))
	}

	return textResult(sb.String()), nil
}

// RecordBlocker records a new blocker.
func (h *Handlers) RecordBlocker(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	desc, err := req.RequireString("description")
	if err != nil {
		return errorResult("description required"), nil
	}

	b := &memdom.Blocker{
		IssueKey:     req.GetString("issue_key", ""),
		Description:  desc,
		BlockedSince: time.Now(),
		Owner:        req.GetString("owner", ""),
	}

	if err := h.Memory.SaveBlocker(ctx, b); err != nil {
		return sanitizedError("failed to save blocker", err), nil
	}

	return textResult(fmt.Sprintf("Blocker recorded: %s\nIssue: %s | Owner: %s",
		desc, b.IssueKey, b.Owner)), nil
}

// ResolveBlocker marks a blocker as resolved.
func (h *Handlers) ResolveBlocker(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id, err := req.RequireInt("blocker_id")
	if err != nil {
		return errorResult("blocker_id required"), nil
	}
	resolution, err := req.RequireString("resolution")
	if err != nil {
		return errorResult("resolution required (how was it resolved)"), nil
	}

	if err := h.Memory.ResolveBlocker(ctx, int64(id), resolution); err != nil {
		return sanitizedError("failed to resolve blocker", err), nil
	}

	return textResult(fmt.Sprintf("Blocker #%d resolved: %s", id, resolution)), nil
}

// GetBlockers shows active blockers.
func (h *Handlers) GetBlockers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	showHistory := req.GetBool("show_history", false)

	var blockers []memdom.Blocker
	var err error

	if showHistory {
		blockers, err = h.Memory.GetBlockerHistory(ctx, 20)
	} else {
		blockers, err = h.Memory.GetActiveBlockers(ctx)
	}
	if err != nil {
		return sanitizedError("failed to get blockers", err), nil
	}

	if len(blockers) == 0 {
		return textResult("No active blockers. Ship it!"), nil
	}

	var sb strings.Builder
	if showHistory {
		sb.WriteString("Blocker History (last 20):\n\n")
	} else {
		sb.WriteString(fmt.Sprintf("Active Blockers (%d):\n\n", len(blockers)))
	}

	for _, b := range blockers {
		days := int(time.Since(b.BlockedSince).Hours() / 24)
		status := "ACTIVE"
		if b.ResolvedAt != nil {
			status = "RESOLVED"
			days = b.DaysBlocked
		}
		sb.WriteString(fmt.Sprintf("#%d [%s] %s\n  Issue: %s | Owner: %s | Days: %d\n",
			b.ID, status, b.Description, b.IssueKey, b.Owner, days))
		if b.Resolution != "" {
			sb.WriteString(fmt.Sprintf("  Resolution: %s\n", b.Resolution))
		}
		sb.WriteString("\n")
	}

	return textResult(sb.String()), nil
}

// RecordTeamMetric records a team member's sprint metrics.
func (h *Handlers) RecordTeamMetric(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	member, err := req.RequireString("member_name")
	if err != nil {
		return errorResult("member_name required"), nil
	}
	sprintName, err := req.RequireString("sprint_name")
	if err != nil {
		return errorResult("sprint_name required"), nil
	}

	m := &memdom.TeamMetric{
		MemberName:     member,
		SprintName:     sprintName,
		RecordedAt:     time.Now(),
		IssuesAssigned: req.GetInt("issues_assigned", 0),
		IssuesDone:     req.GetInt("issues_done", 0),
		BlockerCount:   req.GetInt("blocker_count", 0),
		CarryoverCount: req.GetInt("carryover_count", 0),
		Notes:          req.GetString("notes", ""),
	}

	if err := h.Memory.SaveTeamMetric(ctx, m); err != nil {
		return sanitizedError("failed to save metric", err), nil
	}

	return textResult(fmt.Sprintf("Metric saved for %s (Sprint: %s)\nAssigned: %d | Done: %d | Blockers: %d | Carryover: %d",
		member, sprintName, m.IssuesAssigned, m.IssuesDone, m.BlockerCount, m.CarryoverCount)), nil
}

// GetTeamHealth shows team workload overview.
func (h *Handlers) GetTeamHealth(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintName := req.GetString("sprint_name", "")
	memberName := req.GetString("member_name", "")

	var metrics []memdom.TeamMetric
	var err error

	if memberName != "" {
		metrics, err = h.Memory.GetTeamMetrics(ctx, memberName, 10)
	} else if sprintName != "" {
		metrics, err = h.Memory.GetTeamOverview(ctx, sprintName)
	} else {
		return errorResult("Provide either sprint_name (for team overview) or member_name (for individual history)"), nil
	}
	if err != nil {
		return sanitizedError("failed to get metrics", err), nil
	}

	if len(metrics) == 0 {
		return textResult("No metrics recorded yet."), nil
	}

	var sb strings.Builder
	if memberName != "" {
		sb.WriteString(fmt.Sprintf("History for %s:\n\n", memberName))
	} else {
		sb.WriteString(fmt.Sprintf("Team Overview - %s:\n\n", sprintName))
	}

	for _, m := range metrics {
		completionRate := 0.0
		if m.IssuesAssigned > 0 {
			completionRate = float64(m.IssuesDone) / float64(m.IssuesAssigned) * 100
		}
		sb.WriteString(fmt.Sprintf("%s | Sprint: %s | Assigned: %d | Done: %d (%.0f%%) | Blockers: %d | Carryover: %d\n",
			m.MemberName, m.SprintName, m.IssuesAssigned, m.IssuesDone, completionRate, m.BlockerCount, m.CarryoverCount))
		if m.Notes != "" {
			sb.WriteString(fmt.Sprintf("  Notes: %s\n", m.Notes))
		}
	}

	return textResult(sb.String()), nil
}

// RecordRetrospective saves a sprint retrospective.
func (h *Handlers) RecordRetrospective(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintName, err := req.RequireString("sprint_name")
	if err != nil {
		return errorResult("sprint_name required"), nil
	}

	r := &memdom.Retrospective{
		SprintName:   sprintName,
		Date:         time.Now(),
		WentWell:     req.GetString("went_well", ""),
		Improvements: req.GetString("improvements", ""),
		ActionItems:  req.GetString("action_items", ""),
		Status:       "open",
	}

	if err := h.Memory.SaveRetrospective(ctx, r); err != nil {
		return sanitizedError("failed to save retro", err), nil
	}

	return textResult(fmt.Sprintf("Retrospective saved for sprint: %s\nWent Well: %s\nImprovements: %s\nAction Items: %s",
		sprintName, r.WentWell, r.Improvements, r.ActionItems)), nil
}

// GetActionItems shows pending action items from retros.
func (h *Handlers) GetActionItems(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	items, err := h.Memory.GetPendingActionItems(ctx)
	if err != nil {
		return sanitizedError("failed to get action items", err), nil
	}

	if len(items) == 0 {
		return textResult("No pending action items. All caught up!"), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Pending Action Items (%d):\n\n", len(items)))
	for _, a := range items {
		due := "no due date"
		if a.DueDate != nil {
			due = a.DueDate.Format("2006-01-02")
		}
		sb.WriteString(fmt.Sprintf("#%d %s\n  Owner: %s | Due: %s\n\n", a.ID, a.Description, a.Owner, due))
	}

	return textResult(sb.String()), nil
}
