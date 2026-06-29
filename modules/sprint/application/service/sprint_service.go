// Package service implements sprint application service.
package service

import (
	"context"
	"fmt"

	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/port"
	"github.com/aldok10/zara-jira-mcp/modules/sprint/domain/planning"
)

// Ensure service implements the interface at compile time.
var _ port.Inbound = (*sprintService)(nil)

// sprintService implements port.Inbound.
type sprintService struct {
	snapshots port.SnapshotRepository
	health    port.HealthRepository
	risks     port.RiskRepository
	blockers  port.BlockerRepository
	goals     port.GoalRepository
	jira      port.JiraClient
	ai        port.AIProvider
	events    port.EventBus
}

// CalculateHealth computes a 0-100 health score for the active sprint.
func (s *sprintService) CalculateHealth(ctx context.Context, boardID int) (*port.HealthResult, error) {
	return nil, fmt.Errorf("not implemented")
}

// Forecast predicts completion sprints using Monte Carlo simulation.
func (s *sprintService) Forecast(ctx context.Context, boardID int, remaining int) (*port.ForecastResult, error) {
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
			throughput[i] = 1
		} else {
			throughput[i] = float64(snaps[i].Done)
		}
	}

	res := planning.Forecast(throughput, remaining, 0)

	return &port.ForecastResult{
		MeanSprints: res.MeanSprints,
		MinSprints:  res.MinSprints,
		MaxSprints:  res.MaxSprints,
		Percentiles: res.Percentiles,
		Remaining:   res.Remaining,
		Simulations: res.Simulations,
	}, nil
}

// DetectAntiPatterns scans for known anti-patterns in sprint execution.
func (s *sprintService) DetectAntiPatterns(ctx context.Context, boardID int) ([]port.AntiPattern, error) {
	return nil, fmt.Errorf("not implemented")
}

// VelocityTrend returns the velocity direction over recent sprints.
func (s *sprintService) VelocityTrend(ctx context.Context, boardID int) (string, error) {
	return "", fmt.Errorf("not implemented")
}

// NewSprintService creates a new port.Inbound with its dependencies.
func NewSprintService(
	snapshots port.SnapshotRepository,
	health port.HealthRepository,
	risks port.RiskRepository,
	blockers port.BlockerRepository,
	goals port.GoalRepository,
	jira port.JiraClient,
	ai port.AIProvider,
	events port.EventBus,
) port.Inbound {
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
