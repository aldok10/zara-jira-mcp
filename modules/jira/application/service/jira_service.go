// Package service implements application use cases for the jira module.
package service

import (
	"context"
	"sync"

	"github.com/aldok10/zara-jira-mcp/modules/jira/application/port"
	"github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// Ensure service implements the interface at compile time.
var _ port.Inbound = (*JiraService)(nil)

// JiraService implements port.Inbound for Jira operations.
type JiraService struct {
	client domain.Client
	mu     sync.RWMutex
	boards []domain.Board // cached board list
}

// NewJiraService creates a new JiraService with its dependencies.
func NewJiraService(client domain.Client) *JiraService {
	return &JiraService{client: client}
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

// GetBoards returns all accessible boards, caching the result for future calls.
func (s *JiraService) GetBoards(ctx context.Context) ([]domain.Board, error) {
	// Fast path: return cached copy
	s.mu.RLock()
	if s.boards != nil {
		cpy := make([]domain.Board, len(s.boards))
		copy(cpy, s.boards)
		s.mu.RUnlock()
		return cpy, nil
	}
	s.mu.RUnlock()

	// Slow path: fetch and cache
	boards, err := s.client.GetBoards(ctx)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	s.boards = make([]domain.Board, len(boards))
	copy(s.boards, boards)
	s.mu.Unlock()

	cpy := make([]domain.Board, len(boards))
	copy(cpy, boards)
	return cpy, nil
}

// GetDefaultBoardID returns the first available board ID.
// Returns 0 if no boards found.
func (s *JiraService) GetDefaultBoardID(ctx context.Context) (int, error) {
	boards, err := s.GetBoards(ctx)
	if err != nil {
		return 0, err
	}
	if len(boards) == 0 {
		return 0, nil
	}
	return boards[0].ID, nil
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

// InvalidateBoardCache clears the cached board list, forcing a fresh fetch on next call.
func (s *JiraService) InvalidateBoardCache() {
	s.mu.Lock()
	s.boards = nil
	s.mu.Unlock()
}
