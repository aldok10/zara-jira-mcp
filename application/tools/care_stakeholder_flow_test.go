package tools_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	jiradom "github.com/aldok10/zara-jira-mcp/domain/jira"
	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
)

// --- Care Handler Tests ---

func TestDailyDelta(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints: []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: []jiradom.Issue{
				{Key: "T-1", Status: "Done", Assignee: "alice", Updated: time.Now()},
				{Key: "T-2", Status: "In Progress", Assignee: "bob", Updated: time.Now()},
				{Key: "T-3", Status: "Blocked", Assignee: "bob", Updated: time.Now()},
			},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.DailyDelta(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Daily Delta") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "CURRENT") {
		t.Errorf("expected current state, got: %s", text)
	}
}

func TestDailyDelta_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.DailyDelta(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestDailyDelta_NoSprint(t *testing.T) {
	h := newTestHandlers(&mockJira{sprints: []jiradom.Sprint{}})
	result, _ := h.DailyDelta(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if !strings.Contains(resultText(result), "No active sprint") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestOverloadCheck(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints: []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: []jiradom.Issue{
				{Key: "T-1", Status: "In Progress", Assignee: "alice"},
				{Key: "T-2", Status: "In Progress", Assignee: "alice"},
				{Key: "T-3", Status: "In Progress", Assignee: "alice"},
				{Key: "T-4", Status: "In Progress", Assignee: "alice"},
				{Key: "T-5", Status: "Done", Assignee: "bob"},
			},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.OverloadCheck(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Team Load") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "alice") {
		t.Errorf("expected alice, got: %s", text)
	}
}

func TestOverloadCheck_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.OverloadCheck(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestOverloadCheck_NoSprint(t *testing.T) {
	h := newTestHandlers(&mockJira{sprints: []jiradom.Sprint{}})
	result, _ := h.OverloadCheck(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if !strings.Contains(resultText(result), "No active sprint") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

// --- Stakeholder Handler Tests ---

func TestSprintScorecard(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints: []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: []jiradom.Issue{
				{Key: "T-1", Status: "Done", Assignee: "alice"},
				{Key: "T-2", Status: "Done", Assignee: "bob"},
				{Key: "T-3", Status: "In Progress", Assignee: "alice"},
			},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.SprintScorecard(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint Scorecard") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "/100") {
		t.Errorf("expected score, got: %s", text)
	}
}

func TestSprintScorecard_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.SprintScorecard(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestSprintScorecard_NoSprint(t *testing.T) {
	h := newTestHandlers(&mockJira{sprints: []jiradom.Sprint{}})
	result, _ := h.SprintScorecard(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if !strings.Contains(resultText(result), "No active sprint") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

// --- Flow Handler Tests ---

func TestSprintComparison_InsufficientData(t *testing.T) {
	h := pmHandlers()
	result, err := h.SprintComparison(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Need at least 2") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestSprintComparison_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.SprintComparison(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

type mockMemoryWithTwoSnapshots struct {
	mockMemory
}

func (m *mockMemoryWithTwoSnapshots) GetSprintSnapshots(_ context.Context, _ int, _ int) ([]memdom.SprintSnapshot, error) {
	return []memdom.SprintSnapshot{
		{SprintName: "Sprint 11", Done: 12, TotalIssues: 15, Velocity: 12, CompletionRate: 80, Blocked: 1, Carryover: 2},
		{SprintName: "Sprint 10", Done: 8, TotalIssues: 14, Velocity: 8, CompletionRate: 57, Blocked: 3, Carryover: 4},
	}, nil
}

func TestSprintComparison_WithData(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithTwoSnapshots{},
	}
	result, err := h.SprintComparison(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint Comparison") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "IMPROVING") && !strings.Contains(text, "STABLE") && !strings.Contains(text, "DECLINING") {
		t.Errorf("expected verdict, got: %s", text)
	}
}

// --- Forecast Handler Tests ---

func TestMonteCarloForecast_WithSnapshots(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "S10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "To Do"}, {Key: "T-2", Status: "In Progress"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithSnapshots{},
	}
	result, err := h.MonteCarloForecast(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Monte Carlo") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "confidence") {
		t.Errorf("expected confidence levels, got: %s", text)
	}
}

func TestForecastSprint_WithData(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done"}, {Key: "T-2", Status: "To Do"}, {Key: "T-3", Status: "In Progress"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithSnapshots{},
	}
	result, err := h.ForecastSprint(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint Forecast") {
		t.Errorf("expected header, got: %s", text)
	}
}

// --- Delivery Confidence with blocked items (RED path) ---

func TestDeliveryConfidenceReport_RED(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints: []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active", Goal: "Ship v2"}},
			sprintIssues: []jiradom.Issue{
				{Key: "T-1", Status: "Blocked"},
				{Key: "T-2", Status: "Blocked"},
				{Key: "T-3", Status: "Blocked"},
				{Key: "T-4", Status: "To Do"},
				{Key: "T-5", Status: "Done"},
			},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.DeliveryConfidenceReport(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "RED") {
		t.Errorf("expected RED status, got: %s", text)
	}
}

// --- VelocityTrend ---

func TestVelocityTrend_InsufficientData(t *testing.T) {
	h := pmHandlers()
	result, err := h.VelocityTrend(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Need at least") && !strings.Contains(text, "No velocity") && !strings.Contains(text, "No sprint snapshot") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestVelocityTrend_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.VelocityTrend(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestVelocityTrend_WithData(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithSnapshots{},
	}
	result, err := h.VelocityTrend(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Velocity") {
		t.Errorf("expected velocity info, got: %s", text)
	}
}

// --- StandupPrep ---

func TestStandupPrep_Basic(t *testing.T) {
	h := pmHandlers()
	result, err := h.StandupPrep(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty standup prep")
	}
}

func TestStandupPrep_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.StandupPrep(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

// --- SprintPlanningSummary ---

func TestSprintPlanningSummary(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active", Goal: "Ship auth"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done"}, {Key: "T-2", Status: "To Do"}},
		},
		AI:     &mockAI{response: "Planning insights"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithSnapshots{},
	}
	result, err := h.SprintPlanningSummary(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty planning summary")
	}
}

func TestSprintPlanningSummary_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.SprintPlanningSummary(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

// --- DetectAntiPatterns ---

func TestDetectAntiPatterns(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: sprintIssues(),
		},
		AI:     &mockAI{response: "Anti-patterns detected"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithSnapshots{},
	}
	result, err := h.DetectAntiPatterns(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty anti-patterns output")
	}
}

func TestDetectAntiPatterns_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.DetectAntiPatterns(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

// --- CoachingAdvice ---

func TestCoachingAdvice(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "Coaching advice: reduce WIP"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.CoachingAdvice(context.Background(), makeReq(map[string]any{
		"topic":     "team_dynamics",
		"situation": "Team is overloaded",
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Coaching") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestCoachingAdvice_MissingTopic(t *testing.T) {
	h := pmHandlers()
	result, _ := h.CoachingAdvice(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing topic")
	}
}

// --- OneOnOnePrep ---

func TestOneOnOnePrep_Basic(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			searchResult: &jiradom.SearchResult{Issues: []jiradom.Issue{{Key: "T-1", Status: "Done", Assignee: "alice"}}},
		},
		AI:     &mockAI{response: "1on1 prep ready"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.OneOnOnePrep(context.Background(), makeReq(map[string]any{"member": "alice"}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty 1on1 prep")
	}
}

func TestOneOnOnePrep_MissingMember(t *testing.T) {
	h := pmHandlers()
	result, _ := h.OneOnOnePrep(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing member")
	}
}

// --- RecordConfidence (additional) ---

func TestRecordConfidence_HighScore(t *testing.T) {
	h := pmHandlers()
	result, _ := h.RecordConfidence(context.Background(), makeReq(map[string]any{
		"sprint_name": "Sprint 10",
		"score":       float64(5),
		"member":      "alice",
	}))
	text := resultText(result)
	if !strings.Contains(text, "confident") {
		t.Errorf("unexpected: %s", text)
	}
}

// --- ReviewExperiments ---

func TestReviewExperiments_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.ReviewExperiments(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty output")
	}
}

// --- NLToJQL ---

func TestNLToJQL_Basic(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			searchResult: &jiradom.SearchResult{Issues: []jiradom.Issue{{Key: "T-1", Status: "Open", Summary: "bug"}}, Total: 1},
		},
		AI:     &mockAI{response: "project = TEST AND status = Open"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.NLToJQL(context.Background(), makeReq(map[string]any{"query": "open bugs in TEST project"}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty output")
	}
}

func TestNLToJQL_MissingQuery(t *testing.T) {
	h := pmHandlers()
	result, _ := h.NLToJQL(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing query")
	}
}

// --- GenerateReleaseNotes ---

func TestGenerateReleaseNotes_Basic(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done", Summary: "Login page", Type: "Story"}},
		},
		AI:     &mockAI{response: "Release notes v1.0"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.GenerateReleaseNotes(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty release notes")
	}
}

func TestGenerateReleaseNotes_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.GenerateReleaseNotes(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

// --- WeeklyDigest ---

func TestWeeklyDigest(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done"}},
		},
		AI:     &mockAI{response: "Weekly summary"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.WeeklyDigest(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty weekly digest")
	}
}

func TestWeeklyDigest_NoBoardID(t *testing.T) {
	// WeeklyDigest uses GetInt with default 0, so no error - just less data
	h := &tools.Handlers{
		Jira:   &mockJira{sprints: []jiradom.Sprint{}},
		AI:     &mockAI{response: "Weekly summary with no sprint data"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.WeeklyDigest(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty output")
	}
}

// --- PMHelp ---

func TestPMHelp(t *testing.T) {
	h := pmHandlers()
	result, err := h.PMHelp(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty help")
	}
}

// --- TeamKnowledgeBase (with data) ---

type mockMemoryWithDecisions struct {
	mockMemory
}

func (m *mockMemoryWithDecisions) GetDecisions(_ context.Context, _ int) ([]memdom.Decision, error) {
	return []memdom.Decision{
		{ID: 1, Title: "Use PostgreSQL", Decision: "Chosen over MongoDB", MadeAt: time.Now()},
	}, nil
}

func TestTeamKnowledgeBase_WithData(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithDecisions{},
	}
	result, err := h.TeamKnowledgeBase(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "PostgreSQL") {
		t.Errorf("expected decisions, got: %s", text)
	}
}

// --- PMTeamRadarHistory ---

func TestPMTeamRadarHistory_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.PMTeamRadarHistory(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "No radar") {
		t.Errorf("unexpected: %s", text)
	}
}
