// Package mcp provides MCP tool handlers for the sprint/PM module.
package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/port"
	memory "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

// Handlers holds dependencies for sprint/PM MCP tool handlers.
type Handlers struct {
	Memory        memory.Store
	SprintService port.Inbound
	AI            port.AIProvider
	Config        *config.Config
	Cache         Cache
	Error         *mcputil.ErrorHandler
}

// Cache interface for sprint module caching.
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Available() bool
}

// NewHandlers creates a new sprint MCP handlers instance.
func NewHandlers(
	memStore memory.Store,
	sprintSvc port.Inbound,
	ai port.AIProvider,
	cfg *config.Config,
	cache Cache,
) *Handlers {
	return &Handlers{
		Memory:        memStore,
		SprintService: sprintSvc,
		AI:            ai,
		Config:        cfg,
		Cache:         cache,
		Error:         mcputil.NewErrorHandler(nil),
	}
}

// --- Health ---

// Health returns server version and status.
func (h *Handlers) Health(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("zara-jira-mcp v0.4.0 | sprint module | status: ok"), nil
}

// --- Forecast ---

// PMForecast runs a Monte Carlo simulation to predict sprint completion.
func (h *Handlers) PMForecast(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available. Set PM_MEMORY_DB_PATH."), nil
	}

	boardID := int(req.GetInt("board_id", 0))
	remaining := int(req.GetInt("remaining_items", 0))

	result, err := h.SprintService.Forecast(ctx, boardID, remaining)
	if err != nil {
		return h.Error.Wrap("forecast", err), nil
	}

	var sb strings.Builder
	sb.WriteString("# Sprint Forecast\n\n")
	sb.WriteString(fmt.Sprintf("Based on **%d** historical sprints\n", result.Simulations))
	sb.WriteString(fmt.Sprintf("Remaining items: **%d**\n\n", result.Remaining))

	sb.WriteString("| Confidence | Sprints |\n")
	sb.WriteString("|------------|--------|\n")
	for _, pct := range []int{50, 70, 85, 95} {
		if val, ok := result.Percentiles[pct]; ok {
			sb.WriteString(fmt.Sprintf("| %d%% | %.1f |\n", pct, val))
		}
	}

	sb.WriteString(fmt.Sprintf("\nMean: **%.1f** sprints (range: %d–%d)", result.MeanSprints, result.MinSprints, result.MaxSprints))

	return mcputil.TextResult(sb.String()), nil
}

// --- PM Quick Actions ---

// PMQuickStatus returns a quick project status overview from memory.
func (h *Handlers) PMQuickStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	now := time.Now()
	todayStr := now.Format("2006-01-02")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# PM Status — %s\n\n", todayStr))

	// Recent snapshot
	snapshots, err := h.Memory.GetSprintSnapshots(ctx, 0, 1)
	if err == nil && len(snapshots) > 0 {
		s := snapshots[0]
		sb.WriteString(fmt.Sprintf("**Latest sprint:** %s\n", s.SprintName))
		sb.WriteString(fmt.Sprintf("Completion: %d/%d (%.0f%%)\n", s.Done, s.TotalIssues, s.CompletionRate))
		if s.IsZombie() {
			sb.WriteString("⚠️ Zombie sprint: carryover >30%\n")
		}
	}

	// Active risks
	risks, err := h.Memory.GetOpenRisks(ctx)
	if err == nil && len(risks) > 0 {
		sb.WriteString(fmt.Sprintf("\n**Open risks:** %d\n", len(risks)))
		for _, r := range risks {
			sb.WriteString(fmt.Sprintf("  [%s] %s — %s\n", strings.ToUpper(r.Severity), r.Title, r.Owner))
		}
	}

	// Active blockers
	blockers, err := h.Memory.GetActiveBlockers(ctx)
	if err == nil && len(blockers) > 0 {
		sb.WriteString(fmt.Sprintf("\n**Active blockers:** %d\n", len(blockers)))
		for _, b := range blockers {
			sb.WriteString(fmt.Sprintf("  %s — %s\n", b.IssueKey, b.Description))
		}
	}

	// Pending actions
	actions, err := h.Memory.GetPendingActionItems(ctx)
	if err == nil && len(actions) > 0 {
		sb.WriteString(fmt.Sprintf("\n**Pending actions:** %d\n", len(actions)))
		for _, a := range actions {
			sb.WriteString(fmt.Sprintf("  %s\n", a.Description))
		}
	}

	if risks == nil && blockers == nil && actions == nil {
		sb.WriteString("\nAll clear. No open risks, blockers, or pending actions.\n")
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMCreate creates a work item. Since we don't have Jira access here,
// it records the intent in memory.
func (h *Handlers) PMCreate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return mcputil.ErrInvalid("title parameter is required"), nil
	}
	desc := req.GetString("description", "")
	project := req.GetString("project", "")
	assignee := req.GetString("assignee", "")

	var notes strings.Builder
	notes.WriteString(fmt.Sprintf("Task: %s\n", title))
	if desc != "" {
		notes.WriteString(fmt.Sprintf("Description: %s\n", desc))
	}
	if project != "" {
		notes.WriteString(fmt.Sprintf("Project: %s\n", project))
	}
	if assignee != "" {
		notes.WriteString(fmt.Sprintf("Assignee: %s\n", assignee))
	}

	// Record in memory as a meeting note for now
	_ = h.Memory.SaveMeetingNote(ctx, &memory.MeetingNote{
		MeetingType: "adhoc",
		Date:        time.Now(),
		Notes:       notes.String(),
		ActionItems: title,
	})

	return mcputil.TextResult(fmt.Sprintf("[PM] Created work item: %s", title)), nil
}

// PMDecide records a decision in PM memory.
func (h *Handlers) PMDecide(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	what, err := req.RequireString("what")
	if err != nil {
		return mcputil.ErrInvalid("what parameter is required"), nil
	}
	who := req.GetString("who", "")
	why := req.GetString("why", "")

	decision := &memory.Decision{
		Title:     what[:min(len(what), 200)],
		Decision:  what,
		Rationale: why,
		MadeBy:    who,
		MadeAt:    time.Now(),
	}

	if err := h.Memory.SaveDecision(ctx, decision); err != nil {
		return h.Error.WrapInternal("save decision", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("[PM] Decision recorded: %s", what)), nil
}

// PMRisk records a risk in PM memory.
func (h *Handlers) PMRisk(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	what, err := req.RequireString("what")
	if err != nil {
		return mcputil.ErrInvalid("what parameter is required"), nil
	}
	severity := req.GetString("severity", "medium")
	owner := req.GetString("owner", "")

	risk := &memory.Risk{
		Title:        what[:min(len(what), 200)],
		Description:  what,
		Severity:     severity,
		Status:       "open",
		Owner:        owner,
		IdentifiedAt: time.Now(),
	}

	if err := h.Memory.SaveRisk(ctx, risk); err != nil {
		return h.Error.WrapInternal("save risk", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("[PM] Risk recorded: %s (severity: %s)", what, severity)), nil
}

// PMNext suggests the next high-priority PM action based on memory state.
func (h *Handlers) PMNext(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var suggestions []string

	// Check pending action items
	actions, err := h.Memory.GetPendingActionItems(ctx)
	if err == nil && len(actions) > 0 {
		suggestions = append(suggestions, fmt.Sprintf("Follow up on %d pending action items (use pm_actions)", len(actions)))
	}

	// Check open risks
	risks, err := h.Memory.GetOpenRisks(ctx)
	if err == nil && len(risks) > 0 {
		for _, r := range risks {
			if r.Severity == "critical" || r.Severity == "high" {
				suggestions = append(suggestions, fmt.Sprintf("Escalate high/critical risk: %s", r.Title))
			}
		}
	}

	// Check active blockers
	blockers, err := h.Memory.GetActiveBlockers(ctx)
	if err == nil && len(blockers) > 0 {
		suggestions = append(suggestions, fmt.Sprintf("Resolve %d active blockers", len(blockers)))
	}

	// Check open dependencies
	deps, err := h.Memory.GetOpenDependencies(ctx)
	if err == nil && len(deps) > 0 {
		suggestions = append(suggestions, fmt.Sprintf("Resolve %d open dependencies", len(deps)))
	}

	if len(suggestions) == 0 {
		return mcputil.TextResult("No urgent items. Try running a sprint health check or planning next sprint."), nil
	}

	return mcputil.TextResult("[PM] Next actions:\n  " + strings.Join(suggestions, "\n  ")), nil
}

// --- PM Memory Tools ---

// PMSnapshotSprint snapshots the current sprint state into memory.
func (h *Handlers) PMSnapshotSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)
	sprintName := req.GetString("sprint_name", "")

	snapshot := &memory.SprintSnapshot{
		SprintName:   sprintName,
		BoardID:      int(boardID),
		SnapshotDate: time.Now(),
		TotalIssues:  req.GetInt("total_issues", 0),
		Done:         req.GetInt("done", 0),
		InProgress:   req.GetInt("in_progress", 0),
		Todo:         req.GetInt("todo", 0),
		Blocked:      req.GetInt("blocked", 0),
		Carryover:    req.GetInt("carryover", 0),
		Velocity:     req.GetInt("velocity", 0),
		Notes:        req.GetString("notes", ""),
	}

	snapshot.CompletionRate = snapshot.CalculateCompletionRate()

	if err := h.Memory.SaveSprintSnapshot(ctx, snapshot); err != nil {
		return h.Error.WrapInternal("save snapshot", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Sprint snapshot saved: %s | %d/%d done (%.0f%%)",
		sprintName, snapshot.Done, snapshot.TotalIssues, snapshot.CompletionRate)), nil
}

// PMRecordBlocker records an impediment/blocker.
func (h *Handlers) PMRecordBlocker(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	desc, err := req.RequireString("description")
	if err != nil {
		return mcputil.ErrInvalid("description parameter is required"), nil
	}
	issueKey := req.GetString("issue_key", "")
	owner := req.GetString("owner", "")

	blocker := &memory.Blocker{
		IssueKey:     issueKey,
		Description:  desc,
		BlockedSince: time.Now(),
		Owner:        owner,
		DaysBlocked:  0,
	}

	if err := h.Memory.SaveBlocker(ctx, blocker); err != nil {
		return h.Error.WrapInternal("save blocker", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("[PM] Blocker recorded: %s", desc)), nil
}

// PMRecordDecision records a decision with context.
func (h *Handlers) PMRecordDecision(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return mcputil.ErrInvalid("title parameter is required"), nil
	}
	decision, err := req.RequireString("decision")
	if err != nil {
		return mcputil.ErrInvalid("decision parameter is required"), nil
	}
	contextStr := req.GetString("context", "")
	rationale := req.GetString("rationale", "")
	madeBy := req.GetString("made_by", "")
	tags := req.GetString("tags", "")

	d := &memory.Decision{
		Title:     title,
		Context:   contextStr,
		Decision:  decision,
		Rationale: rationale,
		MadeBy:    madeBy,
		MadeAt:    time.Now(),
		Tags:      tags,
	}

	if err := h.Memory.SaveDecision(ctx, d); err != nil {
		return h.Error.WrapInternal("save decision", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("[PM] Decision recorded: %s — %s", title, decision)), nil
}

// PMRecordRisk records a risk with full details.
func (h *Handlers) PMRecordRisk(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return mcputil.ErrInvalid("title parameter is required"), nil
	}
	desc := req.GetString("description", "")
	severity := req.GetString("severity", "medium")
	owner := req.GetString("owner", "")
	mitigation := req.GetString("mitigation", "")

	r := &memory.Risk{
		Title:        title,
		Description:  desc,
		Severity:     severity,
		Status:       "open",
		Owner:        owner,
		Mitigation:   mitigation,
		IdentifiedAt: time.Now(),
	}

	if err := h.Memory.SaveRisk(ctx, r); err != nil {
		return h.Error.WrapInternal("save risk", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("[PM] Risk recorded: %s (severity: %s)", title, severity)), nil
}

// PMRecordRetro records a sprint retrospective.
func (h *Handlers) PMRecordRetro(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintName, err := req.RequireString("sprint_name")
	if err != nil {
		return mcputil.ErrInvalid("sprint_name parameter is required"), nil
	}
	wentWell := req.GetString("went_well", "")
	improvements := req.GetString("improvements", "")
	actionItems := req.GetString("action_items", "")

	retro := &memory.Retrospective{
		SprintName:   sprintName,
		Date:         time.Now(),
		WentWell:     wentWell,
		Improvements: improvements,
		ActionItems:  actionItems,
		Status:       "open",
	}

	if err := h.Memory.SaveRetrospective(ctx, retro); err != nil {
		return h.Error.WrapInternal("save retro", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("[PM] Retro recorded for %s", sprintName)), nil
}

// PMRecordMeeting records meeting notes.
func (h *Handlers) PMRecordMeeting(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	meetingType, err := req.RequireString("meeting_type")
	if err != nil {
		return mcputil.ErrInvalid("meeting_type parameter is required"), nil
	}
	notes := req.GetString("notes", "")
	attendees := req.GetString("attendees", "")
	decisions := req.GetString("decisions", "")
	actionItems := req.GetString("action_items", "")
	sprintName := req.GetString("sprint_name", "")

	m := &memory.MeetingNote{
		MeetingType: meetingType,
		Date:        time.Now(),
		Attendees:   attendees,
		Notes:       notes,
		Decisions:   decisions,
		ActionItems: actionItems,
		SprintName:  sprintName,
	}

	if err := h.Memory.SaveMeetingNote(ctx, m); err != nil {
		return h.Error.WrapInternal("save meeting note", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("[PM] Meeting note saved (%s)", meetingType)), nil
}

// PMRiskDashboard shows all open risks.
func (h *Handlers) PMRiskDashboard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	risks, err := h.Memory.GetAllRisks(ctx, 50)
	if err != nil {
		return h.Error.WrapInternal("get risks", err), nil
	}

	var sb strings.Builder
	sb.WriteString("# Risk Dashboard\n\n")
	if len(risks) == 0 {
		return mcputil.TextResult("No risks recorded."), nil
	}

	// Group by severity
	bySeverity := make(map[string][]memory.Risk)
	for _, r := range risks {
		bySeverity[r.Severity] = append(bySeverity[r.Severity], r)
	}

	for _, sev := range []string{"critical", "high", "medium", "low"} {
		items, ok := bySeverity[sev]
		if !ok {
			continue
		}
		sb.WriteString(fmt.Sprintf("**[%s]** (%d)\n", strings.ToUpper(sev), len(items)))
		for _, r := range items {
			sb.WriteString(fmt.Sprintf("  %s — %s [%s]\n", r.Title, r.Owner, r.Status))
		}
		sb.WriteString("\n")
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMBlockers shows active blockers.
func (h *Handlers) PMBlockers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	showHistory := req.GetBool("show_history", false)

	var blockers []memory.Blocker
	var err error

	if showHistory {
		blockers, err = h.Memory.GetBlockerHistory(ctx, 20)
	} else {
		blockers, err = h.Memory.GetActiveBlockers(ctx)
	}

	if err != nil {
		return h.Error.WrapInternal("get blockers", err), nil
	}

	var sb strings.Builder
	if showHistory {
		sb.WriteString("# Blocker History\n\n")
	} else {
		sb.WriteString("# Active Blockers\n\n")
	}
	if len(blockers) == 0 {
		return mcputil.TextResult("No blockers found."), nil
	}

	for _, b := range blockers {
		sb.WriteString(fmt.Sprintf("- %s (%s) — %s\n", b.IssueKey, b.Owner, b.Description))
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMDecisions shows recent decisions.
func (h *Handlers) PMDecisions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	limit := int(req.GetInt("limit", 10))

	decisions, err := h.Memory.GetDecisions(ctx, limit)
	if err != nil {
		return h.Error.WrapInternal("get decisions", err), nil
	}

	var sb strings.Builder
	sb.WriteString("# Recent Decisions\n\n")
	if len(decisions) == 0 {
		return mcputil.TextResult("No decisions recorded."), nil
	}

	for _, d := range decisions {
		sb.WriteString(fmt.Sprintf("- **%s** (by %s)\n  %s\n", d.Title, d.MadeBy, d.Decision))
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMActionItems shows pending action items.
func (h *Handlers) PMActionItems(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	actions, err := h.Memory.GetPendingActionItems(ctx)
	if err != nil {
		return h.Error.WrapInternal("get action items", err), nil
	}

	var sb strings.Builder
	sb.WriteString("# Pending Action Items\n\n")
	if len(actions) == 0 {
		return mcputil.TextResult("No pending action items. Good job! 🎯"), nil
	}

	for _, a := range actions {
		sb.WriteString(fmt.Sprintf("- %s\n", a.Description))
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMDependencies shows open dependencies.
func (h *Handlers) PMDependencies(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	issueKey := req.GetString("issue_key", "")

	var deps []memory.Dependency
	var err error

	if issueKey != "" {
		deps, err = h.Memory.GetDependenciesForIssue(ctx, issueKey)
	} else {
		deps, err = h.Memory.GetOpenDependencies(ctx)
	}

	if err != nil {
		return h.Error.WrapInternal("get dependencies", err), nil
	}

	var sb strings.Builder
	sb.WriteString("# Dependencies\n\n")
	if len(deps) == 0 {
		return mcputil.TextResult("No open dependencies."), nil
	}

	for _, d := range deps {
		sb.WriteString(fmt.Sprintf("- %s → %s [%s]\n", d.FromIssueKey, d.ToIssueKey, d.DependencyType))
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMSprintHealth returns sprint health from memory.
func (h *Handlers) PMSprintHealth(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}

	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available. Set PM_MEMORY_DB_PATH."), nil
	}

	result, err := h.SprintService.CalculateHealth(ctx, boardID)
	if err != nil {
		// Fallback: show historical health scores.
		scores, hErr := h.Memory.GetHealthScores(ctx, boardID, 5)
		if hErr != nil {
			return h.Error.Wrap("health", err), nil
		}
		var sb strings.Builder
		sb.WriteString("# Sprint Health History\n\n")
		sb.WriteString(fmt.Sprintf("Live health unavailable (%s). Showing stored scores:\n\n", err.Error()))
		if len(scores) == 0 {
			return mcputil.TextResult("No health data yet. Complete a sprint and snapshot first."), nil
		}
		for _, s := range scores {
			sb.WriteString(fmt.Sprintf("- %s: **%d/100** (%s)\n", s.SprintName, s.OverallScore, s.Rating()))
		}
		return mcputil.TextResult(sb.String()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Sprint: %s\n\n", result.SprintName))
	sb.WriteString(fmt.Sprintf("**Overall: %d/100** — %s\n\n", result.Score, result.Rating))
	sb.WriteString("| Dimension | Score | Max |\n")
	sb.WriteString("|-----------|-------|-----|\n")
	sb.WriteString(fmt.Sprintf("| Velocity | %d | 25 |\n", result.VelocityScore))
	sb.WriteString(fmt.Sprintf("| Blockers | %d | 25 |\n", result.BlockerScore))
	sb.WriteString(fmt.Sprintf("| Scope | %d | 25 |\n", result.ScopeScore))
	sb.WriteString(fmt.Sprintf("| Team | %d | 25 |\n", result.TeamScore))
	sb.WriteString(fmt.Sprintf("\n⚡ Weakest dimension: **%s**\n", result.WeakestDim))

	return mcputil.TextResult(sb.String()), nil
}

// PMVelocityTrend returns velocity trend over recent sprints.
func (h *Handlers) PMVelocityTrend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available. Set PM_MEMORY_DB_PATH."), nil
	}

	trend, err := h.SprintService.VelocityTrend(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("velocity trend", err), nil
	}

	return mcputil.TextResult(trend), nil
}

// PMAntiPatterns detects and returns anti-patterns for a board.
func (h *Handlers) PMAntiPatterns(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available. Set PM_MEMORY_DB_PATH."), nil
	}

	patterns, err := h.SprintService.DetectAntiPatterns(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("anti-patterns", err), nil
	}

	var sb strings.Builder
	sb.WriteString("# Anti-Pattern Scan\n\n")
	if len(patterns) == 0 {
		sb.WriteString("✅ No anti-patterns detected. Team looks healthy!\n")
		return mcputil.TextResult(sb.String()), nil
	}

	for _, p := range patterns {
		severityIcon := "🟡"
		switch p.Severity {
		case "High":
			severityIcon = "🔴"
		case "Medium":
			severityIcon = "🟠"
		}
		sb.WriteString(fmt.Sprintf("%s **%s** (%s)\n", severityIcon, p.Name, p.Severity))
		sb.WriteString(fmt.Sprintf("   %s\n", p.Description))
		sb.WriteString(fmt.Sprintf("   💡 %s\n\n", p.Suggestion))
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMFlowMetrics shows WIP, throughput, cycle time, and completion rate.
func (h *Handlers) PMFlowMetrics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available. Set PM_MEMORY_DB_PATH."), nil
	}

	m, err := h.SprintService.FlowMetrics(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("flow metrics", err), nil
	}

	var sb strings.Builder
	sb.WriteString("# Flow Metrics\n\n")
	sb.WriteString(fmt.Sprintf("**Trend**: %s\n\n", m.Trend))
	sb.WriteString("| Metric | Value |\n")
	sb.WriteString("|--------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Current WIP | %d |\n", m.CurrentWIP))
	sb.WriteString(fmt.Sprintf("| Total Issues | %d |\n", m.TotalIssues))
	sb.WriteString(fmt.Sprintf("| Done | %d |\n", m.DoneIssues))
	sb.WriteString(fmt.Sprintf("| Blocked | %d |\n", m.BlockedIssues))
	sb.WriteString(fmt.Sprintf("| Completion Rate | %.0f%% |\n", m.CompletionPct))
	sb.WriteString(fmt.Sprintf("| Avg Throughput | %.1f items/sprint |\n", m.AvgThroughput))
	sb.WriteString(fmt.Sprintf("| Inferred Cycle Time | %.1f days |\n", m.AvgCycleTime))

	// Health signal based on WIP vs throughput.
	if m.CurrentWIP > 0 && m.AvgThroughput > 0 {
		ratio := float64(m.CurrentWIP) / m.AvgThroughput
		sb.WriteString("\n**Flow Health**: ")
		switch {
		case ratio > 1.5:
			sb.WriteString("⚠️ High WIP relative to throughput. Consider limiting WIP to reduce cycle time.\n")
		case ratio > 0.8:
			sb.WriteString("✓ Healthy flow. WIP is balanced with delivery capacity.\n")
		default:
			sb.WriteString("✓ Low WIP — team has room to take on more work.\n")
		}
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMSprintCompare compares current vs previous sprint.
func (h *Handlers) PMSprintCompare(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available."), nil
	}

	result, err := h.SprintService.SprintCompare(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("sprint compare", err), nil
	}
	return mcputil.TextResult(result), nil
}

// PMPredictability shows sprint predictability score.
func (h *Handlers) PMPredictability(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available."), nil
	}

	result, err := h.SprintService.Predictability(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("predictability", err), nil
	}
	return mcputil.TextResult(result), nil
}

// PMScorecard shows comprehensive sprint scorecard (0-100).
func (h *Handlers) PMScorecard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available."), nil
	}

	result, err := h.SprintService.Scorecard(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("scorecard", err), nil
	}
	return mcputil.TextResult(result), nil
}

// PMCalibration shows forecast accuracy — committed vs delivered over time.
func (h *Handlers) PMCalibration(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available."), nil
	}

	result, err := h.SprintService.Calibration(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("calibration", err), nil
	}
	return mcputil.TextResult(result), nil
}

// PMSetSprintGoal defines a sprint goal with key results.
func (h *Handlers) PMSetSprintGoal(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	goal := req.GetString("goal", "")
	if goal == "" {
		return mcputil.ErrorResult("goal is required"), nil
	}
	boardID := int(req.GetInt("board_id", 0))
	sprintName := req.GetString("sprint_name", "")
	keyResults := req.GetString("key_results", "")

	if sprintName == "" {
		sprintName = "current"
	}

	g := &memory.SprintGoal{
		SprintName: sprintName,
		BoardID:    boardID,
		Goal:       goal,
		KeyResults: keyResults,
		Status:     "active",
		CreatedAt:  time.Now(),
	}

	if err := h.Memory.SaveSprintGoal(ctx, g); err != nil {
		return h.Error.WrapInternal("save sprint goal", err), nil
	}

	var sb strings.Builder
	sb.WriteString("# Sprint Goal Set\n\n")
	sb.WriteString(fmt.Sprintf("**Sprint**: %s\n", sprintName))
	sb.WriteString(fmt.Sprintf("**Goal**: %s\n", goal))
	if keyResults != "" {
		sb.WriteString(fmt.Sprintf("**Key Results**:\n%s\n", keyResults))
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMSprintGoals shows active sprint goals or goal history.
func (h *Handlers) PMSprintGoals(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	showHistory := req.GetBool("show_history", false)

	var goals []memory.SprintGoal
	var err error

	if showHistory {
		goals, err = h.Memory.GetGoalHistory(ctx, boardID, 20)
	} else {
		goals, err = h.Memory.GetActiveGoals(ctx, boardID)
	}
	if err != nil {
		return h.Error.WrapInternal("get goals", err), nil
	}

	var sb strings.Builder
	if showHistory {
		sb.WriteString("# Sprint Goal History\n\n")
	} else {
		sb.WriteString("# Active Sprint Goals\n\n")
	}

	if len(goals) == 0 {
		sb.WriteString("No goals found. Set one with pm_set_sprint_goal.")
		return mcputil.TextResult(sb.String()), nil
	}

	for _, g := range goals {
		statusIcon := "🟢"
		switch g.Status {
		case "achieved":
			statusIcon = "✅"
		case "partially_achieved":
			statusIcon = "🟡"
		case "missed":
			statusIcon = "🔴"
		}
		sb.WriteString(fmt.Sprintf("%s **%s** — %s\n", statusIcon, g.SprintName, g.Goal))
		sb.WriteString(fmt.Sprintf("   Status: *%s*\n", g.Status))
		if g.KeyResults != "" {
			sb.WriteString(fmt.Sprintf("   KRs:\n%s\n", g.KeyResults))
		}
		if g.Outcome != "" {
			sb.WriteString(fmt.Sprintf("   Outcome: %s\n", g.Outcome))
		}
		sb.WriteString("\n")
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMGoalCheck evaluates sprint goal progress using Jira data.
func (h *Handlers) PMGoalCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available."), nil
	}

	var sb strings.Builder
	sb.WriteString("# Sprint Goal Check\n\n")

	// Get active goals.
	goals, err := h.Memory.GetActiveGoals(ctx, boardID)
	if err != nil {
		return h.Error.WrapInternal("get goals", err), nil
	}
	if len(goals) == 0 {
		sb.WriteString("No active goals set. Use `pm_set_sprint_goal` to define one.\n")
		return mcputil.TextResult(sb.String()), nil
	}

	// Get sprint health for status context.
	health, err := h.SprintService.CalculateHealth(ctx, boardID)
	if err != nil {
		sb.WriteString(fmt.Sprintf("⚠️ Could not compute sprint health: %s\n\n", err.Error()))
	} else {
		sb.WriteString(fmt.Sprintf("**Overall Health**: %d/100 (%s)\n", health.Score, health.Rating))
		sb.WriteString(fmt.Sprintf("**Velocity**: %d/25 | **Blockers**: %d/25 | **Scope**: %d/25 | **Team**: %d/25\n\n",
			health.VelocityScore, health.BlockerScore, health.ScopeScore, health.TeamScore))
	}

	for i, g := range goals {
		sb.WriteString(fmt.Sprintf("**Goal %d**: %s\n", i+1, g.Goal))
		if g.KeyResults != "" {
			sb.WriteString(fmt.Sprintf("   KRs: %s\n", g.KeyResults))
		}

		// AI-based progress assessment if available.
		if h.AI != nil {
			prompt := fmt.Sprintf(`Evaluate this sprint goal progress:
Goal: %s
Key Results: %s
Sprint Health: %d/100 (%s)
Velocity subscore: %d/25
Blocker subscore: %d/25

Give a brief assessment (1-2 sentences) of whether this goal is on track, at risk, or off track.`,
				g.Goal, g.KeyResults, health.Score, health.Rating, health.VelocityScore, health.BlockerScore)
			assessment, aiErr := h.AI.Complete(ctx,
				"You are a Scrum Master evaluating sprint goal progress. Be concise.",
				prompt)
			if aiErr == nil {
				sb.WriteString(fmt.Sprintf("   **Assessment**: %s\n", assessment))
			}
		}
		sb.WriteString("\n")
	}

	return mcputil.TextResult(sb.String()), nil
}

// PMTrackDaily captures today's sprint progress from Jira and saves as DailyProgress.
func (h *Handlers) PMTrackDaily(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available."), nil
	}

	result, err := h.SprintService.TrackDaily(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("track daily", err), nil
	}
	return mcputil.TextResult(result), nil
}

// PMBurndown shows sprint burndown chart from daily progress data.
func (h *Handlers) PMBurndown(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available."), nil
	}

	result, err := h.SprintService.Burndown(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("burndown", err), nil
	}
	return mcputil.TextResult(result), nil
}

// PMLearnWorkflow scans board statuses, classifies them, and stores patterns.
func (h *Handlers) PMLearnWorkflow(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.SprintService == nil {
		return mcputil.ErrorResult("Sprint service not available."), nil
	}

	result, err := h.SprintService.LearnWorkflow(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("learn workflow", err), nil
	}
	return mcputil.TextResult(result), nil
}

// PMWorkflowMap displays the learned workflow status classifications for a board.
func (h *Handlers) PMWorkflowMap(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := int(req.GetInt("board_id", 0))
	if boardID <= 0 {
		return mcputil.ErrorResult("board_id is required. Call jira_boards first."), nil
	}
	if h.Memory == nil {
		return mcputil.ErrorResult("Memory store not available."), nil
	}

	patterns, err := h.Memory.GetWorkflowPatterns(ctx, boardID)
	if err != nil {
		return h.Error.Wrap("get workflow patterns", err), nil
	}

	if len(patterns) == 0 {
		return mcputil.TextResult(fmt.Sprintf(
			"**No workflow patterns found** for board **%d**.\n\n"+
				"Run `pm_learn_board` with board_id=%d to auto-discover status classifications.", boardID, boardID)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Workflow Status Map — Board %d\n\n", boardID))
	sb.WriteString(fmt.Sprintf("Total patterns: **%d** | Auto-detected: **%d**\n\n",
		len(patterns), countAuto(patterns)))
	sb.WriteString("| Status | Classification | Pattern | Auto |\n")
	sb.WriteString("|--------|---------------|---------|------|\n")

	for _, p := range patterns {
		auto := "✓"
		if !p.IsAuto {
			auto = ""
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", p.StatusName, p.Classification, p.Pattern, auto))
	}

	sb.WriteString(fmt.Sprintf("\n**Summary**: %s\n", summarizePatterns(patterns)))
	sb.WriteString("\n> Status classifications feed into all sprint analytics. Re-run `pm_learn_board` if board workflow changes.\n")

	return mcputil.TextResult(sb.String()), nil
}

func countAuto(patterns []memory.WorkflowPattern) int {
	var n int
	for _, p := range patterns {
		if p.IsAuto {
			n++
		}
	}
	return n
}

func summarizePatterns(patterns []memory.WorkflowPattern) string {
	var done, blocked, progress, todo int
	for _, p := range patterns {
		switch p.Classification {
		case "done":
			done++
		case "blocked":
			blocked++
		case "progress":
			progress++
		default:
			todo++
		}
	}
	return fmt.Sprintf("%d done · %d blocked · %d progress · %d todo", done, blocked, progress, todo)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
