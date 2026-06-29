package sprint

import (
	"context"

	"github.com/aldok10/zara-jira-mcp/shared/domain/event"
	"github.com/aldok10/zara-jira-mcp/shared/domain/jira"
	memory "github.com/aldok10/zara-jira-mcp/shared/domain/memory"
)

// SnapshotRepository provides sprint snapshot persistence.
type SnapshotRepository interface {
	Save(ctx context.Context, s *memory.SprintSnapshot) error
	FindByBoard(ctx context.Context, boardID int, limit int) ([]memory.SprintSnapshot, error)
	FindLatest(ctx context.Context, boardID int) (*memory.SprintSnapshot, error)
}

// HealthRepository provides health score persistence.
type HealthRepository interface {
	Save(ctx context.Context, h *memory.HealthScore) error
	FindByBoard(ctx context.Context, boardID int, limit int) ([]memory.HealthScore, error)
}

// RiskRepository provides risk persistence.
type RiskRepository interface {
	FindOpen(ctx context.Context) ([]memory.Risk, error)
	Save(ctx context.Context, r *memory.Risk) error
}

// BlockerRepository provides blocker persistence.
type BlockerRepository interface {
	FindActive(ctx context.Context) ([]memory.Blocker, error)
	Save(ctx context.Context, b *memory.Blocker) error
}

// GoalRepository provides sprint goal persistence.
type GoalRepository interface {
	FindActive(ctx context.Context, boardID int) ([]memory.SprintGoal, error)
	Save(ctx context.Context, g *memory.SprintGoal) error
}

// JiraClient provides Jira data access for sprint operations.
type JiraClient interface {
	GetBoards(ctx context.Context) ([]jira.Board, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]jira.Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]jira.Issue, error)
	SearchIssues(ctx context.Context, jql string, maxResults int, startAt int) (*jira.SearchResult, error)
}

// AIProvider provides AI analysis for sprint insights.
type AIProvider interface {
	Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}

// EventBus publishes domain events.
type EventBus interface {
	Publish(ctx context.Context, event event.Event) error
}
