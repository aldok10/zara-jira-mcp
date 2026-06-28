// Package sprintstore provides adapter implementations of sprint.SnapshotRepository
// and other focused port interfaces, bridging internal implementations to use case ports.
package sprintstore

import (
	"context"

	imemory "github.com/aldok10/zara-jira-mcp/internal/memory"
	domain "github.com/aldok10/zara-jira-mcp/shared/domain/memory"
)

// SnapshotRepository adapts *memory.SQLiteStore to the sprint.SnapshotRepository interface.
type SnapshotRepository struct {
	store *imemory.SQLiteStore
}

// NewSnapshotRepository wraps a SQLiteStore as a SnapshotRepository.
func NewSnapshotRepository(store *imemory.SQLiteStore) *SnapshotRepository {
	return &SnapshotRepository{store: store}
}

func (r *SnapshotRepository) Save(ctx context.Context, s *domain.SprintSnapshot) error {
	return r.store.SaveSprintSnapshot(ctx, s)
}

func (r *SnapshotRepository) FindByBoard(ctx context.Context, boardID int, limit int) ([]domain.SprintSnapshot, error) {
	return r.store.GetSprintSnapshots(ctx, boardID, limit)
}

func (r *SnapshotRepository) FindLatest(ctx context.Context, boardID int) (*domain.SprintSnapshot, error) {
	snaps, err := r.store.GetSprintSnapshots(ctx, boardID, 1)
	if err != nil {
		return nil, err
	}
	if len(snaps) == 0 {
		return nil, nil
	}
	return &snaps[0], nil
}
