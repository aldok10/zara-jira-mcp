package tools_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	jiradom "github.com/aldok10/zara-jira-mcp/domain/jira"
	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
)

// --- Reporting Handler Tests ---

func TestReportToPO_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ReportToPO(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}

func TestEscalationBrief_NoBlockers_Basic(t *testing.T) {
	h := pmHandlers()
	result, err := h.EscalationBrief(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No active impediments") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

type mockMemoryWithBlockers struct {
	mockMemory
}

func (m *mockMemoryWithBlockers) GetActiveBlockers(_ context.Context) ([]memdom.Blocker, error) {
	return []memdom.Blocker{
		{ID: 1, IssueKey: "T-5", Description: "Waiting on vendor", Owner: "alice"},
	}, nil
}

func TestEscalationBrief_WithBlockers(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{sprints: []jiradom.Sprint{{ID: 1, Name: "S10", State: "active"}}, sprintIssues: []jiradom.Issue{{Key: "T-1"}}},
		AI:     &mockAI{response: "Escalation: vendor dependency"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithBlockers{},
	}
	result, err := h.EscalationBrief(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty escalation brief")
	}
}

func TestCrossTeamDependencyReport_NoDeps(t *testing.T) {
	h := pmHandlers()
	result, err := h.CrossTeamDependencyReport(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No cross-team dependencies") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestDeliveryConfidenceReport_Basic(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active", Goal: "Deliver v2"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done"}, {Key: "T-2", Status: "In Progress"}, {Key: "T-3", Status: "To Do"}},
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
	if !strings.Contains(text, "DELIVERY CONFIDENCE") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "GREEN") && !strings.Contains(text, "AMBER") && !strings.Contains(text, "RED") {
		t.Errorf("expected color status, got: %s", text)
	}
}

func TestDeliveryConfidenceReport_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.DeliveryConfidenceReport(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}

func TestResourcePlanningReport_NoData(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "S10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Assignee: "alice"}, {Key: "T-2", Assignee: "bob"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.ResourcePlanningReport(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "RESOURCE") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "Insufficient historical data") {
		t.Errorf("expected insufficient data msg, got: %s", text)
	}
}

func TestResourcePlanningReport_WithHistory(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "S10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Assignee: "alice"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithSnapshots{},
	}
	result, err := h.ResourcePlanningReport(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Team Throughput") {
		t.Errorf("expected throughput data, got: %s", text)
	}
}

func TestResourcePlanningReport_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ResourcePlanningReport(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}

// --- Process Handler Tests ---

func TestManageDoR_List_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.ManageDoR(context.Background(), makeReq(map[string]any{"action": "list"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No Definition of Ready") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestManageDoR_Add(t *testing.T) {
	h := pmHandlers()
	result, err := h.ManageDoR(context.Background(), makeReq(map[string]any{
		"action": "add",
		"item":   "Story has acceptance criteria",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "DoR item added") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestManageDoR_AddMissingItem(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ManageDoR(context.Background(), makeReq(map[string]any{"action": "add"}))
	if !result.IsError {
		t.Error("expected error for missing item")
	}
}

func TestManageDoR_Remove(t *testing.T) {
	h := pmHandlers()
	result, err := h.ManageDoR(context.Background(), makeReq(map[string]any{"action": "remove", "item_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "removed") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestManageDoR_RemoveMissingID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ManageDoR(context.Background(), makeReq(map[string]any{"action": "remove"}))
	if !result.IsError {
		t.Error("expected error for missing item_id")
	}
}

func TestCheckStoryReady(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			issue: &jiradom.Issue{Key: "T-1", Summary: "Add login page", Description: "As a user I want to login so that I can access my account", Status: "To Do", Type: "Story"},
		},
		AI:     &mockAI{response: "Story meets DoR criteria"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.CheckStoryReady(context.Background(), makeReq(map[string]any{"key": "T-1"}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty readiness check")
	}
}

func TestCheckStoryReady_MissingKey(t *testing.T) {
	h := pmHandlers()
	result, _ := h.CheckStoryReady(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing key")
	}
}

// --- Stakeholder Handler Tests ---

func TestExecutiveReport_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ExecutiveReport(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}
