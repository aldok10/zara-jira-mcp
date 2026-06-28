package port

import (
	"context"

	"github.com/aldok10/zara-jira-mcp/shared/kernel/event"
	jira "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
	memory "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
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

// Inbound defines the sprint use cases exposed by this module.
type Inbound interface {
	CalculateHealth(ctx context.Context, boardID int) (*HealthResult, error)
	Forecast(ctx context.Context, boardID int, remaining int) (*ForecastResult, error)
	DetectAntiPatterns(ctx context.Context, boardID int) ([]AntiPattern, error)
	VelocityTrend(ctx context.Context, boardID int) (string, error)
}

// HealthResult holds a sprint health assessment.
type HealthResult struct {
	Score         int    // 0-100
	Rating        string // Healthy/Fair/At Risk/Critical
	WeakestDim    string
	SprintName    string
	VelocityScore int
	BlockerScore  int
	ScopeScore    int
	TeamScore     int
}

// AntiPattern describes a detected anti-pattern.
type AntiPattern struct {
	Name        string
	Description string
	Severity    string // High/Medium/Low
	Suggestion  string
}

// ForecastResult holds Monte Carlo forecast results.
type ForecastResult struct {
	MeanSprints float64
	MinSprints  int
	MaxSprints  int
	Percentiles map[int]float64
	Remaining   int
	Simulations int
}
