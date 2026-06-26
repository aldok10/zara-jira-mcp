package tools_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	jiradom "github.com/aldok10/zara-jira-mcp/domain/jira"
)

// --- Deep Handler Tests ---

func TestGetBurndown_NoData(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetBurndown(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No daily progress") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestGetBurndown_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.GetBurndown(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}

func TestSetSprintGoal_Basic(t *testing.T) {
	h := pmHandlers()
	result, err := h.SetSprintGoal(context.Background(), makeReq(map[string]any{
		"board_id":    float64(1),
		"goal":        "Ship authentication module",
		"key_results": "Login works\nOAuth integrated",
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint goal set") {
		t.Errorf("unexpected: %s", text)
	}
	if !strings.Contains(text, "Ship authentication") {
		t.Errorf("expected goal text, got: %s", text)
	}
}

func TestSetSprintGoal_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.SetSprintGoal(context.Background(), makeReq(map[string]any{"goal": "x"}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestSetSprintGoal_MissingGoal(t *testing.T) {
	h := pmHandlers()
	result, _ := h.SetSprintGoal(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestCloseSprintGoal(t *testing.T) {
	h := pmHandlers()
	result, err := h.CloseSprintGoal(context.Background(), makeReq(map[string]any{
		"goal_id": float64(1),
		"status":  "achieved",
		"outcome": "Shipped on time",
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "closed") || !strings.Contains(text, "achieved") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestCloseSprintGoal_MissingID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.CloseSprintGoal(context.Background(), makeReq(map[string]any{"status": "achieved"}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestCloseSprintGoal_MissingStatus(t *testing.T) {
	h := pmHandlers()
	result, _ := h.CloseSprintGoal(context.Background(), makeReq(map[string]any{"goal_id": float64(1)}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestGetSprintGoals_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetSprintGoals(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No goals") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestManageDoD_Add_Basic(t *testing.T) {
	h := pmHandlers()
	result, err := h.ManageDoD(context.Background(), makeReq(map[string]any{
		"action": "add",
		"item":   "All tests pass",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "added") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestManageDoD_AddMissingItem(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ManageDoD(context.Background(), makeReq(map[string]any{"action": "add"}))
	if !result.IsError {
		t.Error("expected error for missing item")
	}
}

func TestManageDoD_Remove(t *testing.T) {
	h := pmHandlers()
	result, err := h.ManageDoD(context.Background(), makeReq(map[string]any{"action": "remove", "item_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "removed") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestManageDoD_RemoveMissingID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ManageDoD(context.Background(), makeReq(map[string]any{"action": "remove"}))
	if !result.IsError {
		t.Error("expected error")
	}
}

// --- Advanced Handler Tests (Dependency, Meeting) ---

func TestRecordDependency(t *testing.T) {
	h := pmHandlers()
	result, err := h.RecordDependency(context.Background(), makeReq(map[string]any{
		"from_issue":  "T-1",
		"to_issue":    "T-2",
		"type":        "blocks",
		"description": "T-1 needs T-2 API",
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Dependency recorded") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestRecordDependency_MissingFrom(t *testing.T) {
	h := pmHandlers()
	result, _ := h.RecordDependency(context.Background(), makeReq(map[string]any{"to_issue": "T-2"}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestRecordDependency_MissingTo(t *testing.T) {
	h := pmHandlers()
	result, _ := h.RecordDependency(context.Background(), makeReq(map[string]any{"from_issue": "T-1"}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestResolveDependency(t *testing.T) {
	h := pmHandlers()
	result, err := h.ResolveDependency(context.Background(), makeReq(map[string]any{"dependency_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "resolved") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestResolveDependency_MissingID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ResolveDependency(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestGetDependencies_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetDependencies(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No open dependencies") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestRecordMeeting(t *testing.T) {
	h := pmHandlers()
	result, err := h.RecordMeeting(context.Background(), makeReq(map[string]any{
		"meeting_type": "standup",
		"notes":        "Discussed progress",
		"decisions":    "Ship by Friday",
		"action_items": "alice: review PR",
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Meeting recorded") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestRecordMeeting_MissingType(t *testing.T) {
	h := pmHandlers()
	result, _ := h.RecordMeeting(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestGetMeetings_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetMeetings(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No meeting notes") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestCapacityPlan_InsufficientData(t *testing.T) {
	h := pmHandlers()
	result, err := h.CapacityPlan(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Not enough") && !strings.Contains(text, "Need at least") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestCapacityPlan_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.CapacityPlan(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestCapacityPlan_WithData(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "S10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Assignee: "alice"}, {Key: "T-2", Assignee: "bob"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithSnapshots{},
	}
	result, err := h.CapacityPlan(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Capacity") {
		t.Errorf("expected capacity info, got: %s", text)
	}
}

// --- Flow Handler Tests ---

func TestFlowMetrics_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.FlowMetrics(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

// --- Forecast Handler Tests ---

func TestForecastSprint_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ForecastSprint(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestMonteCarloForecast_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.MonteCarloForecast(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}
