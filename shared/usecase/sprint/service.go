// Package sprint provides sprint health, forecasting, and analysis use cases.
package sprint

import (
	"context"
	"fmt"

	"github.com/aldok10/zara-jira-mcp/shared/domain/planning"
)

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

// Service defines sprint-related operations.
type Service interface {
	// CalculateHealth computes a 0-100 health score for the active sprint.
	CalculateHealth(ctx context.Context, boardID int) (*HealthResult, error)

	// Forecast predicts completion sprints using Monte Carlo simulation.
	Forecast(ctx context.Context, boardID int, remaining int) (*ForecastResult, error)

	// DetectAntiPatterns scans for known anti-patterns in sprint execution.
	DetectAntiPatterns(ctx context.Context, boardID int) ([]AntiPattern, error)

	// VelocityTrend returns the velocity direction over recent sprints.
	VelocityTrend(ctx context.Context, boardID int) (string, error)
}

// Ensure service implements the interface at compile time.
var _ Service = (*sprintService)(nil)

// sprintService implements Service.
type sprintService struct {
	snapshots SnapshotRepository
	health    HealthRepository
	risks     RiskRepository
	blockers  BlockerRepository
	goals     GoalRepository
	jira      JiraClient
	ai        AIProvider
	events    EventBus
}

// CalculateHealth computes a 0-100 health score for the active sprint.
func (s *sprintService) CalculateHealth(ctx context.Context, boardID int) (*HealthResult, error) {
	// TODO: Implement health calculation
	return nil, fmt.Errorf("not implemented")
}

// Forecast predicts completion sprints using Monte Carlo simulation.
func (s *sprintService) Forecast(ctx context.Context, boardID int, remaining int) (*ForecastResult, error) {
	snaps, err := s.snapshots.FindByBoard(ctx, boardID, 20)
	if err != nil {
		return nil, fmt.Errorf("fetch snapshots: %w", err)
	}
	if len(snaps) < 3 {
		return nil, fmt.Errorf("need at least 3 sprint snapshots, got %d", len(snaps))
	}

	throughput := make([]float64, len(snaps))
	for i := range snaps {
		if snaps[i].Done <= 0 {
			throughput[i] = 1 // floor of 1 for simulation stability
		} else {
			throughput[i] = float64(snaps[i].Done)
		}
	}

	res := planning.Forecast(throughput, remaining, 0)

	return &ForecastResult{
		MeanSprints: res.MeanSprints,
		MinSprints:  res.MinSprints,
		MaxSprints:  res.MaxSprints,
		Percentiles: res.Percentiles,
		Remaining:   res.Remaining,
		Simulations: res.Simulations,
	}, nil
}

// DetectAntiPatterns scans for known anti-patterns in sprint execution.
func (s *sprintService) DetectAntiPatterns(ctx context.Context, boardID int) ([]AntiPattern, error) {
	// TODO: Implement anti-pattern detection
	return nil, fmt.Errorf("not implemented")
}

// VelocityTrend returns the velocity direction over recent sprints.
func (s *sprintService) VelocityTrend(ctx context.Context, boardID int) (string, error) {
	// TODO: Implement velocity trend
	return "", fmt.Errorf("not implemented")
}

// NewSprintService creates a new Service with its dependencies.
func NewSprintService(
	snapshots SnapshotRepository,
	health HealthRepository,
	risks RiskRepository,
	blockers BlockerRepository,
	goals GoalRepository,
	jira JiraClient,
	ai AIProvider,
	events EventBus,
) Service {
	return &sprintService{
		snapshots: snapshots,
		health:    health,
		risks:     risks,
		blockers:  blockers,
		goals:     goals,
		jira:      jira,
		ai:        ai,
		events:    events,
	}
}
