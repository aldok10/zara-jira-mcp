package tools_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	jiradom "github.com/aldok10/zara-jira-mcp/domain/jira"
	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
)

// --- Memory Handler Tests ---

func TestUpdateRisk(t *testing.T) {
	h := pmHandlers()
	result, err := h.UpdateRisk(context.Background(), makeReq(map[string]any{
		"risk_id": float64(1),
		"status":  "resolved",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if result.IsError {
		t.Fatal(resultText(result))
	}
	if !strings.Contains(resultText(result), "updated to status: resolved") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestUpdateRisk_MissingID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.UpdateRisk(context.Background(), makeReq(map[string]any{"status": "resolved"}))
	if !result.IsError {
		t.Error("expected error for missing risk_id")
	}
}

func TestUpdateRisk_MissingStatus(t *testing.T) {
	h := pmHandlers()
	result, _ := h.UpdateRisk(context.Background(), makeReq(map[string]any{"risk_id": float64(1)}))
	if !result.IsError {
		t.Error("expected error for missing status")
	}
}

func TestGetRiskDashboard_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetRiskDashboard(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No open risks") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

type mockMemoryWithRisks struct {
	mockMemory
}

func (m *mockMemoryWithRisks) GetOpenRisks(_ context.Context) ([]memdom.Risk, error) {
	return []memdom.Risk{
		{ID: 1, Title: "API failure", Severity: "high", Owner: "alice", SprintName: "S10"},
	}, nil
}

func TestGetRiskDashboard_WithRisks(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithRisks{},
	}
	result, err := h.GetRiskDashboard(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Open Risks (1)") {
		t.Errorf("expected risks header, got: %s", text)
	}
	if !strings.Contains(text, "API failure") {
		t.Errorf("expected risk title, got: %s", text)
	}
}

func TestSearchDecisions_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.SearchDecisions(context.Background(), makeReq(map[string]any{"query": "anything"}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No decisions found") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestRecordBlocker_MissingDescription(t *testing.T) {
	h := pmHandlers()
	result, _ := h.RecordBlocker(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing description")
	}
}

func TestResolveBlocker(t *testing.T) {
	h := pmHandlers()
	result, err := h.ResolveBlocker(context.Background(), makeReq(map[string]any{
		"blocker_id": float64(1),
		"resolution": "Fixed by infra team",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "Blocker #1 resolved") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestResolveBlocker_MissingID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ResolveBlocker(context.Background(), makeReq(map[string]any{"resolution": "done"}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestResolveBlocker_MissingResolution(t *testing.T) {
	h := pmHandlers()
	result, _ := h.ResolveBlocker(context.Background(), makeReq(map[string]any{"blocker_id": float64(1)}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestGetBlockers_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetBlockers(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No active blockers") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestGetBlockers_ShowHistory(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetBlockers(context.Background(), makeReq(map[string]any{"show_history": true}))
	if err != nil {
		t.Fatal(err)
	}
	// mockMemory returns nil/empty, so either "No active" or empty history
	text := resultText(result)
	if strings.Contains(text, "error") {
		t.Errorf("unexpected error: %s", text)
	}
}

func TestRecordTeamMetric(t *testing.T) {
	h := pmHandlers()
	result, err := h.RecordTeamMetric(context.Background(), makeReq(map[string]any{
		"member_name":     "alice",
		"sprint_name":     "Sprint 10",
		"issues_assigned": float64(8),
		"issues_done":     float64(6),
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "alice") || !strings.Contains(text, "Sprint 10") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestRecordTeamMetric_MissingMember(t *testing.T) {
	h := pmHandlers()
	result, _ := h.RecordTeamMetric(context.Background(), makeReq(map[string]any{"sprint_name": "S1"}))
	if !result.IsError {
		t.Error("expected error for missing member_name")
	}
}

func TestRecordTeamMetric_MissingSprint(t *testing.T) {
	h := pmHandlers()
	result, _ := h.RecordTeamMetric(context.Background(), makeReq(map[string]any{"member_name": "a"}))
	if !result.IsError {
		t.Error("expected error for missing sprint_name")
	}
}

func TestGetTeamHealth_MissingParams(t *testing.T) {
	h := pmHandlers()
	result, _ := h.GetTeamHealth(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error when neither member nor sprint provided")
	}
}

func TestGetTeamHealth_ByMember(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetTeamHealth(context.Background(), makeReq(map[string]any{"member_name": "alice"}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "No metrics") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestGetTeamHealth_BySprint(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetTeamHealth(context.Background(), makeReq(map[string]any{"sprint_name": "S10"}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "No metrics") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestRecordRetrospective(t *testing.T) {
	h := pmHandlers()
	result, err := h.RecordRetrospective(context.Background(), makeReq(map[string]any{
		"sprint_name":  "Sprint 10",
		"went_well":    "Good collaboration",
		"improvements": "Reduce meetings",
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint 10") || !strings.Contains(text, "Good collaboration") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestRecordRetrospective_MissingSprint(t *testing.T) {
	h := pmHandlers()
	result, _ := h.RecordRetrospective(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing sprint_name")
	}
}

func TestGetActionItems_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.GetActionItems(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No pending action items") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

// --- Routing Handler Tests ---

func TestNotifyRouted_NoChannels(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Memory: &mockMemory{},
		// No Slack, Telegram, Teams, Discord, Email, Lark
	}
	result, err := h.NotifyRouted(context.Background(), makeReq(map[string]any{
		"content":  "test message",
		"severity": "info",
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "No channels configured") {
		t.Errorf("expected no channels message, got: %s", text)
	}
}

func TestNotifyRouted_MissingContent(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.NotifyRouted(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing content")
	}
}

func TestNotifyRouted_LarkChannel(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.NotifyRouted(context.Background(), makeReq(map[string]any{
		"content":  "urgent message",
		"severity": "medium",
	}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Lark: sent") {
		t.Errorf("expected Lark sent, got: %s", text)
	}
}

func TestDailyDigest_NothingToReport(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{searchResult: &jiradom.SearchResult{}},
		AI:     &mockAI{response: "ok"},
		Memory: &mockMemory{},
	}
	result, err := h.DailyDigest(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Nothing to report") {
		t.Errorf("expected empty digest, got: %s", text)
	}
}

// --- Database Handler Tests ---

func TestDatabaseQuery_NilDatabase(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
		// Database is nil
	}
	result, err := h.DatabaseQuery(context.Background(), makeReq(map[string]any{"query": "SELECT 1"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Error("expected error for nil database")
	}
	if !strings.Contains(resultText(result), "No database configured") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestDatabaseListTables_NilDatabase(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.DatabaseListTables(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Error("expected error for nil database")
	}
	if !strings.Contains(resultText(result), "No database configured") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestMongoQuery_NilDatabase(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.MongoQuery(context.Background(), makeReq(map[string]any{"collection": "test"}))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Error("expected error for nil database")
	}
	if !strings.Contains(resultText(result), "MongoDB not configured") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestMongoListCollections_NilDatabase(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.MongoListCollections(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Error("expected error for nil database")
	}
}

// --- Management Handler Tests ---

func TestManagementBrief(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active", Goal: "Ship v2"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done"}, {Key: "T-2", Status: "In Progress"}, {Key: "T-3", Status: "To Do"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.ManagementBrief(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint 10") {
		t.Errorf("expected sprint name, got: %s", text)
	}
	if !strings.Contains(text, "On Track") && !strings.Contains(text, "At Risk") && !strings.Contains(text, "Behind") {
		t.Errorf("expected status, got: %s", text)
	}
}

func TestManagementBrief_MissingBoardID(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.ManagementBrief(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}

func TestManagementBrief_NoSprint(t *testing.T) {
	h := newTestHandlers(&mockJira{sprints: []jiradom.Sprint{}})
	result, _ := h.ManagementBrief(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if !result.IsError {
		t.Error("expected error for no active sprint")
	}
}

func TestManagementBrief_DirectorAudience(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done"}, {Key: "T-2", Status: "Blocked"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.ManagementBrief(context.Background(), makeReq(map[string]any{"board_id": float64(1), "audience": "director"}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint 10") {
		t.Errorf("expected sprint, got: %s", text)
	}
}

func TestDependencyReport_NoDeps(t *testing.T) {
	h := pmHandlers()
	result, err := h.DependencyReport(context.Background(), makeReq(map[string]any{}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "No dependencies recorded") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestEscalationReport_NothingToEscalate(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done"}},
			searchResult: &jiradom.SearchResult{Issues: nil, Total: 0},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.EscalationReport(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "No items currently require escalation") {
		t.Errorf("unexpected: %s", text)
	}
}

func TestResourceUtilization(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints: []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: []jiradom.Issue{
				{Key: "T-1", Status: "Done", Assignee: "alice"},
				{Key: "T-2", Status: "In Progress", Assignee: "alice"},
				{Key: "T-3", Status: "In Progress", Assignee: "bob"},
				{Key: "T-4", Status: "In Progress", Assignee: "bob"},
				{Key: "T-5", Status: "In Progress", Assignee: "bob"},
				{Key: "T-6", Status: "In Progress", Assignee: "bob"},
				{Key: "T-7", Status: "In Progress", Assignee: "bob"},
			},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemory{},
	}
	result, err := h.ResourceUtilization(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Resource Utilization") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "OVERLOADED") {
		t.Errorf("expected overloaded flag for bob, got: %s", text)
	}
}

func TestResourceUtilization_MissingBoardID(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.ResourceUtilization(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}

func TestBlockerAgingReport_NilMemory(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: nil,
	}
	result, err := h.BlockerAgingReport(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !result.IsError {
		t.Error("expected error for nil memory")
	}
}

func TestBlockerAgingReport_Empty(t *testing.T) {
	h := pmHandlers()
	result, err := h.BlockerAgingReport(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resultText(result), "No active blockers") {
		t.Errorf("unexpected: %s", resultText(result))
	}
}

func TestSprintCommitmentReport(t *testing.T) {
	h := &tools.Handlers{
		Jira: &mockJira{
			sprints:      []jiradom.Sprint{{ID: 1, Name: "Sprint 10", State: "active"}},
			sprintIssues: []jiradom.Issue{{Key: "T-1", Status: "Done"}, {Key: "T-2", Status: "In Progress"}},
		},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: &mockMemoryWithSnapshots{},
	}
	result, err := h.SprintCommitmentReport(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if err != nil {
		t.Fatal(err)
	}
	text := resultText(result)
	if !strings.Contains(text, "Sprint Commitment") {
		t.Errorf("expected header, got: %s", text)
	}
	if !strings.Contains(text, "Historical") {
		t.Errorf("expected historical data, got: %s", text)
	}
}

func TestSprintCommitmentReport_NilMemory(t *testing.T) {
	h := &tools.Handlers{
		Jira:   &mockJira{sprints: []jiradom.Sprint{{ID: 1, Name: "S10", State: "active"}}},
		AI:     &mockAI{response: "ok"},
		Lark:   &mockLark{},
		Memory: nil,
	}
	result, _ := h.SprintCommitmentReport(context.Background(), makeReq(map[string]any{"board_id": float64(1)}))
	if !result.IsError {
		t.Error("expected error for nil memory")
	}
}

func TestSprintCommitmentReport_MissingBoardID(t *testing.T) {
	h := pmHandlers()
	result, _ := h.SprintCommitmentReport(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error for missing board_id")
	}
}
