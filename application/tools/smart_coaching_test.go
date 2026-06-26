package tools_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	jiradom "github.com/aldok10/zara-jira-mcp/domain/jira"
	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
)

// --- Smart Handler Tests ---

func TestPMSmart_Blockers(t *testing.T) {
	h := pmHandlers()
	result, err := h.PMSmart(context.Background(), makeReq(map[string]any{"ask": "show me blockers"}))
	if err != nil {
		t.Fatal(err)
	}
	// Routes to GetBlockers which returns "No active blockers"
	if !strings.Contains(resultText(result), "No active blockers") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMSmart_Risk(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMSmart(context.Background(), makeReq(map[string]any{"ask": "any risk?"}))
	if !strings.Contains(resultText(result), "No open risks") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMSmart_MyIssues(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{searchResult: &jiradom.SearchResult{Issues: []jiradom.Issue{{Key: "T-1", Summary: "task", Status: "Open"}}, Total: 1}},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, _ := h.PMSmart(context.Background(), makeReq(map[string]any{"ask": "my issues assigned to me"}))
	if !strings.Contains(resultText(result), "T-1") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMSmart_ActionItems(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMSmart(context.Background(), makeReq(map[string]any{"ask": "pending action items"}))
	if !strings.Contains(resultText(result), "No pending") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMSmart_Help(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMSmart(context.Background(), makeReq(map[string]any{"ask": "help me"}))
	text := resultText(result)
	if text == "" {
		t.Error("expected help text")
	}
}

func TestPMSmart_MissingAsk(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMSmart(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing ask")
	}
}

func TestPMSmart_Fallback(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{searchResult: &jiradom.SearchResult{}},
		AI:     &mockAI{response: "AI fallback response"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, _ := h.PMSmart(context.Background(), makeReq(map[string]any{"ask": "something random"}))
	text := resultText(result)
	if !strings.Contains(text, "AI fallback response") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestPMDo_Risk(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMDo(context.Background(), makeReq(map[string]any{"what": "record risk", "title": "API outage", "severity": "high"}))
	if !strings.Contains(resultText(result), "Risk recorded") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMDo_Decision(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMDo(context.Background(), makeReq(map[string]any{"what": "record decision", "title": "Use Kafka", "decision": "Kafka over RabbitMQ"}))
	if !strings.Contains(resultText(result), "Decision recorded") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMDo_Blocker(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMDo(context.Background(), makeReq(map[string]any{"what": "record blocker", "description": "waiting on infra"}))
	if !strings.Contains(resultText(result), "Blocker recorded") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMDo_Unknown(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMDo(context.Background(), makeReq(map[string]any{"what": "something unknown"}))
	if !result.IsError {
		t.Error("expected error for unknown action")
	}
}

func TestPMDo_MissingWhat(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMDo(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing what")
	}
}

func TestPMReport_Status(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMReport(context.Background(), makeReq(map[string]any{"type": "status", "board_id": float64(1)}))
	text := resultText(result)
	if !strings.Contains(text, "Overall:") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestPMReport_Unknown(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMReport(context.Background(), makeReq(map[string]any{"type": "unknown"}))
	if !result.IsError {
		t.Error("expected error for unknown type")
	}
}

func TestPMTeam_Unknown(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMTeam(context.Background(), makeReq(map[string]any{"action": "unknown"}))
	if !result.IsError {
		t.Error("expected error for unknown action")
	}
}

func TestPMPlan_Unknown(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMPlan(context.Background(), makeReq(map[string]any{"action": "unknown"}))
	if !result.IsError {
		t.Error("expected error for unknown action")
	}
}

func TestPMRetro_Actions(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMRetro(context.Background(), makeReq(map[string]any{"action": "actions"}))
	if !strings.Contains(resultText(result), "No pending") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMRetro_Unknown(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMRetro(context.Background(), makeReq(map[string]any{"action": "unknown"}))
	if !result.IsError {
		t.Error("expected error for unknown action")
	}
}

func TestPMSearch_Decisions(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMSearch(context.Background(), makeReq(map[string]any{"query": "decision about db"}))
	if !strings.Contains(resultText(result), "No decisions") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMSearch_Empty(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMSearch(context.Background(), makeReq(map[string]any{"query": ""}))
	if !result.IsError {
		t.Error("expected error for empty query")
	}
}

// --- Coaching Handler Tests ---

func TestPMTeamPulse(t *testing.T) {
	h := pmHandlers()
	result, err := h.PMTeamPulse(context.Background(), makeReq(map[string]any{
		"sprint_name": "Sprint 10",
		"ratings":     `{"alice": 4, "bob": 3}`,
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Team pulse") {
		t.Errorf("unexpected: %s", text)
	}
	if !strings.Contains(text, "2 member ratings") {
		t.Errorf("expected 2 members, got: %s", text)
	}
}

func TestPMTeamPulse_MissingSprint(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMTeamPulse(context.Background(), makeReq(map[string]any{"ratings": `{}`}))
	if !result.IsError {
		t.Error("expected error for missing sprint_name")
	}
}

func TestPMTeamPulse_MissingRatings(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMTeamPulse(context.Background(), makeReq(map[string]any{"sprint_name": "S10"}))
	if !result.IsError {
		t.Error("expected error for missing ratings")
	}
}

func TestPMTeamPulse_InvalidJSON(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMTeamPulse(context.Background(), makeReq(map[string]any{"sprint_name": "S10", "ratings": "not json"}))
	if !result.IsError {
		t.Error("expected error for invalid JSON")
	}
}

func TestPMTeamPulseHistory_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.PMTeamPulseHistory(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No pulse data") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMPredictability_InsufficientData(t *testing.T) {
	h := pmHandlers()
	result, err := h.PMPredictability(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Need at least 3") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestPMPredictability_WithData(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithSnapshots{},
	}
	result, err := h.PMPredictability(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Predictability") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestPMPredictability_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMPredictability(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}

func TestPMMeetingEffectiveness(t *testing.T) {
	h := pmHandlers()
	result, err := h.PMMeetingEffectiveness(context.Background(), makeReq(map[string]any{
		"ceremony":         "standup",
		"duration_minutes": float64(10),
		"score":            float64(4),
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if text == "" {
		t.Error("expected non-empty result")
	}
}

func TestPMMeetingEffectiveness_MissingCeremony(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMMeetingEffectiveness(context.Background(), makeReq(map[string]any{"duration_minutes": float64(10), "score": float64(4)}))
	if !result.IsError {
		t.Error("expected error for missing ceremony")
	}
}

func TestPMMeetingEffectiveness_MissingDuration(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMMeetingEffectiveness(context.Background(), makeReq(map[string]any{"ceremony": "standup", "score": float64(4)}))
	if !result.IsError {
		t.Error("expected error for missing duration")
	}
}

func TestPMMeetingEffectiveness_MissingScore(t *testing.T) {
	h := pmHandlers()
	result, _ := h.PMMeetingEffectiveness(context.Background(), makeReq(map[string]any{"ceremony": "standup", "duration_minutes": float64(10)}))
	if !result.IsError {
		t.Error("expected error for missing score")
	}
}

// --- Advanced Handler Tests ---

func TestSprintHealthScore_NoSprint(t *testing.T) {
	h := newTestHandlers(&mockJira{sprints: []jiradom.Sprint{}})
	result, _ := h.SprintHealthScore(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if !strings.Contains(resultText(result), "No active sprint") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestSprintHealthScore_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.SprintHealthScore(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}

// --- Dependency Recording Tests ---

type mockMemoryWithDeps struct {
	mockMemory
}

func (m *mockMemoryWithDeps) GetOpenDependencies(_ context.Context) ([]memdom.Dependency, error) {
	return []memdom.Dependency{
		{ID: 1, FromIssueKey: "T-1", ToIssueKey: "T-2", DependencyType: "blocks", Description: "blocks api"},
		{ID: 2, FromIssueKey: "T-3", ToIssueKey: "EXT-1", DependencyType: "external", Description: "external dep"},
	}, nil
}

func TestCrossTeamDependencyReport_WithDeps(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithDeps{},
	}
	result, err := h.CrossTeamDependencyReport(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "CROSS-TEAM") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "EXTERNAL") {
		t.Errorf("expected external deps, got: %s", text)
	}
}
