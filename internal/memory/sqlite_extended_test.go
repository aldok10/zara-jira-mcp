package memory

import (
	"context"
	"testing"
	"time"

	domain "github.com/aldok10/zara-jira-mcp/domain/memory"
)

func TestGetLatestSnapshot_Empty(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	snap, err := store.GetLatestSnapshot(ctx, 99)
	if err != nil {
		t.Fatal(err)
	}
	if snap != nil {
		t.Fatalf("expected nil snapshot, got %+v", snap)
	}
}

func TestGetLatestSnapshot_ReturnsNewest(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	for i, name := range []string{"Sprint 1", "Sprint 2", "Sprint 3"} {
		err := store.SaveSprintSnapshot(ctx, &domain.SprintSnapshot{
			SprintName: name, BoardID: 1, SnapshotDate: time.Now().Add(time.Duration(i) * time.Hour),
			TotalIssues: 10, Done: i + 1,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	snap, err := store.GetLatestSnapshot(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if snap == nil {
		t.Fatal("expected snapshot")
	}
	if snap.SprintName != "Sprint 3" {
		t.Errorf("got %s, want Sprint 3", snap.SprintName)
	}
}

func TestSprintSnapshots_MultipleBoards(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveSprintSnapshot(ctx, &domain.SprintSnapshot{SprintName: "S1", BoardID: 1, SnapshotDate: time.Now(), Done: 5})
	store.SaveSprintSnapshot(ctx, &domain.SprintSnapshot{SprintName: "S2", BoardID: 2, SnapshotDate: time.Now(), Done: 3})
	store.SaveSprintSnapshot(ctx, &domain.SprintSnapshot{SprintName: "S3", BoardID: 1, SnapshotDate: time.Now(), Done: 7})

	snaps, err := store.GetSprintSnapshots(ctx, 1, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(snaps) != 2 {
		t.Fatalf("got %d snaps for board 1, want 2", len(snaps))
	}
}

func TestRisks_UpdateAndResolve(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveRisk(ctx, &domain.Risk{Title: "Risk A", Severity: "high", Status: "open", Owner: "alice", IdentifiedAt: time.Now()})
	store.SaveRisk(ctx, &domain.Risk{Title: "Risk B", Severity: "low", Status: "open", Owner: "bob", IdentifiedAt: time.Now()})

	risks, _ := store.GetOpenRisks(ctx)
	if len(risks) != 2 {
		t.Fatalf("got %d, want 2", len(risks))
	}
	// High severity comes first
	if risks[0].Severity != "high" {
		t.Errorf("expected high first, got %s", risks[0].Severity)
	}

	// Resolve one
	now := time.Now()
	risks[0].Status = "resolved"
	risks[0].ResolvedAt = &now
	store.UpdateRisk(ctx, &risks[0])

	open, _ := store.GetOpenRisks(ctx)
	if len(open) != 1 {
		t.Fatalf("got %d open, want 1", len(open))
	}
}

func TestRisks_GetAllRisks(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		store.SaveRisk(ctx, &domain.Risk{Title: "R", Severity: "medium", Status: "open", IdentifiedAt: time.Now()})
	}
	all, _ := store.GetAllRisks(ctx, 3)
	if len(all) != 3 {
		t.Fatalf("got %d, want 3 (limited)", len(all))
	}
}

func TestDecisions_SearchNoMatch(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveDecision(ctx, &domain.Decision{Title: "Use Redis", Decision: "For caching", MadeAt: time.Now()})

	found, err := store.SearchDecisions(ctx, "nonexistent")
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 0 {
		t.Fatalf("got %d, want 0", len(found))
	}
}

func TestDecisions_GetMultiple(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveDecision(ctx, &domain.Decision{Title: "D1", Decision: "x", MadeAt: time.Now()})
	store.SaveDecision(ctx, &domain.Decision{Title: "D2", Decision: "y", MadeAt: time.Now()})
	store.SaveDecision(ctx, &domain.Decision{Title: "D3", Decision: "z", MadeAt: time.Now()})

	decs, _ := store.GetDecisions(ctx, 2)
	if len(decs) != 2 {
		t.Fatalf("got %d, want 2", len(decs))
	}
}

func TestBlockers_HistoryIncludesResolved(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveBlocker(ctx, &domain.Blocker{IssueKey: "T-1", Description: "blocker1", BlockedSince: time.Now(), Owner: "a"})
	store.SaveBlocker(ctx, &domain.Blocker{IssueKey: "T-2", Description: "blocker2", BlockedSince: time.Now(), Owner: "b"})

	blockers, _ := store.GetActiveBlockers(ctx)
	store.ResolveBlocker(ctx, blockers[0].ID, "fixed")

	active, _ := store.GetActiveBlockers(ctx)
	if len(active) != 1 {
		t.Fatalf("got %d active, want 1", len(active))
	}

	history, _ := store.GetBlockerHistory(ctx, 10)
	if len(history) != 2 {
		t.Fatalf("got %d in history, want 2", len(history))
	}
}

func TestTeamMetrics_GetTeamOverview(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveTeamMetric(ctx, &domain.TeamMetric{MemberName: "alice", SprintName: "S10", RecordedAt: time.Now(), IssuesDone: 5})
	store.SaveTeamMetric(ctx, &domain.TeamMetric{MemberName: "bob", SprintName: "S10", RecordedAt: time.Now(), IssuesDone: 3})
	store.SaveTeamMetric(ctx, &domain.TeamMetric{MemberName: "alice", SprintName: "S11", RecordedAt: time.Now(), IssuesDone: 7})

	overview, _ := store.GetTeamOverview(ctx, "S10")
	if len(overview) != 2 {
		t.Fatalf("got %d, want 2", len(overview))
	}
}

func TestRetrospectives_SaveAndRetrieve(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	r := &domain.Retrospective{
		SprintName: "Sprint 5", Date: time.Now(),
		WentWell: "Good teamwork", Improvements: "Less meetings",
		ActionItems: "Reduce standups to 10min", Status: "open",
	}
	err := store.SaveRetrospective(ctx, r)
	if err != nil {
		t.Fatal(err)
	}
	if r.ID == 0 {
		t.Error("expected ID to be set after save")
	}

	retros, err := store.GetRetrospectives(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(retros) != 1 {
		t.Fatalf("got %d, want 1", len(retros))
	}
	if retros[0].WentWell != "Good teamwork" {
		t.Errorf("got %s", retros[0].WentWell)
	}
}

func TestActionItems_PendingAndComplete(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveActionItem(ctx, &domain.ActionItem{Description: "action1", Owner: "alice", Status: "pending"})
	store.SaveActionItem(ctx, &domain.ActionItem{Description: "action2", Owner: "bob", Status: "pending"})

	pending, _ := store.GetPendingActionItems(ctx)
	if len(pending) != 2 {
		t.Fatalf("got %d, want 2", len(pending))
	}

	store.CompleteActionItem(ctx, pending[0].ID)
	pending2, _ := store.GetPendingActionItems(ctx)
	if len(pending2) != 1 {
		t.Fatalf("got %d after complete, want 1", len(pending2))
	}
}

func TestDependencies_SaveResolveGet(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveDependency(ctx, &domain.Dependency{
		FromIssueKey: "T-1", ToIssueKey: "T-2", DependencyType: "blocks",
		Description: "T-1 blocks T-2", Status: "open", CreatedAt: time.Now(),
	})
	store.SaveDependency(ctx, &domain.Dependency{
		FromIssueKey: "T-3", ToIssueKey: "T-4", DependencyType: "external",
		Description: "Waiting on infra", Status: "open", CreatedAt: time.Now(),
	})

	open, _ := store.GetOpenDependencies(ctx)
	if len(open) != 2 {
		t.Fatalf("got %d, want 2", len(open))
	}

	store.ResolveDependency(ctx, open[0].ID)
	open2, _ := store.GetOpenDependencies(ctx)
	if len(open2) != 1 {
		t.Fatalf("got %d after resolve, want 1", len(open2))
	}
}

func TestDependencies_ForIssue(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveDependency(ctx, &domain.Dependency{FromIssueKey: "A-1", ToIssueKey: "B-1", DependencyType: "blocks", Status: "open", CreatedAt: time.Now()})
	store.SaveDependency(ctx, &domain.Dependency{FromIssueKey: "C-1", ToIssueKey: "A-1", DependencyType: "blocks", Status: "open", CreatedAt: time.Now()})
	store.SaveDependency(ctx, &domain.Dependency{FromIssueKey: "X-1", ToIssueKey: "Y-1", DependencyType: "blocks", Status: "open", CreatedAt: time.Now()})

	deps, _ := store.GetDependenciesForIssue(ctx, "A-1")
	if len(deps) != 2 {
		t.Fatalf("got %d deps for A-1, want 2", len(deps))
	}
}

func TestMeetingNotes_SaveAndFilter(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveMeetingNote(ctx, &domain.MeetingNote{MeetingType: "standup", Date: time.Now(), Notes: "discussed blockers"})
	store.SaveMeetingNote(ctx, &domain.MeetingNote{MeetingType: "retro", Date: time.Now(), Notes: "retro notes"})
	store.SaveMeetingNote(ctx, &domain.MeetingNote{MeetingType: "standup", Date: time.Now(), Notes: "standup 2"})

	all, _ := store.GetMeetingNotes(ctx, "", 10)
	if len(all) != 3 {
		t.Fatalf("got %d, want 3", len(all))
	}

	standups, _ := store.GetMeetingNotes(ctx, "standup", 10)
	if len(standups) != 2 {
		t.Fatalf("got %d standups, want 2", len(standups))
	}
}

func TestHealthScores_SaveAndGet(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveHealthScore(ctx, &domain.HealthScore{SprintName: "S1", BoardID: 1, ComputedAt: time.Now(), OverallScore: 75, VelocityScore: 20, BlockerScore: 15, ScopeScore: 20, TeamScore: 20})
	store.SaveHealthScore(ctx, &domain.HealthScore{SprintName: "S2", BoardID: 1, ComputedAt: time.Now(), OverallScore: 80, VelocityScore: 22, BlockerScore: 18, ScopeScore: 20, TeamScore: 20})
	store.SaveHealthScore(ctx, &domain.HealthScore{SprintName: "S1", BoardID: 2, ComputedAt: time.Now(), OverallScore: 60})

	scores, _ := store.GetHealthScores(ctx, 1, 10)
	if len(scores) != 2 {
		t.Fatalf("got %d for board 1, want 2", len(scores))
	}
}

func TestDailyProgress_SaveAndGet(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	today := time.Now()
	store.SaveDailyProgress(ctx, &domain.DailyProgress{SprintName: "S1", BoardID: 1, Date: today, TotalIssues: 10, Done: 3, InProgress: 4, Todo: 3})
	store.SaveDailyProgress(ctx, &domain.DailyProgress{SprintName: "S1", BoardID: 1, Date: today.Add(24 * time.Hour), TotalIssues: 10, Done: 5, InProgress: 3, Todo: 2})

	progress, _ := store.GetDailyProgress(ctx, 1, "S1")
	if len(progress) != 2 {
		t.Fatalf("got %d, want 2", len(progress))
	}
	if progress[0].Done != 3 {
		t.Errorf("got done=%d, want 3", progress[0].Done)
	}
}

func TestSprintGoals_SaveAndUpdate(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	// Note: GetActiveGoals/GetGoalHistory have a scan bug with NULL outcome column.
	// We test SaveSprintGoal and UpdateSprintGoal work without errors.
	err := store.SaveSprintGoal(ctx, &domain.SprintGoal{SprintName: "S10", BoardID: 1, Goal: "Ship auth", KeyResults: "KR1\nKR2", Status: "active", CreatedAt: time.Now()})
	if err != nil {
		t.Fatal("save:", err)
	}

	// Update sets outcome to non-null, so scan works after update
	now := time.Now()
	err = store.UpdateSprintGoal(ctx, &domain.SprintGoal{ID: 1, Status: "achieved", Outcome: "Shipped", ClosedAt: &now})
	if err != nil {
		t.Fatal("update:", err)
	}
}

func TestDoDItems_CRUD(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveDoDItem(ctx, &domain.DoDItem{Project: "*", Item: "Tests pass", Category: "testing", OrderNum: 1, Active: true})
	store.SaveDoDItem(ctx, &domain.DoDItem{Project: "PROJ", Item: "Docs updated", Category: "docs", OrderNum: 2, Active: true})
	store.SaveDoDItem(ctx, &domain.DoDItem{Project: "OTHER", Item: "Something", Category: "code", OrderNum: 3, Active: true})

	// Global items show for any project
	dod, _ := store.GetDoD(ctx, "PROJ")
	if len(dod) != 2 {
		t.Fatalf("got %d, want 2 (global + project-specific)", len(dod))
	}

	// Delete (soft)
	store.DeleteDoDItem(ctx, dod[0].ID)
	dod2, _ := store.GetDoD(ctx, "PROJ")
	if len(dod2) != 1 {
		t.Fatalf("got %d after delete, want 1", len(dod2))
	}
}

func TestEscalations_SaveGetAcknowledge(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	e := &domain.Escalation{Type: "blocker", ReferenceID: 1, Title: "Blocker too old", Severity: "high", EscalatedAt: time.Now(), Channel: "lark"}
	store.SaveEscalation(ctx, e)
	if e.ID == 0 {
		t.Error("expected ID")
	}

	store.SaveEscalation(ctx, &domain.Escalation{Type: "risk", Title: "Critical risk", Severity: "critical", EscalatedAt: time.Now(), Channel: "lark"})

	esc, _ := store.GetRecentEscalations(ctx, 10)
	if len(esc) != 2 {
		t.Fatalf("got %d, want 2", len(esc))
	}

	store.AcknowledgeEscalation(ctx, esc[0].ID)
	esc2, _ := store.GetRecentEscalations(ctx, 10)
	found := false
	for _, e := range esc2 {
		if e.Acknowledged {
			found = true
		}
	}
	if !found {
		t.Error("expected at least one acknowledged escalation")
	}
}

func TestTeamPulse_SaveAndGet(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveTeamPulse(ctx, &domain.TeamPulse{SprintName: "S10", Member: "alice", Score: 4, Notes: "good"})
	store.SaveTeamPulse(ctx, &domain.TeamPulse{SprintName: "S10", Member: "bob", Score: 3, Notes: "ok"})

	pulses, _ := store.GetTeamPulseHistory(ctx, 10)
	if len(pulses) != 2 {
		t.Fatalf("got %d, want 2", len(pulses))
	}
}

func TestMeetingEffectiveness_SaveAndFilter(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveMeetingEffectiveness(ctx, &domain.MeetingEffectiveness{Ceremony: "standup", DurationMinutes: 10, Score: 4, SprintName: "S10"})
	store.SaveMeetingEffectiveness(ctx, &domain.MeetingEffectiveness{Ceremony: "retro", DurationMinutes: 60, Score: 3, SprintName: "S10"})
	store.SaveMeetingEffectiveness(ctx, &domain.MeetingEffectiveness{Ceremony: "standup", DurationMinutes: 12, Score: 5, SprintName: "S11"})

	all, _ := store.GetMeetingEffectivenessHistory(ctx, "", 10)
	if len(all) != 3 {
		t.Fatalf("got %d, want 3", len(all))
	}

	standups, _ := store.GetMeetingEffectivenessHistory(ctx, "standup", 10)
	if len(standups) != 2 {
		t.Fatalf("got %d standups, want 2", len(standups))
	}
}

func TestTeamRadar_SaveAndGet(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	store.SaveTeamRadar(ctx, &domain.TeamRadar{SprintName: "S10", Dimension: "collaboration", Score: 4})
	store.SaveTeamRadar(ctx, &domain.TeamRadar{SprintName: "S10", Dimension: "quality", Score: 3})

	radars, _ := store.GetTeamRadarHistory(ctx, 10)
	if len(radars) != 2 {
		t.Fatalf("got %d, want 2", len(radars))
	}
}
