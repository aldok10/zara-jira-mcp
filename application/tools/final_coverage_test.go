package tools_test

import (
	"context"
	"testing"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	jiradom "github.com/aldok10/zara-jira-mcp/domain/jira"
)

func TestCommunicateDecision(t *testing.T) {
	h := &tools.Handlers{Jira: &mockJira{}, AI: &mockAI{response: "Decision communicated"}, Lark: &mockLark{}, Memory: &mockMemory{}}
	result, err := h.CommunicateDecision(context.Background(), makeReq(map[string]any{"decision": "Use REST", "audience": "team"}))
	if err != nil {
		t.Fatal(err)
	}
	if resultText(result) == "" {
		t.Error("expected output")
	}
}

func TestEscalateWithSCQA(t *testing.T) {
	h := &tools.Handlers{Jira: &mockJira{}, AI: &mockAI{response: "Escalation formatted"}, Lark: &mockLark{}, Memory: &mockMemory{}}
	result, err := h.EscalateWithSCQA(context.Background(), makeReq(map[string]any{"situation": "blocked", "complication": "no response", "question": "what do we do?"}))
	if err != nil {
		t.Fatal(err)
	}
	if resultText(result) == "" {
		t.Error("expected output")
	}
}

func TestAdaptMessage(t *testing.T) {
	h := &tools.Handlers{Jira: &mockJira{}, AI: &mockAI{response: "Adapted message"}, Lark: &mockLark{}, Memory: &mockMemory{}}
	result, err := h.AdaptMessage(context.Background(), makeReq(map[string]any{"message": "sprint is behind", "audience": "executive"}))
	if err != nil {
		t.Fatal(err)
	}
	if resultText(result) == "" {
		t.Error("expected output")
	}
}

func TestGiveFeedback(t *testing.T) {
	h := &tools.Handlers{Jira: &mockJira{}, AI: &mockAI{response: "Feedback crafted"}, Lark: &mockLark{}, Memory: &mockMemory{}}
	result, err := h.GiveFeedback(context.Background(), makeReq(map[string]any{"situation": "missed deadline", "behavior": "no communication", "impact": "team blocked"}))
	if err != nil {
		t.Fatal(err)
	}
	if resultText(result) == "" {
		t.Error("expected output")
	}
}

func TestWriteUpdate(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "S10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done"}},
		},
		AI: &mockAI{response: "Update written"}, Lark: &mockLark{}, Memory: &mockMemory{},
	}
	result, err := h.WriteUpdate(context.Background(), makeReq(map[string]any{"board_id": float64(1), "audience": "team"}))
	if err != nil {
		t.Fatal(err)
	}
	if resultText(result) == "" {
		t.Error("expected output")
	}
}

func TestEscalate(t *testing.T) {
	h := pmHandlers()
	result, err := h.Escalate(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if resultText(result) == "" {
		t.Error("expected output")
	}
}

func TestGetEscalations_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetEscalations(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if resultText(result) == "" {
		t.Error("expected output")
	}
}

func TestSprintGoalCheck(t *testing.T) {
	h := pmHandlers()
	result, err := h.SprintGoalCheck(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if resultText(result) == "" {
		t.Error("expected output")
	}
}

func TestPMMaturityAssessment(t *testing.T) {
	h := pmHandlers()
	result, err := h.PMMaturityAssessment(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if resultText(result) == "" {
		t.Error("expected output")
	}
}

func TestPMDailyDigestCoaching(t *testing.T) {
	h := pmHandlers()
	result, err := h.PMDailyDigestCoaching(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if resultText(result) == "" {
		t.Error("expected output")
	}
}
