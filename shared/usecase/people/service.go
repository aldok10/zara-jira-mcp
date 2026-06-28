// Package people provides team health, workload, and collaboration use cases.
package people

import (
	"context"
	"fmt"
)

// Service defines people/team related operations.
type Service interface {
	// WorkloadCheck analyzes team workload distribution.
	WorkloadCheck(ctx context.Context, boardID int) (string, error)

	// OverloadCheck identifies overloaded team members.
	OverloadCheck(ctx context.Context, boardID int) (string, error)

	// CollaborationSignal evaluates cross-area collaboration.
	CollaborationSignal(ctx context.Context, boardID int) (string, error)

	// ExecutiveReport generates an executive-friendly status report.
	ExecutiveReport(ctx context.Context, boardID int) (string, error)
}

var _ Service = (*peopleService)(nil)

type peopleService struct {
	metrics  TeamMetricRepository
	health   HealthRepository
	blockers BlockerRepository
	jira     JiraClient
	ai       AIProvider
}

func NewPeopleService(
	metrics TeamMetricRepository,
	health HealthRepository,
	blockers BlockerRepository,
	jira JiraClient,
	ai AIProvider,
) Service {
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
