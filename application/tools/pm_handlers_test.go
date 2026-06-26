package tools_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	jiradom "github.com/aldok10/zara-jira-mcp/domain/jira"
)

func sprintIssues() []jiradom.Issue {
	return []jiradom.Issue{
		{Key: "T-1", Summary: "Login page", Status: "Done", Type: "Story", Assignee: "alice", Created: time.Now().Add(-10 * 24 * time.Hour), Updated: time.Now().Add(-2 * 24 * time.Hour)},
		{Key: "T-2", Summary: "Fix crash", Status: "Done", Type: "Bug", Assignee: "bob", Created: time.Now().Add(-8 * 24 * time.Hour), Updated: time.Now().Add(-1 * 24 * time.Hour)},
		{Key: "T-3", Summary: "API endpoint", Status: "In Progress", Type: "Task", Assignee: "alice", Created: time.Now().Add(-5 * 24 * time.Hour), Updated: time.Now()},
		{Key: "T-4", Summary: "Blocked thing", Status: "Blocked", Type: "Task", Assignee: "charlie", Created: time.Now().Add(-4 * 24 * time.Hour), Updated: time.Now()},
		{Key: "T-5", Summary: "Todo item", Status: "To Do", Type: "Story", Assignee: "bob", Created: time.Now().Add(-3 * 24 * time.Hour), Updated: time.Now()},
	}
}

func pmHandlers() *tools.Handlers {
	return &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active", Goal: "Ship login"}},
			sprintIssues: sprintIssues(),
		},
		AI:     &mockAI{response: "AI analysis result"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
}

func TestSnapshotSprint(t *testing.T) {
	h := pmHandlers()
	result, err := h.SnapshotSprint(context.Background(), makeReq(map[string]any{
		"board_id": float64(1),
		"velocity": float64(21),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint 10") {
		t.Errorf("expected sprint name, got: %s", text)
	}
	if !strings.Contains(text, "Done: 2") {
		t.Errorf("expected done count, got: %s", text)
	}
}

func TestRecordRisk(t *testing.T) {
	h := pmHandlers()
	result, err := h.RecordRisk(context.Background(), makeReq(map[string]any{
		"title":      "External API might break",
		"severity":   "high",
		"owner":      "alice",
		"mitigation": "Add fallback",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "External API") {
		t.Errorf("expected risk title, got: %s", text)
	}
	if !strings.Contains(text, "high") {
		t.Errorf("expected severity, got: %s", text)
	}
}

func TestRecordDecision(t *testing.T) {
	h := pmHandlers()
	result, err := h.RecordDecision(context.Background(), makeReq(map[string]any{
		"title":     "Use PostgreSQL",
		"decision":  "PostgreSQL over MongoDB",
		"rationale": "ACID compliance needed",
		"tags":      "architecture,database",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "PostgreSQL") {
		t.Errorf("expected decision, got: %s", text)
	}
}

func TestRecordBlocker(t *testing.T) {
	h := pmHandlers()
	result, err := h.RecordBlocker(context.Background(), makeReq(map[string]any{
		"description": "Waiting for design approval",
		"issue_key":   "T-3",
		"owner":       "pm-lead",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Waiting for design") {
		t.Errorf("expected blocker desc, got: %s", text)
	}
}

func TestSprintHealthScore(t *testing.T) {
	h := pmHandlers()
	result, err := h.SprintHealthScore(context.Background(), makeReq(map[string]any{
		"board_id": float64(1),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint Health:") {
		t.Errorf("expected health header, got: %s", text)
	}
	if !strings.Contains(text, "/100") {
		t.Errorf("expected score, got: %s", text)
	}
}

func TestFlowMetrics(t *testing.T) {
	h := pmHandlers()
	result, err := h.FlowMetrics(context.Background(), makeReq(map[string]any{
		"board_id": float64(1),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "WIP") {
		t.Errorf("expected WIP metric, got: %s", text)
	}
	if !strings.Contains(text, "Throughput") {
		t.Errorf("expected throughput, got: %s", text)
	}
}

func TestCeremonyFacilitator(t *testing.T) {
	h := pmHandlers()
	result, err := h.CeremonyFacilitator(context.Background(), makeReq(map[string]any{
		"ceremony": "retro",
		"board_id": float64(1),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty result")
	}
}

func TestCeremonyFacilitator_InvalidCeremony(t *testing.T) {
	h := pmHandlers()
	result, _ := h.CeremonyFacilitator(context.Background(), makeReq(map[string]any{
		"ceremony": "invalid",
	}))
	if !result.IsError {
		t.Error("expected error for invalid ceremony")
	}
}

func TestPMDashboard(t *testing.T) {
	h := pmHandlers()
	result, err := h.PMDashboard(context.Background(), makeReq(map[string]any{
		"board_id": float64(1),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "PM DASHBOARD") {
		t.Errorf("expected dashboard header, got: %s", text)
	}
	if !strings.Contains(text, "Sprint 10") {
		t.Errorf("expected sprint info, got: %s", text)
	}
}

func TestRecordDailyProgress(t *testing.T) {
	h := pmHandlers()
	result, err := h.TrackDailyProgress(context.Background(), makeReq(map[string]any{
		"board_id": float64(1),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "Daily progress") || !strings.Contains(text, "Sprint 10") {
		t.Errorf("expected daily progress output, got: %s", text)
	}
}

func TestMonteCarloForecast_InsufficientData(t *testing.T) {
	h := pmHandlers() // mockMemory returns nil snapshots
	result, err := h.MonteCarloForecast(context.Background(), makeReq(map[string]any{
		"board_id": float64(1),
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "at least 3") {
		t.Errorf("expected insufficient data message, got: %s", text)
	}
}

func TestManageDoD_EmptyList(t *testing.T) {
	h := pmHandlers()
	result, err := h.ManageDoD(context.Background(), makeReq(map[string]any{
		"action": "list",
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "No Definition of Done") {
		t.Errorf("expected empty message, got: %s", text)
	}
}

func TestRecordConfidence(t *testing.T) {
	h := pmHandlers()
	result, err := h.RecordConfidence(context.Background(), makeReq(map[string]any{
		"sprint_name": "Sprint 10",
		"score":       float64(4),
		"member":      "alice",
		"note":        "feeling good about scope",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	text := resultText(result)
	if !strings.Contains(text, "confident") {
		t.Errorf("expected confidence label, got: %s", text)
	}
}

func TestRecordConfidence_InvalidScore(t *testing.T) {
	h := pmHandlers()
	result, _ := h.RecordConfidence(context.Background(), makeReq(map[string]any{
		"sprint_name": "Sprint 10",
		"score":       float64(0),
	}))
	if !result.IsError {
		t.Error("expected error for invalid score")
	}
}

func TestExecutiveReport(t *testing.T) {
	h := pmHandlers()
	result, err := h.ExecutiveReport(context.Background(), makeReq(map[string]any{
		"board_id": float64(1),
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	// Should return AI response
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty report")
	}
}

func TestTeamKnowledgeBase_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.TeamKnowledgeBase(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "empty") {
		t.Errorf("expected empty KB message, got: %s", text)
	}
}
