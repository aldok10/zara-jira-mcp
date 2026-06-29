// Package sprintstore provides adapter implementations of sprint domain port interfaces.
package sprintstore

import (
	"context"

	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/port"
	memory "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
	"github.com/aldok10/zara-jira-mcp/modules/sprint/infrastructure/persistence"
	"github.com/aldok10/zara-jira-mcp/shared/domain/event"
)

// compile-time interface checks
var _ port.HealthRepository = (*HealthRepository)(nil)
var _ port.RiskRepository = (*RiskRepository)(nil)
var _ port.BlockerRepository = (*BlockerRepository)(nil)
var _ port.GoalRepository = (*GoalRepository)(nil)

// --- HealthRepository adapter ---

// HealthRepository adapts *persistence.SQLiteStore to port.HealthRepository.
type HealthRepository struct {
	store *persistence.SQLiteStore
}

// NewHealthRepository wraps a SQLiteStore as a HealthRepository.
func NewHealthRepository(store *persistence.SQLiteStore) *HealthRepository {
	return &HealthRepository{store: store}
}

func (r *HealthRepository) Save(ctx context.Context, h *memory.HealthScore) error {
	return r.store.SaveHealthScore(ctx, h)
}

func (r *HealthRepository) FindByBoard(ctx context.Context, boardID int, limit int) ([]memory.HealthScore, error) {
	return r.store.GetHealthScores(ctx, boardID, limit)
}

// --- RiskRepository adapter ---

// RiskRepository adapts *persistence.SQLiteStore to port.RiskRepository.
type RiskRepository struct {
	store *persistence.SQLiteStore
}

// NewRiskRepository wraps a SQLiteStore as a RiskRepository.
func NewRiskRepository(store *persistence.SQLiteStore) *RiskRepository {
	return &RiskRepository{store: store}
}

func (r *RiskRepository) FindOpen(ctx context.Context) ([]memory.Risk, error) {
	return r.store.GetOpenRisks(ctx)
}

func (r *RiskRepository) Save(ctx context.Context, risk *memory.Risk) error {
	return r.store.SaveRisk(ctx, risk)
}

// --- BlockerRepository adapter ---

// BlockerRepository adapts *persistence.SQLiteStore to port.BlockerRepository.
type BlockerRepository struct {
	store *persistence.SQLiteStore
}

// NewBlockerRepository wraps a SQLiteStore as a BlockerRepository.
func NewBlockerRepository(store *persistence.SQLiteStore) *BlockerRepository {
	return &BlockerRepository{store: store}
}

func (r *BlockerRepository) FindActive(ctx context.Context) ([]memory.Blocker, error) {
	return r.store.GetActiveBlockers(ctx)
}

func (r *BlockerRepository) Save(ctx context.Context, b *memory.Blocker) error {
	return r.store.SaveBlocker(ctx, b)
}

// --- GoalRepository adapter ---

// GoalRepository adapts *persistence.SQLiteStore to port.GoalRepository.
type GoalRepository struct {
	store *persistence.SQLiteStore
}

// NewGoalRepository wraps a SQLiteStore as a GoalRepository.
func NewGoalRepository(store *persistence.SQLiteStore) *GoalRepository {
	return &GoalRepository{store: store}
}

func (r *GoalRepository) FindActive(ctx context.Context, boardID int) ([]memory.SprintGoal, error) {
	return r.store.GetActiveGoals(ctx, boardID)
}

func (r *GoalRepository) Save(ctx context.Context, g *memory.SprintGoal) error {
	return r.store.SaveSprintGoal(ctx, g)
}

// --- NoopEventBus ---

// NoopEventBus implements port.EventBus as a no-op (for when events are not yet wired).
type NoopEventBus struct{}

func (b *NoopEventBus) Publish(ctx context.Context, e event.Event) error {
	return nil
}
