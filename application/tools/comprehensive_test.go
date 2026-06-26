package tools_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	jiradom "github.com/aldok10/zara-jira-mcp/domain/jira"
	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// --- Enriched mocks ---

type richMemory struct {
	mockMemory
}

func (m *richMemory) GetSprintSnapshots(_ context.Context, _ int, _ int) ([]memdom.SprintSnapshot, error) {
	return []memdom.SprintSnapshot{
		{SprintName: "Sprint 5", Done: 12, InProgress: 2, Todo: 1, Blocked: 0, Velocity: 24, CompletionRate: 80, TotalIssues: 15, Carryover: 1},
		{SprintName: "Sprint 4", Done: 10, InProgress: 3, Todo: 2, Blocked: 1, Velocity: 20, CompletionRate: 66, TotalIssues: 16, Carryover: 2},
		{SprintName: "Sprint 3", Done: 11, InProgress: 2, Todo: 1, Blocked: 0, Velocity: 22, CompletionRate: 78, TotalIssues: 14, Carryover: 1},
		{SprintName: "Sprint 2", Done: 9, InProgress: 4, Todo: 2, Blocked: 1, Velocity: 18, CompletionRate: 56, TotalIssues: 16, Carryover: 3},
		{SprintName: "Sprint 1", Done: 8, InProgress: 3, Todo: 3, Blocked: 2, Velocity: 16, CompletionRate: 50, TotalIssues: 16, Carryover: 4},
	}, nil
}

func (m *richMemory) GetLatestSnapshot(_ context.Context, _ int) (*memdom.SprintSnapshot, error) {
	return &memdom.SprintSnapshot{SprintName: "Sprint 5", TotalIssues: 15, Done: 12, Velocity: 24, CompletionRate: 80}, nil
}

func (m *richMemory) GetOpenRisks(_ context.Context) ([]memdom.Risk, error) {
	return []memdom.Risk{
		{ID: 1, Title: "External API scope change", Severity: "high", Status: "open", Owner: "alice", IdentifiedAt: time.Now().Add(-5 * 24 * time.Hour), Mitigation: "remove from sprint"},
	}, nil
}

func (m *richMemory) GetActiveBlockers(_ context.Context) ([]memdom.Blocker, error) {
	return []memdom.Blocker{
		{ID: 1, IssueKey: "T-4", Description: "Waiting for infra team deploy", BlockedSince: time.Now().Add(-4 * 24 * time.Hour), Owner: "bob"},
	}, nil
}

func (m *richMemory) GetBlockerHistory(_ context.Context, _ int) ([]memdom.Blocker, error) {
	resolved := time.Now().Add(-1 * 24 * time.Hour)
	return []memdom.Blocker{
		{ID: 1, IssueKey: "T-4", Description: "Waiting for infra team deploy", BlockedSince: time.Now().Add(-4 * 24 * time.Hour), Owner: "bob"},
		{ID: 2, IssueKey: "T-2", Description: "Old blocker", BlockedSince: time.Now().Add(-10 * 24 * time.Hour), ResolvedAt: &resolved, DaysBlocked: 3, Resolution: "Fixed"},
	}, nil
}

func (m *richMemory) GetDecisions(_ context.Context, _ int) ([]memdom.Decision, error) {
	return []memdom.Decision{
		{ID: 1, Title: "Use PostgreSQL", Decision: "PostgreSQL over MongoDB", Rationale: "ACID", MadeAt: time.Now().Add(-3 * 24 * time.Hour)},
	}, nil
}

func (m *richMemory) GetPendingActionItems(_ context.Context) ([]memdom.ActionItem, error) {
	return []memdom.ActionItem{
		{ID: 1, Description: "Set up monitoring dashboard", Owner: "charlie", Status: "pending"},
	}, nil
}

func (m *richMemory) GetRetrospectives(_ context.Context, _ int) ([]memdom.Retrospective, error) {
	return []memdom.Retrospective{
		{ID: 1, SprintName: "Sprint 4", WentWell: "Good collaboration", Improvements: "communication delays", Status: "closed"},
		{ID: 2, SprintName: "Sprint 5", WentWell: "Fast delivery", Improvements: "testing coverage", Status: "open"},
	}, nil
}

func (m *richMemory) GetOpenDependencies(_ context.Context) ([]memdom.Dependency, error) {
	return []memdom.Dependency{
		{ID: 1, FromIssueKey: "T-3", ToIssueKey: "EXT-1", DependencyType: "blocked_by", Description: "Need API from platform team", CreatedAt: time.Now().Add(-6 * 24 * time.Hour)},
	}, nil
}

func (m *richMemory) GetTeamMetrics(_ context.Context, _ string, _ int) ([]memdom.TeamMetric, error) {
	return []memdom.TeamMetric{
		{MemberName: "alice", SprintName: "Sprint 5", IssuesAssigned: 5, IssuesDone: 4, BlockerCount: 0, CarryoverCount: 1},
	}, nil
}

func (m *richMemory) GetHealthScores(_ context.Context, _ int, _ int) ([]memdom.HealthScore, error) {
	return []memdom.HealthScore{
		{SprintName: "Sprint 5", BoardID: 1, OverallScore: 75, VelocityScore: 20, BlockerScore: 18, ScopeScore: 20, TeamScore: 17},
	}, nil
}

func (m *richMemory) GetActiveGoals(_ context.Context, _ int) ([]memdom.SprintGoal, error) {
	return []memdom.SprintGoal{
		{ID: 1, SprintName: "Sprint 5", BoardID: 1, Goal: "Ship auth module", Status: "active"},
	}, nil
}

func (m *richMemory) GetGoalHistory(_ context.Context, _ int, _ int) ([]memdom.SprintGoal, error) {
	return []memdom.SprintGoal{
		{ID: 1, SprintName: "Sprint 5", Goal: "Ship auth module", Status: "active"},
	}, nil
}

func (m *richMemory) GetDailyProgress(_ context.Context, _ int, _ string) ([]memdom.DailyProgress, error) {
	return []memdom.DailyProgress{
		{SprintName: "Sprint 5", Date: time.Now().Add(-2 * 24 * time.Hour), TotalIssues: 15, Done: 10, InProgress: 3, Todo: 1, Blocked: 1},
		{SprintName: "Sprint 5", Date: time.Now().Add(-1 * 24 * time.Hour), TotalIssues: 15, Done: 11, InProgress: 2, Todo: 1, Blocked: 1},
		{SprintName: "Sprint 5", Date: time.Now(), TotalIssues: 15, Done: 12, InProgress: 2, Todo: 1, Blocked: 0},
	}, nil
}

type richJira struct {
	mockJira
}

func (m *richJira) GetActiveSprints(_ context.Context, _ int) ([]jiradom.Sprint, error) {
	return []jiradom.Sprint{{ID: 1, Name: "Sprint 5", State: "active", Goal: "Ship auth module"}}, nil
}

func (m *richJira) GetSprintIssues(_ context.Context, _ int) ([]jiradom.Issue, error) {
	return []jiradom.Issue{
		{Key: "T-1", Summary: "Auth login", Status: "Done", Type: "Story", Assignee: "alice", Created: time.Now().Add(-10 * 24 * time.Hour), Updated: time.Now().Add(-1 * 24 * time.Hour)},
		{Key: "T-2", Summary: "Fix crash", Status: "Done", Type: "Bug", Assignee: "bob", Created: time.Now().Add(-8 * 24 * time.Hour), Updated: time.Now().Add(-1 * 24 * time.Hour)},
		{Key: "T-3", Summary: "API endpoint", Status: "In Progress", Type: "Task", Assignee: "alice", Created: time.Now().Add(-5 * 24 * time.Hour), Updated: time.Now()},
		{Key: "T-4", Summary: "Blocked thing", Status: "Blocked", Type: "Task", Assignee: "charlie", Created: time.Now().Add(-4 * 24 * time.Hour), Updated: time.Now()},
		{Key: "T-5", Summary: "Todo item", Status: "To Do", Type: "Story", Assignee: "bob", Created: time.Now().Add(-3 * 24 * time.Hour), Updated: time.Now()},
	}, nil
}

func (m *richJira) SearchIssues(_ context.Context, _ string, _ int, _ int) (*jiradom.SearchResult, error) {
	return &jiradom.SearchResult{
		Issues: []jiradom.Issue{
			{Key: "STALE-1", Summary: "Old ticket", Status: "Open", Type: "Bug", Assignee: "dev1", Updated: time.Now().Add(-100 * 24 * time.Hour)},
		},
		Total: 1,
	}, nil
}

func (m *richJira) GetProjects(_ context.Context) ([]jiradom.Project, error) {
	return []jiradom.Project{
		{Key: "PROJ", Name: "Project Alpha", Lead: "Alice", Type: "software"},
	}, nil
}

// --- Test helpers ---

func richHandlers() *tools.Handlers {
	return &tools.Handlers{
		Jira:   &richJira{},
		AI:     &mockAI{response: "AI analysis"},
		Lark:   &mockLark{},
		Memory: &richMemory{},
	}
}

func assertOK(t *testing.T, result *mcp.CallToolResult, err error) string {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("got error result: %s", resultText(result))
	}
	text := resultText(result)
	if text == "" {
		t.Fatal("expected non-empty result")
	}
	return text
}

// --- Reporting Tests ---

func TestCompReportToPO(t *testing.T) {
	h := richHandlers()
	result, err := h.ReportToPO(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

func TestCompEscalationBrief(t *testing.T) {
	h := richHandlers()
	result, err := h.EscalationBrief(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

func TestCompEscalationBrief_NoBlockers(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &richJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.EscalationBrief(context.Background(), makeReq(map[string]any{}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "No active impediments") {
		t.Errorf("expected clean message, got: %s", text)
	}
}

func TestCompCrossTeamDependencyReport(t *testing.T) {
	h := richHandlers()
	result, err := h.CrossTeamDependencyReport(context.Background(), makeReq(nil))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "CROSS-TEAM") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "T-3") {
		t.Errorf("expected dependency issue key, got: %s", text)
	}
}

func TestCompCrossTeamDependencyReport_Empty(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &richJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.CrossTeamDependencyReport(context.Background(), makeReq(nil))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "No cross-team") {
		t.Errorf("expected empty message, got: %s", text)
	}
}

func TestCompDeliveryConfidenceReport(t *testing.T) {
	h := richHandlers()
	result, err := h.DeliveryConfidenceReport(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

func TestCompResourcePlanningReport(t *testing.T) {
	h := richHandlers()
	result, err := h.ResourcePlanningReport(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "RESOURCE") {
		t.Errorf("expected header, got: %s", text)
	}
}

// --- SM Leverage Tests ---

func TestCompTeamMaturityAssess(t *testing.T) {
	h := richHandlers()
	result, err := h.TeamMaturityAssess(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Maturity") {
		t.Errorf("expected maturity header, got: %s", text)
	}
}

func TestCompImprovementVelocity(t *testing.T) {
	h := richHandlers()
	result, err := h.ImprovementVelocity(context.Background(), makeReq(nil))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Improvement Velocity") {
		t.Errorf("expected header, got: %s", text)
	}
}

func TestCompImprovementVelocity_NoRetros(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &richJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.ImprovementVelocity(context.Background(), makeReq(nil))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "No retrospectives") {
		t.Errorf("expected empty message, got: %s", text)
	}
}

func TestCompMeetingROI(t *testing.T) {
	h := richHandlers()
	result, err := h.MeetingROI(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Meeting ROI") {
		t.Errorf("expected header, got: %s", text)
	}
}

func TestCompSprintCommitmentAdvisor(t *testing.T) {
	h := richHandlers()
	result, err := h.SprintCommitmentAdvisor(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

func TestCompDysfunctionDetector(t *testing.T) {
	h := richHandlers()
	result, err := h.DysfunctionDetector(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

// --- WhatNext & Related ---

func TestCompWhatNext(t *testing.T) {
	h := richHandlers()
	result, err := h.WhatNext(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

func TestCompOneOnOnePrep(t *testing.T) {
	h := richHandlers()
	result, err := h.OneOnOnePrep(context.Background(), makeReq(map[string]any{"member": "alice"}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "alice") {
		t.Errorf("expected member name, got: %s", text)
	}
}

func TestCompSprintNarrative(t *testing.T) {
	h := richHandlers()
	result, err := h.SprintNarrative(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

// --- PM Intel Tests ---

func TestCompPMRecommendations(t *testing.T) {
	h := richHandlers()
	result, err := h.PMRecommendations(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

func TestCompVelocityTrend(t *testing.T) {
	h := richHandlers()
	result, err := h.VelocityTrend(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Velocity Trend") {
		t.Errorf("expected header, got: %s", text)
	}
}

func TestCompStandupPrep(t *testing.T) {
	h := richHandlers()
	result, err := h.StandupPrep(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

func TestCompSprintRetroAnalysis(t *testing.T) {
	h := richHandlers()
	result, err := h.SprintRetroAnalysis(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

// --- Flow Tests ---

func TestCompFlowMetrics(t *testing.T) {
	h := richHandlers()
	result, err := h.FlowMetrics(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "WIP") {
		t.Errorf("expected WIP, got: %s", text)
	}
}

func TestCompSprintComparison(t *testing.T) {
	h := richHandlers()
	result, err := h.SprintComparison(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Sprint Comparison") {
		t.Errorf("expected header, got: %s", text)
	}
}

func TestCompCeremonyFacilitator_Planning(t *testing.T) {
	h := richHandlers()
	result, err := h.CeremonyFacilitator(context.Background(), makeReq(map[string]any{"ceremony": "planning", "board_id": float64(1)}))
	assertOK(t, result, err)
}

// --- Deep Handlers Tests ---

func TestCompTrackDailyProgress(t *testing.T) {
	h := richHandlers()
	result, err := h.TrackDailyProgress(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Sprint 5") {
		t.Errorf("expected sprint name, got: %s", text)
	}
}

func TestCompGetBurndown(t *testing.T) {
	h := richHandlers()
	result, err := h.GetBurndown(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

func TestCompSetSprintGoal(t *testing.T) {
	h := richHandlers()
	result, err := h.SetSprintGoal(context.Background(), makeReq(map[string]any{
		"board_id": float64(1),
		"goal":     "Deliver auth flow",
	}))
	assertOK(t, result, err)
}

func TestCompManageDoD_Add(t *testing.T) {
	h := richHandlers()
	result, err := h.ManageDoD(context.Background(), makeReq(map[string]any{
		"action":   "add",
		"item":     "All tests pass",
		"category": "testing",
	}))
	assertOK(t, result, err)
}

func TestCompPMDashboard(t *testing.T) {
	h := richHandlers()
	result, err := h.PMDashboard(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Overall:") {
		t.Errorf("expected dashboard header, got: %s", text)
	}
}

func TestCompGenerateReleaseNotes(t *testing.T) {
	h := richHandlers()
	result, err := h.GenerateReleaseNotes(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

func TestCompSprintHealthScore(t *testing.T) {
	h := richHandlers()
	result, err := h.SprintHealthScore(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Sprint Health:") {
		t.Errorf("expected health header, got: %s", text)
	}
}

func TestCompAutoDetectRisks(t *testing.T) {
	h := richHandlers()
	result, err := h.AutoDetectRisks(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

// --- Memory Tests ---

func TestCompSnapshotSprint(t *testing.T) {
	h := richHandlers()
	result, err := h.SnapshotSprint(context.Background(), makeReq(map[string]any{"board_id": float64(1), "velocity": float64(24)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Sprint 5") {
		t.Errorf("expected sprint name, got: %s", text)
	}
}

func TestCompRecordRisk(t *testing.T) {
	h := richHandlers()
	result, err := h.RecordRisk(context.Background(), makeReq(map[string]any{
		"title": "Dependency risk", "severity": "high", "owner": "alice",
	}))
	assertOK(t, result, err)
}

func TestCompRecordDecision(t *testing.T) {
	h := richHandlers()
	result, err := h.RecordDecision(context.Background(), makeReq(map[string]any{
		"title": "Go with REST", "decision": "REST over gRPC", "rationale": "simpler",
	}))
	assertOK(t, result, err)
}

func TestCompRecordBlocker(t *testing.T) {
	h := richHandlers()
	result, err := h.RecordBlocker(context.Background(), makeReq(map[string]any{
		"description": "Waiting for access", "owner": "pm",
	}))
	assertOK(t, result, err)
}

func TestCompGetBlockers_Active(t *testing.T) {
	h := richHandlers()
	result, err := h.GetBlockers(context.Background(), makeReq(map[string]any{}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Active Blockers") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "T-4") {
		t.Errorf("expected blocker issue key, got: %s", text)
	}
}

func TestCompRecordTeamMetric(t *testing.T) {
	h := richHandlers()
	result, err := h.RecordTeamMetric(context.Background(), makeReq(map[string]any{
		"member_name": "alice", "sprint_name": "Sprint 5", "issues_assigned": float64(5), "issues_done": float64(4),
	}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "alice") {
		t.Errorf("expected member name, got: %s", text)
	}
}

// --- Forecast Tests ---

func TestCompForecastSprint(t *testing.T) {
	h := richHandlers()
	result, err := h.ForecastSprint(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Sprint Forecast") || !strings.Contains(text, "confidence") {
		t.Errorf("expected forecast data, got: %s", text)
	}
}

func TestCompScopeCreep_WithSnapshot(t *testing.T) {
	h := richHandlers()
	result, err := h.ScopeCreep(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	assertOK(t, result, err)
}

func TestCompBacklogGroom(t *testing.T) {
	h := richHandlers()
	result, err := h.BacklogGroom(context.Background(), makeReq(map[string]any{"days": float64(90)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "stale") {
		t.Errorf("expected stale items, got: %s", text)
	}
}

func TestCompNLToJQL(t *testing.T) {
	h := richHandlers()
	result, err := h.NLToJQL(context.Background(), makeReq(map[string]any{"query": "my open bugs"}))
	assertOK(t, result, err)
}

// --- Portfolio Tests ---

func TestCompPortfolioOverview(t *testing.T) {
	h := richHandlers()
	result, err := h.PortfolioOverview(context.Background(), makeReq(nil))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Portfolio Overview") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "PROJ") {
		t.Errorf("expected project, got: %s", text)
	}
}

func TestCompPortfolioBlockers_WithData(t *testing.T) {
	h := richHandlers()
	result, err := h.PortfolioBlockers(context.Background(), makeReq(nil))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Active Blockers") {
		t.Errorf("expected blockers section, got: %s", text)
	}
}

func TestCompPortfolioWorkload(t *testing.T) {
	h := richHandlers()
	result, err := h.PortfolioWorkload(context.Background(), makeReq(nil))
	assertOK(t, result, err)
}

// --- Monte Carlo with data ---

func TestCompMonteCarloForecast_WithData(t *testing.T) {
	h := richHandlers()
	result, err := h.MonteCarloForecast(context.Background(), makeReq(map[string]any{"board_id": float64(1), "remaining_items": float64(10)}))
	text := assertOK(t, result, err)
	if !strings.Contains(text, "Monte Carlo") || !strings.Contains(text, "%") {
		t.Errorf("expected monte carlo output, got: %s", text)
	}
}
