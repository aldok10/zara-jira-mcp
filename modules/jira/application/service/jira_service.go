// Package service implements application use cases for the jira module.
package service

import (
	"context"

	"github.com/aldok10/zara-jira-mcp/modules/jira/application/port"
	"github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// Ensure service implements the interface at compile time.
var _ port.Inbound = (*JiraService)(nil)

// JiraService implements port.Inbound for Jira operations.
type JiraService struct {
	client domain.Client
	cache  port.Cache
}

// NewJiraService creates a new JiraService with its dependencies.
func NewJiraService(client domain.Client, cache port.Cache) *JiraService {
	return &JiraService{client: client, cache: cache}
}

func (s *JiraService) SearchIssues(ctx context.Context, jql string, maxResults int) (*domain.SearchResult, error) {
	if maxResults <= 0 {
		maxResults = 50
	}
	return s.client.SearchIssues(ctx, jql, maxResults, 0)
}

func (s *JiraService) GetIssue(ctx context.Context, key string) (*domain.Issue, error) {
	return s.client.GetIssue(ctx, key)
}

func (s *JiraService) GetBoards(ctx context.Context) ([]domain.Board, error) {
	return s.client.GetBoards(ctx)
}

func (s *JiraService) GetActiveSprints(ctx context.Context, boardID int) ([]domain.Sprint, error) {
	return s.client.GetActiveSprints(ctx, boardID)
}

func (s *JiraService) GetSprintIssues(ctx context.Context, sprintID int) ([]domain.Issue, error) {
	return s.client.GetSprintIssues(ctx, sprintID)
}

func (s *JiraService) CreateIssue(ctx context.Context, input *domain.CreateIssueInput) (*domain.Issue, error) {
	return s.client.CreateIssue(ctx, input)
}

func (s *JiraService) UpdateIssue(ctx context.Context, input *domain.UpdateIssueInput) error {
	return s.client.UpdateIssue(ctx, input)
}

func (s *JiraService) AddComment(ctx context.Context, issueKey, body string) error {
	return s.client.AddComment(ctx, issueKey, body)
}

func (s *JiraService) DeleteIssue(ctx context.Context, issueKey string) error {
	return s.client.DeleteIssue(ctx, issueKey)
}

func (s *JiraService) CreateSprint(ctx context.Context, boardID int, name, goal string) (*domain.Sprint, error) {
	return s.client.CreateSprint(ctx, boardID, name, goal)
}

func (s *JiraService) CloseSprint(ctx context.Context, sprintID int) error {
	return s.client.CloseSprint(ctx, sprintID)
}
