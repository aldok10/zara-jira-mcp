// Package okr provides OKR health tracking, suggestions, and sync use cases.
package service

import (
	"context"
	"fmt"

	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/okrport"
	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/port"
)

// HealthResult holds the health status for an OKR.
type HealthResult struct {
	Objective   string
	KeyResult   string
	Progress    float64
	Target      float64
	DaysElapsed int
	DaysTotal   int
	Status      string // ON_TRACK / NEEDS_ATTENTION / AT_RISK
}

// OKRService defines OKR-related operations.
type OKRService interface {
	// Health checks progress of all OKR signals.
	Health(ctx context.Context) ([]HealthResult, error)

	// Suggest recommends OKRs from sprint data.
	Suggest(ctx context.Context, boardID int) (string, error)

	// SyncWithLark bi-directionally syncs OKR progress with Lark.
	SyncWithLark(ctx context.Context) error
}

var _ OKRService = (*okrService)(nil)

type okrService struct {
	okrs    okrport.Repository
	sprints okrport.SprintRepository
	jira    okrport.JiraClient
	ai      port.AIProvider
}

func NewOKRService(
	okrs okrport.Repository,
	sprints okrport.SprintRepository,
	jira okrport.JiraClient,
	ai port.AIProvider,
) OKRService {
	return &okrService{
		okrs:    okrs,
		sprints: sprints,
		jira:    jira,
		ai:      ai,
	}
}

// Health implements Service.Health.
func (o *okrService) Health(ctx context.Context) ([]HealthResult, error) {
	// TODO: Implement OKR health check
	return nil, fmt.Errorf("not implemented")
}

// Suggest implements Service.Suggest.
func (o *okrService) Suggest(ctx context.Context, boardID int) (string, error) {
	// TODO: Implement OKR suggestion
	return "", fmt.Errorf("not implemented")
}

// SyncWithLark implements Service.SyncWithLark.
func (o *okrService) SyncWithLark(ctx context.Context) error {
	// TODO: Implement Lark sync
	return fmt.Errorf("not implemented")
}
