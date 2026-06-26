package tools_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	jiradom "github.com/aldok10/zara-jira-mcp/domain/jira"
	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
)

// --- Mock overrides for this file ---

type mockJiraWithProjects struct {
	mockJira
}

func (m *mockJiraWithProjects) GetProjects(_ context.Context) ([]jiradom.Project, error) {
	return []jiradom.Project{
		{Key: "PROJ", Name: "Project Alpha", Lead: "Alice", Type: "software"},
		{Key: "BETA", Name: "Project Beta", Lead: "Bob", Type: "software"},
	}, nil
}

type mockMemoryWithSnapshots struct {
	mockMemory
}

func (m *mockMemoryWithSnapshots) GetSprintSnapshots(_ context.Context, _ int, _ int) ([]memdom.SprintSnapshot, error) {
	return []memdom.SprintSnapshot{
		{Done: 10, Velocity: 10},
		{Done: 8, Velocity: 8},
		{Done: 12, Velocity: 12},
		{Done: 9, Velocity: 9},
		{Done: 11, Velocity: 11},
	}, nil
}

// --- Epic Tests ---

func TestEpicIssues(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			searchResult: &jiradom.SearchResult{
				Issues: []jiradom.Issue{
					{Key: "PROJ-1", Summary: "Child 1", Status: "Done", Type: "Story", Assignee: "dev1"},
					{Key: "PROJ-2", Summary: "Child 2", Status: "In Progress", Type: "Task", Assignee: "dev2"},
				},
				Total: 2,
			},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.EpicIssues(context.Background(), makeReq(map[string]any{"epic_key": "PROJ-100"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "PROJ-100 has 2 issues") {
		t.Errorf("expected epic summary, got: %s", text)
	}
}

func TestEpicAdd(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.EpicAdd(context.Background(), makeReq(map[string]any{
		"issue_keys": "PROJ-1,PROJ-2",
		"epic_key":   "PROJ-100",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Added 2 issue(s) to epic PROJ-100") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestEpicRemove(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.EpicRemove(context.Background(), makeReq(map[string]any{
		"issue_keys": "PROJ-1,PROJ-2,PROJ-3",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Removed 3 issue(s)") {
		t.Errorf("unexpected: %s", text)
	}
}

// --- Sprint Tests ---

func TestListSprints(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints: []jiradom.Sprint{
				{ID: 1, Name: "Sprint 1", State: "active", Goal: "Ship v1"},
				{ID: 2, Name: "Sprint 2", State: "future"},
			},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.ListSprints(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint 1") || !strings.Contains(text, "Sprint 2") {
		t.Errorf("expected both sprints, got: %s", text)
	}
}

func TestCreateSprintTool(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.CreateSprintTool(context.Background(), makeReq(map[string]any{
		"board_id": float64(1),
		"name":     "Sprint 99",
		"goal":     "Test goal",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint 99") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestMoveIssuesToSprintTool(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.MoveIssuesToSprintTool(context.Background(), makeReq(map[string]any{
		"sprint_id":  float64(5),
		"issue_keys": "PROJ-1,PROJ-2",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Moved 2 issue(s) to sprint 5") {
		t.Errorf("unexpected: %s", text)
	}
}

// --- Bulk Ops Tests ---

func TestBulkTransition(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.BulkTransition(context.Background(), makeReq(map[string]any{
		"issue_keys":    "T-1,T-2,T-3",
		"transition_id": "21",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "3 succeeded") {
		t.Errorf("expected 3 succeeded, got: %s", text)
	}
}

func TestBulkAssign(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.BulkAssign(context.Background(), makeReq(map[string]any{
		"issue_keys":  "T-1,T-2",
		"assignee_id": "user123",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "2 succeeded") {
		t.Errorf("expected 2 succeeded, got: %s", text)
	}
}

// --- Recipe Tests ---

func TestRecipeStartWork(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			issue:       &jiradom.Issue{Key: "PROJ-123", Summary: "Fix login", Status: "To Do", Assignee: "John"},
			transitions: []jiradom.Transition{{ID: "21", Name: "In Progress"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.RecipeStartWork(context.Background(), makeReq(map[string]any{"key": "PROJ-123"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "PROJ-123") {
		t.Errorf("expected issue key, got: %s", text)
	}
	if !strings.Contains(text, "feature/") {
		t.Errorf("expected branch suggestion, got: %s", text)
	}
}

func TestRecipeDone(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			transitions: []jiradom.Transition{{ID: "31", Name: "Done"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.RecipeDone(context.Background(), makeReq(map[string]any{"key": "PROJ-50"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Completed: PROJ-50") {
		t.Errorf("unexpected: %s", text)
	}
	if !strings.Contains(text, "Transitioned to Done") {
		t.Errorf("expected transition action, got: %s", text)
	}
}

func TestRecipeBlock(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.RecipeBlock(context.Background(), makeReq(map[string]any{
		"key":    "PROJ-7",
		"reason": "Waiting for API key",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Blocked: PROJ-7") {
		t.Errorf("unexpected: %s", text)
	}
	if !strings.Contains(text, "Waiting for API key") {
		t.Errorf("expected reason, got: %s", text)
	}
}

// --- Forecast Tests ---

func TestForecastSprint(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 5", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "To Do"}, {Key: "T-2", Status: "In Progress"}, {Key: "T-3", Status: "Done"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithSnapshots{},
	}
	result, err := h.ForecastSprint(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint Forecast") {
		t.Errorf("expected forecast header, got: %s", text)
	}
	if !strings.Contains(text, "confidence") {
		t.Errorf("expected confidence data, got: %s", text)
	}
}

func TestScopeCreep(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 5", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1"}, {Key: "T-2"}, {Key: "T-3"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.ScopeCreep(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	// Without a snapshot, it should say no baseline
	if !strings.Contains(text, "No baseline snapshot") {
		t.Errorf("expected no baseline message, got: %s", text)
	}
}

func TestBacklogGroom(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			searchResult: &jiradom.SearchResult{
				Issues: []jiradom.Issue{
					{Key: "OLD-1", Summary: "Ancient ticket", Type: "Bug", Assignee: "nobody"},
				},
				Total: 1,
			},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.BacklogGroom(context.Background(), makeReq(map[string]any{"days": float64(90)}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "stale items") {
		t.Errorf("expected stale items header, got: %s", text)
	}
	if !strings.Contains(text, "OLD-1") {
		t.Errorf("expected issue key, got: %s", text)
	}
}

// --- GitHub Tests ---

func TestIssueFromBranch(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			issue: &jiradom.Issue{Key: "PROJ-123", Summary: "Fix login", Status: "To Do", Assignee: "John"},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.IssueFromBranch(context.Background(), makeReq(map[string]any{"branch": "feature/proj-123-fix-login"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "PROJ-123") {
		t.Errorf("expected issue key, got: %s", text)
	}
	if !strings.Contains(text, "Fix login") {
		t.Errorf("expected summary, got: %s", text)
	}
}

func TestSmartCommit(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			transitions: []jiradom.Transition{{ID: "31", Name: "Done"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.SmartCommit(context.Background(), makeReq(map[string]any{
		"message": "PROJ-45 fix auth #done #time 2h",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "PROJ-45") {
		t.Errorf("expected issue key, got: %s", text)
	}
	if !strings.Contains(text, "Transitioned to Done") {
		t.Errorf("expected transition action, got: %s", text)
	}
	if !strings.Contains(text, "Logged 2h") {
		t.Errorf("expected time log, got: %s", text)
	}
}

// --- Portfolio Tests ---

func TestPortfolioOverview(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJiraWithProjects{mockJira: mockJira{
			searchResult: &jiradom.SearchResult{Total: 5},
		}},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.PortfolioOverview(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Portfolio Overview") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "PROJ") {
		t.Errorf("expected project key, got: %s", text)
	}
}

func TestPortfolioBlockers_Clean(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.PortfolioBlockers(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Clean") {
		t.Errorf("expected clean message, got: %s", text)
	}
}

// --- Issue Ops Tests ---

func TestAssignIssue(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.AssignIssue(context.Background(), makeReq(map[string]any{"key": "T-1", "assignee_id": "user1"}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	if !strings.Contains(resultText(result), "T-1 assigned") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestUnassignIssue(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.UnassignIssue(context.Background(), makeReq(map[string]any{"key": "T-1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "unassigned") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestDeleteIssue(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.DeleteIssue(context.Background(), makeReq(map[string]any{"key": "T-99"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "T-99 deleted") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestCreateSubtask(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.CreateSubtask(context.Background(), makeReq(map[string]any{
		"parent_key": "PROJ-10",
		"summary":    "Sub item",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "SUB-1") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestFindUser(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.FindUser(context.Background(), makeReq(map[string]any{"query": "alice"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No users found") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

// --- Link & Worklog Tests ---

func TestLinkIssues(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.LinkIssues(context.Background(), makeReq(map[string]any{
		"inward_key":  "T-1",
		"outward_key": "T-2",
		"link_type":   "Blocks",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Linked T-1 -> T-2") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestWorklogAdd(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.WorklogAdd(context.Background(), makeReq(map[string]any{
		"key":        "T-1",
		"time_spent": "2h",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Logged 2h on T-1") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestWorklogList_Empty(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.WorklogList(context.Background(), makeReq(map[string]any{"key": "T-1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No worklogs") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestWatch(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.Watch(context.Background(), makeReq(map[string]any{"key": "T-1", "account_id": "user1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Added watcher") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestWatchers_Empty(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.Watchers(context.Background(), makeReq(map[string]any{"key": "T-1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No watchers") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestLabelsSet(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.LabelsSet(context.Background(), makeReq(map[string]any{"key": "T-1", "labels": "bug,urgent"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Labels set on T-1") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

// --- Bulk Project Tests ---

func TestListProjects(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJiraWithProjects{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.ListProjects(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "PROJ") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestProjectDetail(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.ProjectDetail(context.Background(), makeReq(map[string]any{"key": "TEST"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "TEST") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestRawRequest(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.RawRequest(context.Background(), makeReq(map[string]any{
		"method": "GET",
		"path":   "/rest/api/2/myself",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "200") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestLinkPR(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.LinkPR(context.Background(), makeReq(map[string]any{
		"key":    "T-1",
		"pr_url": "https://github.com/org/repo/pull/42",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Linked PR to T-1") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestBulkLabel(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.BulkLabel(context.Background(), makeReq(map[string]any{
		"issue_keys": "T-1,T-2",
		"label":      "tech-debt",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "2 succeeded") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestStartSprintTool(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.StartSprintTool(context.Background(), makeReq(map[string]any{
		"sprint_id":  float64(1),
		"start_date": "2026-01-01",
		"end_date":   "2026-01-14",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Sprint 1 started") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestCloseSprintTool(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.CloseSprintTool(context.Background(), makeReq(map[string]any{"sprint_id": float64(5)}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Sprint 5 closed") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestLinkTypes(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.LinkTypes(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "link types") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}
