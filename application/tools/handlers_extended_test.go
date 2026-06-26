package tools_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	jiradom "github.com/aldok10/zara-jira-mcp/domain/jira"
)

// --- GetBoards ---

func TestGetBoards(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			boards: []jiradom.Board{
				{ID: 1, Name: "Dev Board", Type: "scrum"},
				{ID: 2, Name: "Kanban", Type: "kanban"},
			},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.GetBoards(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Dev Board") || !strings.Contains(text, "Kanban") {
		t.Errorf("expected both boards, got: %s", text)
	}
}

func TestGetBoards_Empty(t *testing.T) {
	h := newTestHandlers(&mockJira{boards: []jiradom.Board{}})
	result, _ := h.GetBoards(context.Background(), makeReq(nil))
	text := resultText(result)
	if text != "" {
		t.Errorf("expected empty output for no boards, got: %s", text)
	}
}

// --- GetSprintSummary ---

func TestGetSprintSummary(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 5", State: "active", Goal: "Ship v2"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done", Summary: "task1"}, {Key: "T-2", Status: "In Progress", Summary: "task2"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.GetSprintSummary(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint 5") {
		t.Errorf("expected sprint name, got: %s", text)
	}
	if !strings.Contains(text, "Ship v2") {
		t.Errorf("expected goal, got: %s", text)
	}
	if !strings.Contains(text, "Total: 2 issues") {
		t.Errorf("expected total count, got: %s", text)
	}
}

func TestGetSprintSummary_MissingBoardID(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.GetSprintSummary(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}

func TestGetSprintSummary_NoActiveSprint(t *testing.T) {
	h := newTestHandlers(&mockJira{sprints: []jiradom.Sprint{}})
	result, _ := h.GetSprintSummary(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	text := resultText(result)
	if !strings.Contains(text, "No active sprints") {
		t.Errorf("expected no sprints message, got: %s", text)
	}
}

// --- AddComment ---

func TestAddComment(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.AddComment(context.Background(), makeReq(map[string]any{"key": "T-1", "body": "test comment"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Comment added to T-1") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestAddComment_MissingKey(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.AddComment(context.Background(), makeReq(map[string]any{"body": "hi"}))
	if !result.IsError {
		t.Error("expected error for missing key")
	}
}

func TestAddComment_MissingBody(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.AddComment(context.Background(), makeReq(map[string]any{"key": "T-1"}))
	if !result.IsError {
		t.Error("expected error for missing body")
	}
}

// --- TransitionIssue ---

func TestTransitionIssue(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.TransitionIssue(context.Background(), makeReq(map[string]any{"key": "T-1", "transition_id": "21"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "T-1 transitioned successfully") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestTransitionIssue_MissingKey(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.TransitionIssue(context.Background(), makeReq(map[string]any{"transition_id": "21"}))
	if !result.IsError {
		t.Error("expected error for missing key")
	}
}

func TestTransitionIssue_MissingTransitionID(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.TransitionIssue(context.Background(), makeReq(map[string]any{"key": "T-1"}))
	if !result.IsError {
		t.Error("expected error for missing transition_id")
	}
}

// --- GetTransitions ---

func TestGetTransitions(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			transitions: []jiradom.Transition{{ID: "21", Name: "In Progress"}, {ID: "31", Name: "Done"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.GetTransitions(context.Background(), makeReq(map[string]any{"key": "T-1"}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "In Progress") || !strings.Contains(text, "Done") {
		t.Errorf("expected transitions, got: %s", text)
	}
}

func TestGetTransitions_MissingKey(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.GetTransitions(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing key")
	}
}

// --- Workload ---

func TestWorkload(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			searchResult: &jiradom.SearchResult{
				Issues: []jiradom.Issue{
					{Key: "T-1", Assignee: "alice"},
					{Key: "T-2", Assignee: "alice"},
					{Key: "T-3", Assignee: "bob"},
				},
				Total: 3,
			},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.Workload(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "alice") || !strings.Contains(text, "bob") {
		t.Errorf("expected team members, got: %s", text)
	}
	if !strings.Contains(text, "3 open issues") {
		t.Errorf("expected count, got: %s", text)
	}
}

// --- NotifyLark ---

func TestNotifyLark(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, err := h.NotifyLark(context.Background(), makeReq(map[string]any{"content": "hello", "title": "Test"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Message sent to Lark") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestNotifyLark_MissingContent(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.NotifyLark(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing content")
	}
}

// --- AIAnalyze ---

func TestAIAnalyze(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			searchResult: &jiradom.SearchResult{
				Issues: []jiradom.Issue{{Key: "T-1", Summary: "bug", Status: "Open"}},
				Total:  1,
			},
		},
		AI:     &mockAI{response: "Analysis: one open bug found"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.AIAnalyze(context.Background(), makeReq(map[string]any{"query": "what are the blockers?"}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Analysis") {
		t.Errorf("expected AI response, got: %s", text)
	}
}

func TestAIAnalyze_MissingQuery(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.AIAnalyze(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing query")
	}
}

// --- AISprintReport ---

func TestAISprintReport(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 5", State: "active", Goal: "Ship"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done"}},
		},
		AI:     &mockAI{response: "Sprint is on track"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.AISprintReport(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Sprint is on track") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestAISprintReport_MissingBoardID(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.AISprintReport(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}

func TestAISprintReport_NoActiveSprint(t *testing.T) {
	h := newTestHandlers(&mockJira{sprints: []jiradom.Sprint{}})
	result, _ := h.AISprintReport(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if !strings.Contains(resultText(result), "No active sprints") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}
