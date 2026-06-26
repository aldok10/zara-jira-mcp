package memory

import (
	"context"
	"os"
	"testing"
	"time"

	domain "github.com/aldok10/zara-jira-mcp/domain/memory"
)

func setupTestDB(t *testing.T) *SQLiteStore {
	t.Helper()
	f, err := os.CreateTemp("", "pm-test-*.db")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	store, err := NewSQLiteStore(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store
}

func TestSprintSnapshots(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	err := store.SaveSprintSnapshot(ctx, &domain.SprintSnapshot{
		SprintName: "Sprint 10", BoardID: 42, SnapshotDate: time.Now(),
		TotalIssues: 20, Done: 12, InProgress: 5, Todo: 3,
		Velocity: 32, CompletionRate: 0.6,
	})
	if err != nil {
		t.Fatal(err)
	}

	snaps, err := store.GetSprintSnapshots(ctx, 42, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(snaps) != 1 {
		t.Fatalf("got %d snapshots, want 1", len(snaps))
	}
	if snaps[0].SprintName != "Sprint 10" {
		t.Errorf("got %s, want Sprint 10", snaps[0].SprintName)
	}
	if snaps[0].Done != 12 {
		t.Errorf("got done=%d, want 12", snaps[0].Done)
	}
}

func TestRisks(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	err := store.SaveRisk(ctx, &domain.Risk{
		Title: "DB migration risk", Severity: "high", Status: "open",
		Owner: "aldo", IdentifiedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	risks, err := store.GetOpenRisks(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(risks) != 1 {
		t.Fatalf("got %d risks, want 1", len(risks))
	}
	if risks[0].Title != "DB migration risk" {
		t.Errorf("got %s", risks[0].Title)
	}
}

func TestDecisions(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	err := store.SaveDecision(ctx, &domain.Decision{
		Title: "Use PostgreSQL", Decision: "Chosen over MongoDB",
		Rationale: "Better for relational data", MadeBy: "team", MadeAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	decisions, err := store.GetDecisions(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(decisions) != 1 {
		t.Fatalf("got %d, want 1", len(decisions))
	}

	found, err := store.SearchDecisions(ctx, "PostgreSQL")
	if err != nil {
		t.Fatal(err)
	}
	if len(found) != 1 {
		t.Fatalf("search got %d, want 1", len(found))
	}
}

func TestBlockers(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	err := store.SaveBlocker(ctx, &domain.Blocker{
		IssueKey: "SIT-123", Description: "Waiting for API key",
		BlockedSince: time.Now(), Owner: "aldo",
	})
	if err != nil {
		t.Fatal(err)
	}

	blockers, err := store.GetActiveBlockers(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(blockers) != 1 {
		t.Fatalf("got %d, want 1", len(blockers))
	}

	err = store.ResolveBlocker(ctx, blockers[0].ID, "Got the key from infra team")
	if err != nil {
		t.Fatal(err)
	}

	active, _ := store.GetActiveBlockers(ctx)
	if len(active) != 0 {
		t.Errorf("got %d active after resolve, want 0", len(active))
	}
}

func TestTeamMetrics(t *testing.T) {
	store := setupTestDB(t)
	ctx := context.Background()

	err := store.SaveTeamMetric(ctx, &domain.TeamMetric{
		MemberName: "aldo", SprintName: "Sprint 10", RecordedAt: time.Now(),
		IssuesAssigned: 8, IssuesDone: 6, BlockerCount: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	metrics, err := store.GetTeamMetrics(ctx, "aldo", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(metrics) != 1 {
		t.Fatalf("got %d, want 1", len(metrics))
	}
	if metrics[0].IssuesDone != 6 {
		t.Errorf("got done=%d, want 6", metrics[0].IssuesDone)
	}
}
