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

// --- Mocks ---

type mockJira struct {
	searchResult *jiradom.SearchResult
	issue        *jiradom.Issue
	boards       []jiradom.Board
	sprints      []jiradom.Sprint
	sprintIssues []jiradom.Issue
	transitions  []jiradom.Transition
	createdIssue *jiradom.Issue
	err          error
}

func (m *mockJira) SearchIssues(_ context.Context, _ string, _ int, _ int) (*jiradom.SearchResult, error) {
	return m.searchResult, m.err
}
func (m *mockJira) GetIssue(_ context.Context, _ string) (*jiradom.Issue, error) {
	return m.issue, m.err
}
func (m *mockJira) GetBoards(_ context.Context) ([]jiradom.Board, error) {
	return m.boards, m.err
}
func (m *mockJira) GetActiveSprints(_ context.Context, _ int) ([]jiradom.Sprint, error) {
	return m.sprints, m.err
}
func (m *mockJira) GetSprintIssues(_ context.Context, _ int) ([]jiradom.Issue, error) {
	return m.sprintIssues, m.err
}
func (m *mockJira) CreateIssue(_ context.Context, _ *jiradom.CreateIssueInput) (*jiradom.Issue, error) {
	return m.createdIssue, m.err
}
func (m *mockJira) UpdateIssue(_ context.Context, _ *jiradom.UpdateIssueInput) error {
	return m.err
}
func (m *mockJira) AddComment(_ context.Context, _, _ string) error { return m.err }
func (m *mockJira) TransitionIssue(_ context.Context, _, _ string) error {
	return m.err
}
func (m *mockJira) GetTransitions(_ context.Context, _ string) ([]jiradom.Transition, error) {
	return m.transitions, m.err
}
func (m *mockJira) AssignIssue(_ context.Context, _, _ string) error  { return m.err }
func (m *mockJira) DeleteIssue(_ context.Context, _ string) error     { return m.err }
func (m *mockJira) CreateSubtask(_ context.Context, _ string, _ *jiradom.CreateIssueInput) (*jiradom.Issue, error) {
	return &jiradom.Issue{Key: "SUB-1"}, m.err
}
func (m *mockJira) FindUser(_ context.Context, _ string) ([]jiradom.User, error) {
	return nil, m.err
}
func (m *mockJira) SetEpicLink(_ context.Context, _, _ string) error    { return m.err }
func (m *mockJira) RemoveEpicLink(_ context.Context, _ string) error    { return m.err }
func (m *mockJira) GetSprints(_ context.Context, _ int, _ string) ([]jiradom.Sprint, error) {
	return m.sprints, m.err
}
func (m *mockJira) CreateSprint(_ context.Context, _ int, name, _ string) (*jiradom.Sprint, error) {
	return &jiradom.Sprint{ID: 1, Name: name, State: "future"}, m.err
}
func (m *mockJira) StartSprint(_ context.Context, _ int, _, _ string) error { return m.err }
func (m *mockJira) CloseSprint(_ context.Context, _ int) error              { return m.err }
func (m *mockJira) MoveIssuesToSprint(_ context.Context, _ int, _ []string) error {
	return m.err
}
func (m *mockJira) AddLabel(_ context.Context, _, _ string) error { return m.err }
func (m *mockJira) GetProjects(_ context.Context) ([]jiradom.Project, error) {
	return nil, m.err
}
func (m *mockJira) GetProject(_ context.Context, _ string) (*jiradom.ProjectDetail, error) {
	return &jiradom.ProjectDetail{Key: "TEST", Name: "Test Project"}, m.err
}
func (m *mockJira) RawRequest(_ context.Context, _, _ string, _ []byte) ([]byte, int, error) {
	return []byte(`{}`), 200, m.err
}
func (m *mockJira) LinkIssues(_ context.Context, _, _, _ string) error { return m.err }
func (m *mockJira) GetLinkTypes(_ context.Context) ([]jiradom.LinkType, error) {
	return nil, m.err
}
func (m *mockJira) AddWorklog(_ context.Context, _, _, _ string) error { return m.err }
func (m *mockJira) GetWorklogs(_ context.Context, _ string) ([]jiradom.Worklog, error) {
	return nil, m.err
}
func (m *mockJira) AddWatcher(_ context.Context, _, _ string) error { return m.err }
func (m *mockJira) GetWatchers(_ context.Context, _ string) ([]jiradom.User, error) {
	return nil, m.err
}
func (m *mockJira) CreateVersion(_ context.Context, _, _, _ string) (*jiradom.Version, error) {
	return &jiradom.Version{ID: "1", Name: "v1.0"}, m.err
}
func (m *mockJira) GetAttachments(_ context.Context, _ string) ([]jiradom.Attachment, error) {
	return nil, m.err
}
func (m *mockJira) GetComponents(_ context.Context, _ string) ([]jiradom.Component, error) {
	return nil, m.err
}
func (m *mockJira) GetFields(_ context.Context) ([]jiradom.Field, error) { return nil, m.err }
func (m *mockJira) GetVersions(_ context.Context, _ string) ([]jiradom.Version, error) {
	return nil, m.err
}
func (m *mockJira) ReleaseVersion(_ context.Context, _ string) error { return m.err }

type mockAI struct {
	response string
	err      error
}

func (m *mockAI) Complete(_ context.Context, _, _ string) (string, error) {
	return m.response, m.err
}

type mockLark struct {
	err error
}

func (m *mockLark) SendText(_ context.Context, _ string) error       { return m.err }
func (m *mockLark) SendMarkdown(_ context.Context, _, _ string) error { return m.err }

type mockMemory struct{}

func (m *mockMemory) SaveSprintSnapshot(_ context.Context, _ *memdom.SprintSnapshot) error {
	return nil
}
func (m *mockMemory) GetSprintSnapshots(_ context.Context, _ int, _ int) ([]memdom.SprintSnapshot, error) {
	return nil, nil
}
func (m *mockMemory) GetLatestSnapshot(_ context.Context, _ int) (*memdom.SprintSnapshot, error) {
	return nil, nil
}
func (m *mockMemory) SaveRisk(_ context.Context, _ *memdom.Risk) error       { return nil }
func (m *mockMemory) UpdateRisk(_ context.Context, _ *memdom.Risk) error     { return nil }
func (m *mockMemory) GetOpenRisks(_ context.Context) ([]memdom.Risk, error)  { return nil, nil }
func (m *mockMemory) GetAllRisks(_ context.Context, _ int) ([]memdom.Risk, error) {
	return nil, nil
}
func (m *mockMemory) SaveDecision(_ context.Context, _ *memdom.Decision) error { return nil }
func (m *mockMemory) GetDecisions(_ context.Context, _ int) ([]memdom.Decision, error) {
	return nil, nil
}
func (m *mockMemory) SearchDecisions(_ context.Context, _ string) ([]memdom.Decision, error) {
	return nil, nil
}
func (m *mockMemory) SaveBlocker(_ context.Context, _ *memdom.Blocker) error { return nil }
func (m *mockMemory) ResolveBlocker(_ context.Context, _ int64, _ string) error {
	return nil
}
func (m *mockMemory) GetActiveBlockers(_ context.Context) ([]memdom.Blocker, error) {
	return nil, nil
}
func (m *mockMemory) GetBlockerHistory(_ context.Context, _ int) ([]memdom.Blocker, error) {
	return nil, nil
}
func (m *mockMemory) SaveTeamMetric(_ context.Context, _ *memdom.TeamMetric) error { return nil }
func (m *mockMemory) GetTeamMetrics(_ context.Context, _ string, _ int) ([]memdom.TeamMetric, error) {
	return nil, nil
}
func (m *mockMemory) GetTeamOverview(_ context.Context, _ string) ([]memdom.TeamMetric, error) {
	return nil, nil
}
func (m *mockMemory) SaveRetrospective(_ context.Context, _ *memdom.Retrospective) error {
	return nil
}
func (m *mockMemory) GetRetrospectives(_ context.Context, _ int) ([]memdom.Retrospective, error) {
	return nil, nil
}
func (m *mockMemory) SaveActionItem(_ context.Context, _ *memdom.ActionItem) error { return nil }
func (m *mockMemory) GetPendingActionItems(_ context.Context) ([]memdom.ActionItem, error) {
	return nil, nil
}
func (m *mockMemory) CompleteActionItem(_ context.Context, _ int64) error { return nil }
func (m *mockMemory) SaveHealthScore(_ context.Context, _ *memdom.HealthScore) error {
	return nil
}
func (m *mockMemory) GetHealthScores(_ context.Context, _ int, _ int) ([]memdom.HealthScore, error) {
	return nil, nil
}
func (m *mockMemory) SaveDependency(_ context.Context, _ *memdom.Dependency) error { return nil }
func (m *mockMemory) ResolveDependency(_ context.Context, _ int64) error            { return nil }
func (m *mockMemory) GetDependenciesForIssue(_ context.Context, _ string) ([]memdom.Dependency, error) {
	return nil, nil
}
func (m *mockMemory) GetOpenDependencies(_ context.Context) ([]memdom.Dependency, error) {
	return nil, nil
}
func (m *mockMemory) SaveMeetingNote(_ context.Context, _ *memdom.MeetingNote) error { return nil }
func (m *mockMemory) GetMeetingNotes(_ context.Context, _ string, _ int) ([]memdom.MeetingNote, error) {
	return nil, nil
}
func (m *mockMemory) SaveDailyProgress(_ context.Context, _ *memdom.DailyProgress) error {
	return nil
}
func (m *mockMemory) GetDailyProgress(_ context.Context, _ int, _ string) ([]memdom.DailyProgress, error) {
	return nil, nil
}
func (m *mockMemory) SaveSprintGoal(_ context.Context, _ *memdom.SprintGoal) error { return nil }
func (m *mockMemory) UpdateSprintGoal(_ context.Context, _ *memdom.SprintGoal) error {
	return nil
}
func (m *mockMemory) GetActiveGoals(_ context.Context, _ int) ([]memdom.SprintGoal, error) {
	return nil, nil
}
func (m *mockMemory) GetGoalHistory(_ context.Context, _ int, _ int) ([]memdom.SprintGoal, error) {
	return nil, nil
}
func (m *mockMemory) SaveDoDItem(_ context.Context, _ *memdom.DoDItem) error { return nil }
func (m *mockMemory) GetDoD(_ context.Context, _ string) ([]memdom.DoDItem, error) {
	return nil, nil
}
func (m *mockMemory) DeleteDoDItem(_ context.Context, _ int64) error { return nil }
func (m *mockMemory) SaveEscalation(_ context.Context, _ *memdom.Escalation) error {
	return nil
}
func (m *mockMemory) GetRecentEscalations(_ context.Context, _ int) ([]memdom.Escalation, error) {
	return nil, nil
}
func (m *mockMemory) AcknowledgeEscalation(_ context.Context, _ int64) error { return nil }
func (m *mockMemory) SaveTeamPulse(_ context.Context, _ *memdom.TeamPulse) error {
	return nil
}
func (m *mockMemory) GetTeamPulseHistory(_ context.Context, _ int) ([]memdom.TeamPulse, error) {
	return nil, nil
}
func (m *mockMemory) SaveMeetingEffectiveness(_ context.Context, _ *memdom.MeetingEffectiveness) error {
	return nil
}
func (m *mockMemory) GetMeetingEffectivenessHistory(_ context.Context, _ string, _ int) ([]memdom.MeetingEffectiveness, error) {
	return nil, nil
}
func (m *mockMemory) SaveTeamRadar(_ context.Context, _ *memdom.TeamRadar) error { return nil }
func (m *mockMemory) GetTeamRadarHistory(_ context.Context, _ int) ([]memdom.TeamRadar, error) {
	return nil, nil
}
func (m *mockMemory) DB() memdom.RawDB { return nil }
func (m *mockMemory) DeleteOKRSignal(_ context.Context, _ int64) error { return nil }
func (m *mockMemory) SaveOKRSignal(_ context.Context, _ *memdom.OKRSignal) error { return nil }
func (m *mockMemory) UpdateOKRSignalProgress(_ context.Context, _ int64, _, _ float64) error {
	return nil
}
func (m *mockMemory) GetOKRSignals(_ context.Context) ([]memdom.OKRSignal, error) {
	return nil, nil
}

// --- Helpers ---

func newTestHandlers(jira *mockJira) *tools.Handlers {
	return &tools.Handlers{
		Jira:   jira,
		AI:     &mockAI{response: "test analysis"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
}

func makeReq(args map[string]any) mcp.CallToolRequest {
	req := mcp.CallToolRequest{}
	req.Params.Arguments = args
	return req
}

func resultText(r *mcp.CallToolResult) string {
	if len(r.Content) == 0 {
		return ""
	}
	if tc, ok := r.Content[0].(mcp.TextContent); ok {
		return tc.Text
	}
	return ""
}

// --- Tests ---

func TestSearchIssues(t *testing.T) {
	jiraMock := &mockJira{
		searchResult: &jiradom.SearchResult{
			Issues: []jiradom.Issue{
				{Key: "TEST-1", Summary: "Fix bug", Status: "Open", Priority: "High", Type: "Bug", Assignee: "dev1"},
				{Key: "TEST-2", Summary: "Add feature", Status: "In Progress", Priority: "Medium", Type: "Story", Assignee: "dev2"},
			},
			Total:   2,
			StartAt: 0,
			HasMore: false,
		},
	}
	h := newTestHandlers(jiraMock)

	result, err := h.SearchIssues(context.Background(), makeReq(map[string]any{"jql": "project = TEST"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("got error result: %s", resultText(result))
	}

	text := resultText(result)
	if !strings.Contains(text, "TEST-1") {
		t.Errorf("expected TEST-1 in output, got: %s", text)
	}
	if !strings.Contains(text, "Found 2 issues") {
		t.Errorf("expected 'Found 2 issues' in output, got: %s", text)
	}
}

func TestSearchIssues_MissingJQL(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.SearchIssues(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing jql")
	}
}

func TestGetIssue(t *testing.T) {
	jiraMock := &mockJira{
		issue: &jiradom.Issue{
			Key:      "TEST-1",
			Summary:  "Fix critical bug",
			Status:   "In Progress",
			Priority: "High",
			Type:     "Bug",
			Assignee: "dev1",
			Updated:  time.Now(),
		},
	}
	h := newTestHandlers(jiraMock)

	result, err := h.GetIssue(context.Background(), makeReq(map[string]any{"key": "TEST-1"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("got error result: %s", resultText(result))
	}

	text := resultText(result)
	if !strings.Contains(text, "TEST-1") {
		t.Errorf("expected TEST-1 in output, got: %s", text)
	}
	if !strings.Contains(text, "Fix critical bug") {
		t.Errorf("expected summary in output, got: %s", text)
	}
}

func TestGetIssue_MissingKey(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.GetIssue(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing key")
	}
}

func TestCreateIssue(t *testing.T) {
	jiraMock := &mockJira{
		createdIssue: &jiradom.Issue{Key: "TEST-3", Summary: "New task"},
	}
	h := newTestHandlers(jiraMock)

	result, err := h.CreateIssue(context.Background(), makeReq(map[string]any{
		"project": "TEST",
		"summary": "New task",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("got error result: %s", resultText(result))
	}

	text := resultText(result)
	if !strings.Contains(text, "TEST-3") {
		t.Errorf("expected TEST-3 in output, got: %s", text)
	}
}

func TestCreateIssue_MissingProject(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.CreateIssue(context.Background(), makeReq(map[string]any{"summary": "x"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error for missing project")
	}
}

func TestUpdateIssue(t *testing.T) {
	h := newTestHandlers(&mockJira{})

	result, err := h.UpdateIssue(context.Background(), makeReq(map[string]any{
		"key":      "TEST-1",
		"summary":  "Updated summary",
		"priority": "Low",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("got error result: %s", resultText(result))
	}

	text := resultText(result)
	if !strings.Contains(text, "TEST-1 updated successfully") {
		t.Errorf("expected success message, got: %s", text)
	}
}

func TestUpdateIssue_MissingKey(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.UpdateIssue(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error for missing key")
	}
}

func TestHealth(t *testing.T) {
	h := newTestHandlers(&mockJira{})

	result, err := h.Health(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatal("unexpected error result")
	}

	text := resultText(result)
	if !strings.Contains(text, "v0.3.0") {
		t.Errorf("expected version in output, got: %s", text)
	}
	if !strings.Contains(text, "ok") {
		t.Errorf("expected 'ok' in output, got: %s", text)
	}
}

func TestMyIssues(t *testing.T) {
	jiraMock := &mockJira{
		searchResult: &jiradom.SearchResult{
			Issues: []jiradom.Issue{
				{Key: "TEST-5", Summary: "My task", Status: "In Progress", Priority: "High"},
			},
			Total: 1,
		},
	}
	h := newTestHandlers(jiraMock)

	result, err := h.MyIssues(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("got error result: %s", resultText(result))
	}

	text := resultText(result)
	if !strings.Contains(text, "TEST-5") {
		t.Errorf("expected TEST-5 in output, got: %s", text)
	}
	if !strings.Contains(text, "Your open issues (1)") {
		t.Errorf("expected issue count in output, got: %s", text)
	}
}

func TestOverdue(t *testing.T) {
	jiraMock := &mockJira{
		searchResult: &jiradom.SearchResult{
			Issues: []jiradom.Issue{
				{Key: "TEST-10", Summary: "Stale ticket", Status: "Open", Assignee: "dev1", Updated: time.Now().Add(-30 * 24 * time.Hour)},
			},
			Total: 1,
		},
	}
	h := newTestHandlers(jiraMock)

	result, err := h.Overdue(context.Background(), makeReq(map[string]any{"days": float64(14)}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.IsError {
		t.Fatalf("got error result: %s", resultText(result))
	}

	text := resultText(result)
	if !strings.Contains(text, "TEST-10") {
		t.Errorf("expected TEST-10 in output, got: %s", text)
	}
	if !strings.Contains(text, "Stale issues") {
		t.Errorf("expected 'Stale issues' in output, got: %s", text)
	}
}
