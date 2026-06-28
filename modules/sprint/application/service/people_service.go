// Package people provides team health, workload, and collaboration use cases.
package service

import (
	"context"
	"fmt"

	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/peopleport"
	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/port"
)

// PeopleService defines people/team related operations.
type PeopleService interface {
	// WorkloadCheck analyzes team workload distribution.
	WorkloadCheck(ctx context.Context, boardID int) (string, error)

	// OverloadCheck identifies overloaded team members.
	OverloadCheck(ctx context.Context, boardID int) (string, error)

	// CollaborationSignal evaluates cross-area collaboration.
	CollaborationSignal(ctx context.Context, boardID int) (string, error)

	// ExecutiveReport generates an executive-friendly status report.
	ExecutiveReport(ctx context.Context, boardID int) (string, error)
}

var _ PeopleService = (*peopleService)(nil)

type peopleService struct {
	metrics  peopleport.TeamMetricRepository
	health   peopleport.HealthRepository
	blockers peopleport.BlockerRepository
	jira     peopleport.JiraClient
	ai       port.AIProvider
}

func NewPeopleService(
	metrics peopleport.TeamMetricRepository,
	health peopleport.HealthRepository,
	blockers peopleport.BlockerRepository,
	jira peopleport.JiraClient,
	ai port.AIProvider,
) PeopleService {
	return &peopleService{
		metrics:  metrics,
		health:   health,
		blockers: blockers,
		jira:     jira,
		ai:       ai,
	}
}

// WorkloadCheck implements PeopleService.WorkloadCheck.
func (p *peopleService) WorkloadCheck(ctx context.Context, boardID int) (string, error) {
	// TODO: Implement workload check
	return "", fmt.Errorf("not implemented")
}

// OverloadCheck implements PeopleService.OverloadCheck.
func (p *peopleService) OverloadCheck(ctx context.Context, boardID int) (string, error) {
	// TODO: Implement overload check
	return "", fmt.Errorf("not implemented")
}

// CollaborationSignal implements PeopleService.CollaborationSignal.
func (p *peopleService) CollaborationSignal(ctx context.Context, boardID int) (string, error) {
	// TODO: Implement collaboration signal
	return "", fmt.Errorf("not implemented")
}

// ExecutiveReport implements PeopleService.ExecutiveReport.
func (p *peopleService) ExecutiveReport(ctx context.Context, boardID int) (string, error) {
	// TODO: Implement executive report
	return "", fmt.Errorf("not implemented")
}
