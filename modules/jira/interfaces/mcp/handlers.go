// Package mcp provides MCP tool handlers for the jira module.
package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/modules/jira/application/port"
	"github.com/aldok10/zara-jira-mcp/modules/jira/domain"
	mem "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/validate"
)

// Handlers holds dependencies for jira MCP tool handlers.
type Handlers struct {
	Jira   port.Inbound
	Memory mem.Store // optional: auto-records blockers/risks to PM memory
}

// NewHandlers creates a new jira MCP handlers instance.
func NewHandlers(jiraService port.Inbound, memStore mem.Store) *Handlers {
	return &Handlers{Jira: jiraService, Memory: memStore}
}

// Health returns server version and status.
func (h *Handlers) Health(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("zara-jira-mcp | status: ok | modular handlers"), nil
}

// recordBlockers scans a list of issues for blocked items and auto-records to PM memory.
// Skips if memory is not configured. Deduplicates against existing active blockers.
// Uses issue-level heuristic for blocking detection.
func (h *Handlers) recordBlockers(ctx context.Context, issues []domain.Issue) {
	if h.Memory == nil {
		return
	}

	for _, issue := range issues {
		if !issue.IsBlocked() {
			continue
		}

		// Check if already recorded (dedup)
		row := h.Memory.DB().QueryRow(
			"SELECT COUNT(*) FROM blockers WHERE issue_key = ? AND resolved_at IS NULL",
			issue.Key,
		)
		var count int
		if err := row.Scan(&count); err != nil || count > 0 {
			continue
		}

		blocker := &mem.Blocker{
			IssueKey:     issue.Key,
			Description:  fmt.Sprintf("Auto-detected: %s — %s", issue.Summary, issue.Status),
			BlockedSince: issue.Updated,
			Owner:        issue.Assignee,
			DaysBlocked:  0,
		}

		if err := h.Memory.SaveBlocker(ctx, blocker); err != nil {
			slog.Warn("auto-record blocker", "issue", issue.Key, "error", err)
		} else {
			slog.Info("auto-recorded blocker", "issue", issue.Key, "assignee", issue.Assignee)
		}
	}
}

// recordSnapshot auto-records a sprint snapshot to PM memory when sprint summary is called.
// Uses board-aware classification when board config is available.
func (h *Handlers) recordSnapshot(ctx context.Context, sprintName string, boardID int, issues []domain.Issue) {
	if h.Memory == nil {
		return
	}
	if sprintName == "" {
		return
	}

	// Build board-aware classifier — fetch config once, use for all issues
	classify := domain.HeuristicClassify // fallback
	if cfg, err := h.Jira.GetBoardConfiguration(ctx, boardID); err == nil && cfg != nil {
		classify = cfg.StatusCategory
	}

	total := len(issues)
	var done, inProg, blocked, todo int
	for _, i := range issues {
		switch classify(i.Status) {
		case "done":
			done++
		case "progress":
			inProg++
		case "blocked":
			blocked++
		default:
			todo++
		}
	}

	snap := &mem.SprintSnapshot{
		SprintName:   sprintName,
		BoardID:      boardID,
		SnapshotDate: time.Now(),
		TotalIssues:  total,
		Done:         done,
		InProgress:   inProg,
		Todo:         todo,
		Blocked:      blocked,
	}
	snap.CompletionRate = snap.CalculateCompletionRate()

	if err := h.Memory.SaveSprintSnapshot(ctx, snap); err != nil {
		slog.Warn("auto-record snapshot", "sprint", sprintName, "error", err)
	} else {
		slog.Info("auto-recorded sprint snapshot", "sprint", sprintName, "done/total", fmt.Sprintf("%d/%d", done, total))
	}
}

// recordBlockersForBoard scans issues for blocked items using board-aware classification.
// Uses the board configuration to detect blocking statuses (e.g. "Stalled", "On Hold").
func (h *Handlers) recordBlockersForBoard(ctx context.Context, boardID int, issues []domain.Issue) {
	if h.Memory == nil {
		return
	}

	// Fetch board config for board-aware blocking detection
	isBlocked := func(status string) bool {
		return h.classifyStatus(ctx, boardID, status) == "blocked"
	}

	for _, issue := range issues {
		if !isBlocked(issue.Status) {
			continue
		}

		// Check if already recorded (dedup)
		row := h.Memory.DB().QueryRow(
			"SELECT COUNT(*) FROM blockers WHERE issue_key = ? AND resolved_at IS NULL",
			issue.Key,
		)
		var count int
		if err := row.Scan(&count); err != nil || count > 0 {
			continue
		}

		blocker := &mem.Blocker{
			IssueKey:     issue.Key,
			Description:  fmt.Sprintf("Auto-detected: %s — %s", issue.Summary, issue.Status),
			BlockedSince: issue.Updated,
			Owner:        issue.Assignee,
			DaysBlocked:  0,
		}

		if err := h.Memory.SaveBlocker(ctx, blocker); err != nil {
			slog.Warn("auto-record blocker (board-aware)", "issue", issue.Key, "error", err)
		} else {
			slog.Info("auto-recorded blocker (board-aware)", "issue", issue.Key, "status", issue.Status, "assignee", issue.Assignee)
		}
	}
}

// reconcileBlockers checks stored blockers against current issue data.
// Auto-resolves blockers whose issues are no longer blocked.
// Uses heuristic blocking detection (use reconcileBlockersForBoard for board-aware).
func (h *Handlers) reconcileBlockers(ctx context.Context, issues []domain.Issue) {
	if h.Memory == nil {
		return
	}

	issueMap := make(map[string]bool, len(issues))
	for _, issue := range issues {
		issueMap[issue.Key] = issue.IsBlocked()
	}

	h.resolveUnblockedIssues(ctx, issueMap)
}

// reconcileBlockersForBoard checks stored blockers against current issue data
// using board-aware classification for accurate blocking detection.
func (h *Handlers) reconcileBlockersForBoard(ctx context.Context, boardID int, issues []domain.Issue) {
	if h.Memory == nil {
		return
	}

	// Fetch board config once
	isBlocked := func(status string) bool {
		return h.classifyStatus(ctx, boardID, status) == "blocked"
	}

	issueMap := make(map[string]bool, len(issues))
	for _, issue := range issues {
		issueMap[issue.Key] = isBlocked(issue.Status)
	}

	h.resolveUnblockedIssues(ctx, issueMap)
}

// resolveUnblockedIssues is shared by reconcileBlockers and reconcileBlockersForBoard.
func (h *Handlers) resolveUnblockedIssues(ctx context.Context, issueMap map[string]bool) {

	blockers, err := h.Memory.GetActiveBlockers(ctx)
	if err != nil {
		slog.Warn("reconcile blockers: get active", "error", err)
		return
	}

	for _, blocker := range blockers {
		if blocker.IssueKey == "" {
			continue // no Jira link, can't reconcile
		}

		isBlocked, found := issueMap[blocker.IssueKey]
		if !found {
			continue // not in current dataset, skip
		}

		if isBlocked {
			continue // still blocked, nothing to do
		}

		// Issue is no longer blocked — auto-resolve
		resolution := "Auto-reconciled: issue no longer blocked"
		if err := h.Memory.ResolveBlocker(ctx, blocker.ID, resolution); err != nil {
			slog.Warn("reconcile blocker: resolve", "issue", blocker.IssueKey, "error", err)
		} else {
			slog.Info("reconciled blocker — auto-resolved", "issue", blocker.IssueKey)
		}
	}
}

// reconcileRisks checks stored risks against current issue data.
// Resolves risks where the issue is done, or mitigates if recently updated.
func (h *Handlers) reconcileRisks(ctx context.Context, issues []domain.Issue) {
	if h.Memory == nil {
		return
	}

	issueMap := make(map[string]domain.Issue, len(issues))
	for _, issue := range issues {
		issueMap[issue.Key] = issue
	}

	risks, err := h.Memory.GetAllRisks(ctx, 100)
	if err != nil {
		slog.Warn("reconcile risks: get all", "error", err)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -staleThreshold)

	for _, risk := range risks {
		if risk.Status != "open" && risk.Status != "mitigating" {
			continue
		}

		issueKey := extractIssueKey(risk.Title)
		if issueKey == "" {
			issueKey = extractIssueKey(risk.Description)
		}
		if issueKey == "" {
			continue
		}

		issue, found := issueMap[issueKey]
		if !found {
			continue
		}

		now := time.Now()

		// Case 1: Issue is done → resolve risk
		if issue.IsDone() {
			risk.Status = "resolved"
			risk.ResolvedAt = &now
			if err := h.Memory.UpdateRisk(ctx, &risk); err != nil {
				slog.Warn("reconcile risk: resolve", "id", risk.ID, "issue", issueKey, "error", err)
			} else {
				slog.Info("reconciled risk — resolved (done)", "issue", issueKey, "risk", risk.Title)
			}
			continue
		}

		// Case 2: Issue recently updated (>7d) → mitigate risk
		if !issue.Updated.IsZero() && issue.Updated.After(cutoff) {
			risk.Status = "mitigating"
			risk.ResolvedAt = &now
			if err := h.Memory.UpdateRisk(ctx, &risk); err != nil {
				slog.Warn("reconcile risk: mitigate", "id", risk.ID, "issue", issueKey, "error", err)
			} else {
				slog.Info("reconciled risk — mitigated (updated)", "issue", issueKey, "risk", risk.Title)
			}
			continue
		}

		// Case 3: Issue is still stale at high priority — update timestamp only
		if risk.IdentifiedAt.IsZero() {
			risk.IdentifiedAt = now
			if err := h.Memory.UpdateRisk(ctx, &risk); err != nil {
				slog.Warn("reconcile risk: update timestamp", "id", risk.ID, "error", err)
			}
		}
	}
}

// extractIssueKey finds a Jira issue key (e.g. PROJ-123) from a string.
func extractIssueKey(s string) string {
	if s == "" {
		return ""
	}
	for _, part := range strings.Fields(s) {
		part = strings.Trim(part, ":,;-.")
		if looksLikeIssueKey(part) {
			return part
		}
	}
	return ""
}

// looksLikeIssueKey checks if a string matches a Jira key pattern (e.g. PROJ-123).
func looksLikeIssueKey(s string) bool {
	hyphenIdx := strings.LastIndex(s, "-")
	if hyphenIdx < 1 || hyphenIdx >= len(s)-1 {
		return false
	}
	for _, ch := range s[:hyphenIdx] {
		if ch < 'A' || ch > 'Z' {
			return false
		}
	}
	for _, ch := range s[hyphenIdx+1:] {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// staleThreshold is the number of days without update before a high-priority issue
// is considered stale and auto-recorded as a risk.
const staleThreshold = 7

// recordRisks scans for stale high-priority issues (Highest/Critical, not Done, >7d no update)
// and auto-records them as risks in PM memory.
func (h *Handlers) recordRisks(ctx context.Context, issues []domain.Issue) {
	if h.Memory == nil {
		return
	}

	cutoff := time.Now().AddDate(0, 0, -staleThreshold)

	for _, issue := range issues {
		if issue.IsDone() {
			continue
		}
		// Only flag Highest/Critical priority items as risks
		prio := strings.ToLower(issue.Priority)
		if prio != "highest" && prio != "critical" && prio != "blocker" {
			continue
		}
		if issue.Updated.After(cutoff) {
			continue // recently updated, not stale
		}
		if issue.Updated.IsZero() {
			continue // no timestamp, skip
		}

		// Dedup: check if risk title already exists for this issue
		riskTitle := fmt.Sprintf("Stale: %s — %s", issue.Key, issue.Summary[:min(len(issue.Summary), 80)])
		row := h.Memory.DB().QueryRow(
			"SELECT COUNT(*) FROM risks WHERE title LIKE ? AND status IN ('open', 'mitigating')",
			fmt.Sprintf("Stale: %s%%", issue.Key),
		)
		var count int
		if err := row.Scan(&count); err != nil || count > 0 {
			continue
		}

		risk := &mem.Risk{
			Title:       riskTitle,
			Description: fmt.Sprintf("Auto-detected: %s has been %s (priority: %s) with no update for %d+ days", issue.Key, issue.Status, issue.Priority, staleThreshold),
			Severity:    "high",
			Status:      "open",
			Owner:       issue.Assignee,
			IdentifiedAt: time.Now(),
		}

		if err := h.Memory.SaveRisk(ctx, risk); err != nil {
			slog.Warn("auto-record risk", "issue", issue.Key, "error", err)
		} else {
			slog.Info("auto-recorded stale risk", "issue", issue.Key, "priority", issue.Priority, "assignee", issue.Assignee)
		}
	}
}

// SearchIssues searches Jira issues using JQL.
func (h *Handlers) SearchIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jql, err := req.RequireString("jql")
	if err != nil {
		return mcputil.ErrInvalid("jql parameter is required"), nil
	}
	if err := validate.JQL(jql); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	maxResults := req.GetInt("max_results", 20)

	results, err := h.Jira.SearchIssues(ctx, jql, int(maxResults))
	if err != nil {
		return mcputil.ErrJira("Jira search", err), nil
	}
	if results == nil {
		return mcputil.TextResult("No results found."), nil
	}

	// Auto-record + reconcile blockers & risks against current Jira state
	h.recordBlockers(ctx, results.Issues)
	h.reconcileBlockers(ctx, results.Issues)
	h.recordRisks(ctx, results.Issues)
	h.reconcileRisks(ctx, results.Issues)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d issues\n", len(results.Issues)))
	for _, issue := range results.Issues {
		sb.WriteString(fmt.Sprintf("  %s - %s [%s]\n", issue.Key, issue.Summary, issue.Status))
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetIssue returns full details of a single Jira issue.
func (h *Handlers) GetIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	issue, err := h.Jira.GetIssue(ctx, key)
	if err != nil {
		return mcputil.ErrJira("get issue", err), nil
	}

	// Auto-record + reconcile for this issue
	issues := []domain.Issue{*issue}
	h.recordBlockers(ctx, issues)
	h.reconcileBlockers(ctx, issues)
	h.recordRisks(ctx, issues)
	h.reconcileRisks(ctx, issues)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s** - %s\n", issue.Key, issue.Summary))
	sb.WriteString(fmt.Sprintf("Type: %s | Status: %s | Priority: %s\n", issue.Type, issue.Status, issue.Priority))
	sb.WriteString(fmt.Sprintf("Assignee: %s\n", issue.Assignee))
	sb.WriteString(fmt.Sprintf("Description: %s\n", issue.Description))
	return mcputil.TextResult(sb.String()), nil
}

// GetBoards lists all accessible Jira boards.
func (h *Handlers) GetBoards(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boards, err := h.Jira.GetBoards(ctx)
	if err != nil {
		return mcputil.ErrJira("get boards", err), nil
	}

	var sb strings.Builder
	for _, b := range boards {
		sb.WriteString(fmt.Sprintf("%d: %s (%s)\n", b.ID, b.Name, b.Type))
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetSprintSummary returns the active sprint status for a board.
// CreateIssue creates a new Jira issue.
func (h *Handlers) CreateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project, err := req.RequireString("project")
	if err != nil {
		return mcputil.ErrInvalid("project parameter is required"), nil
	}
	if err := validate.ProjectKey(project); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	summary, err := req.RequireString("summary")
	if err != nil {
		return mcputil.ErrInvalid("summary parameter is required"), nil
	}

	issueType := req.GetString("issue_type", "Task")
	priority := req.GetString("priority", "")
	desc := req.GetString("description", "")
	assigneeID := req.GetString("assignee_id", "")
	labelsStr := req.GetString("labels", "")

	var labels []string
	if labelsStr != "" {
		labels, err = validate.Labels(labelsStr)
		if err != nil {
			return mcputil.ErrInvalid(err.Error()), nil
		}
	}

	input := &domain.CreateIssueInput{
		Project:     project,
		Summary:     summary,
		IssueType:   issueType,
		Description: desc,
		Priority:    priority,
		Assignee:    assigneeID,
		Labels:      labels,
	}

	issue, err := h.Jira.CreateIssue(ctx, input)
	if err != nil {
		return mcputil.ErrJira("create issue", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Created %s - %s", issue.Key, issue.Summary)), nil
}

// TransitionIssue transitions a Jira issue to a new status.
func (h *Handlers) TransitionIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	transitionID, err := req.RequireString("transition_id")
	if err != nil {
		return mcputil.ErrInvalid("transition_id parameter is required"), nil
	}
	if err := validate.TransitionID(transitionID); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	if err := h.Jira.TransitionIssue(ctx, key, transitionID); err != nil {
		return mcputil.ErrJira("transition issue", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Transitioned %s (transition: %s)", key, transitionID)), nil
}

// GetTransitions returns available transitions for an issue.
func (h *Handlers) GetTransitions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	transitions, err := h.Jira.GetTransitions(ctx, key)
	if err != nil {
		return mcputil.ErrJira("get transitions", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Available transitions for %s:\n", key))
	for _, t := range transitions {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", t.ID, t.Name))
	}
	return mcputil.TextResult(sb.String()), nil
}

// AssignIssue assigns a Jira issue to a user.
func (h *Handlers) AssignIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	accountID, err := req.RequireString("account_id")
	if err != nil {
		return mcputil.ErrInvalid("account_id parameter is required"), nil
	}

	if err := h.Jira.AssignIssue(ctx, key, accountID); err != nil {
		return mcputil.ErrJira("assign issue", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Assigned %s to %s", key, accountID)), nil
}

// FindUser searches for Jira users.
func (h *Handlers) FindUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return mcputil.ErrInvalid("query parameter is required"), nil
	}

	users, err := h.Jira.FindUser(ctx, query)
	if err != nil {
		return mcputil.ErrJira("find user", err), nil
	}

	var sb strings.Builder
	if len(users) == 0 {
		return mcputil.TextResult("No users found."), nil
	}
	for _, u := range users {
		sb.WriteString(fmt.Sprintf("%s: %s (%s)\n", u.AccountID, u.DisplayName, u.Email))
	}
	return mcputil.TextResult(sb.String()), nil
}

// AddComment adds a comment to a Jira issue.
func (h *Handlers) AddComment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	body, err := req.RequireString("body")
	if err != nil {
		return mcputil.ErrInvalid("body parameter is required"), nil
	}

	if err := h.Jira.AddComment(ctx, key, body); err != nil {
		return mcputil.ErrJira("add comment", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Comment added to %s", key)), nil
}

// GetSprints lists sprints for a board with optional state filter.
func (h *Handlers) GetSprints(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return mcputil.ErrInvalid("board_id parameter is required"), nil
	}
	state := req.GetString("state", "")

	sprints, err := h.Jira.GetSprints(ctx, int(boardID), state)
	if err != nil {
		return mcputil.ErrJira("get sprints", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprints for board %d (state: %s):\n", int(boardID), state))
	for _, s := range sprints {
		sb.WriteString(fmt.Sprintf("  %d: %s [%s]\n", s.ID, s.Name, s.State))
	}
	return mcputil.TextResult(sb.String()), nil
}

// StartSprint starts a sprint with start/end dates.
func (h *Handlers) StartSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintID, err := req.RequireInt("sprint_id")
	if err != nil {
		return mcputil.ErrInvalid("sprint_id parameter is required"), nil
	}
	startDate := req.GetString("start_date", "")
	endDate := req.GetString("end_date", "")
	if startDate == "" || endDate == "" {
		return mcputil.ErrInvalid("start_date and end_date are required"), nil
	}

	if err := h.Jira.StartSprint(ctx, int(sprintID), startDate, endDate); err != nil {
		return mcputil.ErrJira("start sprint", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Sprint %d started (%s → %s)", int(sprintID), startDate, endDate)), nil
}

// MoveIssuesToSprint moves issues into a sprint.
func (h *Handlers) MoveIssuesToSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintID, err := req.RequireInt("sprint_id")
	if err != nil {
		return mcputil.ErrInvalid("sprint_id parameter is required"), nil
	}
	issueKeysStr, err := req.RequireString("issue_keys")
	if err != nil {
		return mcputil.ErrInvalid("issue_keys parameter is required"), nil
	}
	issueKeys, err := validate.IssueKeys(issueKeysStr)
	if err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	if err := h.Jira.MoveIssuesToSprint(ctx, int(sprintID), issueKeys); err != nil {
		return mcputil.ErrJira("move to sprint", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Moved %d issues to sprint %d", len(issueKeys), int(sprintID))), nil
}

// LinkIssues creates a link between two issues.
func (h *Handlers) LinkIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	inwardKey, err := req.RequireString("inward_key")
	if err != nil {
		return mcputil.ErrInvalid("inward_key parameter is required"), nil
	}
	outwardKey, err := req.RequireString("outward_key")
	if err != nil {
		return mcputil.ErrInvalid("outward_key parameter is required"), nil
	}
	linkType, err := req.RequireString("link_type")
	if err != nil {
		return mcputil.ErrInvalid("link_type parameter is required"), nil
	}

	if err := h.Jira.LinkIssues(ctx, inwardKey, outwardKey, linkType); err != nil {
		return mcputil.ErrJira("link issues", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Linked %s → %s (%s)", inwardKey, outwardKey, linkType)), nil
}

// GetLinkTypes returns available link types.
func (h *Handlers) GetLinkTypes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	types, err := h.Jira.GetLinkTypes(ctx)
	if err != nil {
		return mcputil.ErrJira("get link types", err), nil
	}

	var sb strings.Builder
	for _, t := range types {
		sb.WriteString(fmt.Sprintf("%s: %s / %s\n", t.Name, t.Inward, t.Outward))
	}
	return mcputil.TextResult(sb.String()), nil
}

// AddWorklog logs time on an issue.
func (h *Handlers) AddWorklog(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	timeSpent, err := req.RequireString("time_spent")
	if err != nil {
		return mcputil.ErrInvalid("time_spent parameter is required"), nil
	}
	comment := req.GetString("comment", "")

	if err := h.Jira.AddWorklog(ctx, key, timeSpent, comment); err != nil {
		return mcputil.ErrJira("add worklog", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Logged %s on %s", timeSpent, key)), nil
}

// GetWorklogs returns worklogs for an issue.
func (h *Handlers) GetWorklogs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	logs, err := h.Jira.GetWorklogs(ctx, key)
	if err != nil {
		return mcputil.ErrJira("get worklogs", err), nil
	}

	var sb strings.Builder
	if len(logs) == 0 {
		return mcputil.TextResult("No worklogs found."), nil
	}
	for _, l := range logs {
		sb.WriteString(fmt.Sprintf("%s: %s - %s\n", l.Author, l.TimeSpent, l.Comment))
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetProjects lists all accessible projects.
func (h *Handlers) GetProjects(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projects, err := h.Jira.GetProjects(ctx)
	if err != nil {
		return mcputil.ErrJira("get projects", err), nil
	}

	var sb strings.Builder
	for _, p := range projects {
		sb.WriteString(fmt.Sprintf("%s: %s (%s)\n", p.Key, p.Name, p.Lead))
	}
	return mcputil.TextResult(sb.String()), nil
}

// --- Epic Handlers ---

// SetEpicLink adds issues to an epic by setting the epic parent link.
func (h *Handlers) SetEpicLink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	epicKey, err := req.RequireString("epic_key")
	if err != nil {
		return mcputil.ErrInvalid("epic_key parameter is required"), nil
	}
	if err := validate.IssueKey(epicKey); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	issueKeysParam, err := req.RequireString("issue_keys")
	if err != nil {
		return mcputil.ErrInvalid("issue_keys parameter is required"), nil
	}
	issueKeys, err := validate.IssueKeys(issueKeysParam)
	if err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	var failed []string
	for _, key := range issueKeys {
		if err := h.Jira.SetEpicLink(ctx, key, epicKey); err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", key, err))
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Added %d issue(s) to epic %s\n", len(issueKeys), epicKey))
	if len(failed) > 0 {
		sb.WriteString("Failed:\n")
		for _, f := range failed {
			sb.WriteString(fmt.Sprintf("  %s\n", f))
		}
	}
	return mcputil.TextResult(sb.String()), nil
}

// RemoveEpicLink removes issues from their epic by clearing the parent link.
func (h *Handlers) RemoveEpicLink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	issueKeysParam, err := req.RequireString("issue_keys")
	if err != nil {
		return mcputil.ErrInvalid("issue_keys parameter is required"), nil
	}
	issueKeys, err := validate.IssueKeys(issueKeysParam)
	if err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	var failed []string
	for _, key := range issueKeys {
		if err := h.Jira.RemoveEpicLink(ctx, key); err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", key, err))
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Removed %d issue(s) from their epic\n", len(issueKeys)))
	if len(failed) > 0 {
		sb.WriteString("Failed:\n")
		for _, f := range failed {
			sb.WriteString(fmt.Sprintf("  %s\n", f))
		}
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetEpicIssues lists issues in an epic using JQL.
func (h *Handlers) GetEpicIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	epicKey, err := req.RequireString("epic_key")
	if err != nil {
		return mcputil.ErrInvalid("epic_key parameter is required"), nil
	}
	if err := validate.IssueKey(epicKey); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	maxResults := int(req.GetInt("max_results", 50))

	jql := fmt.Sprintf(`"Epic Link" = %s`, epicKey)
	results, err := h.Jira.SearchIssues(ctx, jql, maxResults)
	if err != nil {
		jql = fmt.Sprintf("issue in parentEpicOf(%s)", epicKey)
		results, err = h.Jira.SearchIssues(ctx, jql, maxResults)
		if err != nil {
			return mcputil.ErrJira("get epic issues", err), nil
		}
	}

	if results == nil || len(results.Issues) == 0 {
		return mcputil.TextResult(fmt.Sprintf("No issues found in epic %s.", epicKey)), nil
	}

	// Auto-record + reconcile blockers & risks against current Jira state
	h.recordBlockers(ctx, results.Issues)
	h.reconcileBlockers(ctx, results.Issues)
	h.recordRisks(ctx, results.Issues)
	h.reconcileRisks(ctx, results.Issues)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Epic %s — %d issue(s):\n", epicKey, len(results.Issues)))
	for _, issue := range results.Issues {
		sb.WriteString(fmt.Sprintf("  %s - %s [%s]\n", issue.Key, issue.Summary, issue.Status))
	}
	return mcputil.TextResult(sb.String()), nil
}

// --- Version Handlers ---

// GetVersions lists project versions/releases.
func (h *Handlers) GetVersions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectKey, err := req.RequireString("project")
	if err != nil {
		return mcputil.ErrInvalid("project parameter is required"), nil
	}
	if err := validate.ProjectKey(projectKey); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	versions, err := h.Jira.GetVersions(ctx, projectKey)
	if err != nil {
		return mcputil.ErrJira("get versions", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Versions for %s:\n", projectKey))
	for _, v := range versions {
		status := "unreleased"
		if v.Released {
			status = fmt.Sprintf("released %s", v.ReleaseDate)
		}
		sb.WriteString(fmt.Sprintf("  %s — %s [%s]\n", v.Name, v.Description, status))
	}
	return mcputil.TextResult(sb.String()), nil
}

// CreateVersion creates a new project version.
func (h *Handlers) CreateVersion(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectKey, err := req.RequireString("project")
	if err != nil {
		return mcputil.ErrInvalid("project parameter is required"), nil
	}
	if err := validate.ProjectKey(projectKey); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	name, err := req.RequireString("name")
	if err != nil {
		return mcputil.ErrInvalid("name parameter is required"), nil
	}
	desc := req.GetString("description", "")

	version, err := h.Jira.CreateVersion(ctx, projectKey, name, desc)
	if err != nil {
		return mcputil.ErrJira("create version", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Created version %s in %s (ID: %s)", version.Name, projectKey, version.ID)), nil
}

// ReleaseVersion marks a version as released.
func (h *Handlers) ReleaseVersion(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	versionID, err := req.RequireString("version_id")
	if err != nil {
		return mcputil.ErrInvalid("version_id parameter is required"), nil
	}

	if err := h.Jira.ReleaseVersion(ctx, versionID); err != nil {
		return mcputil.ErrJira("release version", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Released version %s", versionID)), nil
}

// --- Component Handlers ---

// GetComponents lists project components.
func (h *Handlers) GetComponents(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectKey, err := req.RequireString("project")
	if err != nil {
		return mcputil.ErrInvalid("project parameter is required"), nil
	}
	if err := validate.ProjectKey(projectKey); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	components, err := h.Jira.GetComponents(ctx, projectKey)
	if err != nil {
		return mcputil.ErrJira("get components", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Components for %s:\n", projectKey))
	for _, c := range components {
		sb.WriteString(fmt.Sprintf("  %s — %s\n", c.Name, c.Lead))
	}
	return mcputil.TextResult(sb.String()), nil
}

// --- Attachment Handler ---

// GetBoardConfig returns the full board configuration including column layout and status mappings.
func (h *Handlers) GetBoardConfig(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return mcputil.ErrInvalid("board_id parameter is required"), nil
	}
	if err := validate.BoardID(int(boardID)); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	cfg, err := h.Jira.GetBoardConfiguration(ctx, int(boardID))
	if err != nil {
		return mcputil.ErrJira("get board configuration", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Board: %s (ID: %d, Type: %s)\n\n", cfg.Name, cfg.ID, cfg.Type))
	sb.WriteString("Columns:\n")
	for _, col := range cfg.ColumnConfig.Columns {
		sb.WriteString(fmt.Sprintf("  %s:\n", col.Name))
		for _, s := range col.Statuses {
			sb.WriteString(fmt.Sprintf("    - %s (ID: %s)\n", s.Name, s.ID))
		}
	}
	return mcputil.TextResult(sb.String()), nil
}

// classifyStatus uses board-aware classification when boardID is available,
// falling back to heuristic on the Issue method.
func (h *Handlers) classifyStatus(ctx context.Context, boardID int, statusName string) string {
	cfg, err := h.Jira.GetBoardConfiguration(ctx, boardID)
	if err != nil || cfg == nil {
		return domain.HeuristicClassify(statusName)
	}
	return cfg.StatusCategory(statusName)
}

// GetAttachments lists attachments on an issue.
func (h *Handlers) GetAttachments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	attachments, err := h.Jira.GetAttachments(ctx, key)
	if err != nil {
		return mcputil.ErrJira("get attachments", err), nil
	}

	var sb strings.Builder
	if len(attachments) == 0 {
		return mcputil.TextResult("No attachments found."), nil
	}
	sb.WriteString(fmt.Sprintf("Attachments for %s:\n", key))
	for _, a := range attachments {
		sb.WriteString(fmt.Sprintf("  %s (%s, %d bytes) — %s\n", a.Filename, a.MimeType, a.Size, a.Author))
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetSprintSummary returns the active sprint status for a board.
func (h *Handlers) GetSprintSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return mcputil.ErrInvalid("board_id parameter is required"), nil
	}
	if err := validate.BoardID(int(boardID)); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, int(boardID))
	if err != nil {
		return mcputil.ErrJira("get active sprints", err), nil
	}

	if len(sprints) == 0 {
		return mcputil.TextResult("No active sprints found for this board."), nil
	}

	sprintResult := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprintResult.ID)
	if err != nil {
		return mcputil.ErrJira("get sprint issues", err), nil
	}

	// Auto-record + reconcile blockers/risks + snapshot to PM memory
	// Sprint summary has board context, so we use board-aware classification
	h.recordBlockersForBoard(ctx, int(boardID), issues)
	h.reconcileBlockersForBoard(ctx, int(boardID), issues)
	h.recordRisks(ctx, issues)
	h.reconcileRisks(ctx, issues)
	h.recordSnapshot(ctx, sprintResult.Name, int(boardID), issues)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprint: %s (ID: %d)\n", sprintResult.Name, sprintResult.ID))
	sb.WriteString(fmt.Sprintf("Goals: %s\n", sprintResult.Goal))
	sb.WriteString(fmt.Sprintf("Start: %s | End: %s\n", sprintResult.StartDate, sprintResult.EndDate))
	sb.WriteString(fmt.Sprintf("Issues: %d\n", len(issues)))

	for _, i := range issues {
		sb.WriteString(fmt.Sprintf("  %s - %s [%s]\n", i.Key, i.Summary, i.Status))
	}
	return mcputil.TextResult(sb.String()), nil
}

// ReconcileMemory performs a full sweep: fetches every issue key from stored blockers/risks
// and reconciles them with current Jira state. Auto-resolves resolved items.
func (h *Handlers) ReconcileMemory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return mcputil.ErrorResult("PM memory not configured. Set PM_MEMORY_DB_PATH."), nil
	}

	var resolvedBlockers, resolvedRisks, mitigatedRisks int

	// Get all active blockers with issue keys
	blockers, err := h.Memory.GetActiveBlockers(ctx)
	if err != nil {
		return mcputil.ErrInternal("reconcile: get blockers", err), nil
	}

	// Get all open/mitigating risks
	allRisks, err := h.Memory.GetAllRisks(ctx, 200)
	if err != nil {
		return mcputil.ErrInternal("reconcile: get risks", err), nil
	}

	// Collect unique issue keys from both stores
	keys := make(map[string]bool)
	for _, b := range blockers {
		if b.IssueKey != "" {
			keys[b.IssueKey] = false
		}
	}
	for _, r := range allRisks {
		if r.Status != "open" && r.Status != "mitigating" {
			continue
		}
		issueKey := extractIssueKey(r.Title)
		if issueKey == "" {
			issueKey = extractIssueKey(r.Description)
		}
		if issueKey != "" {
			keys[issueKey] = true
		}
	}

	if len(keys) == 0 {
		return mcputil.TextResult("No stored items to reconcile. All data is consistent."), nil
	}

	// Fetch each issue and reconcile
	total := len(keys)
	processed := 0
	for issueKey := range keys {
		processed++
		issue, err := h.Jira.GetIssue(ctx, issueKey)
		if err != nil {
			slog.Warn("reconcile: fetch issue", "issue", issueKey, "error", err)
			continue
		}

		issues := []domain.Issue{*issue}
		prevBlockers, _ := h.Memory.GetActiveBlockers(ctx)
		prevCount := len(prevBlockers)

		h.reconcileBlockers(ctx, issues)
		h.reconcileRisks(ctx, issues)

		// Count changes
		currBlockers, _ := h.Memory.GetActiveBlockers(ctx)
		if len(currBlockers) < prevCount {
			resolvedBlockers += prevCount - len(currBlockers)
		}

		slog.Info("reconcile progress", "processed", processed, "total", total)
	}

	var sb strings.Builder
	sb.WriteString("# PM Memory Reconciliation\n\n")
	sb.WriteString(fmt.Sprintf("Issues checked: %d\n", total))
	sb.WriteString(fmt.Sprintf("Blockers resolved: %d\n", resolvedBlockers))
	sb.WriteString(fmt.Sprintf("Risks resolved/mitigated: %d\n", resolvedRisks+mitigatedRisks))
	sb.WriteString("\nMemory is now consistent with current Jira state.\n")

	return mcputil.TextResult(sb.String()), nil
}
