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

func (s *JiraService) TransitionIssue(ctx context.Context, issueKey, transitionID string) error {
	return s.client.TransitionIssue(ctx, issueKey, transitionID)
}

func (s *JiraService) GetTransitions(ctx context.Context, issueKey string) ([]domain.Transition, error) {
	return s.client.GetTransitions(ctx, issueKey)
}

func (s *JiraService) AssignIssue(ctx context.Context, issueKey, accountID string) error {
	return s.client.AssignIssue(ctx, issueKey, accountID)
}

func (s *JiraService) FindUser(ctx context.Context, query string) ([]domain.User, error) {
	return s.client.FindUser(ctx, query)
}

func (s *JiraService) GetSprints(ctx context.Context, boardID int, state string) ([]domain.Sprint, error) {
	return s.client.GetSprints(ctx, boardID, state)
}

func (s *JiraService) StartSprint(ctx context.Context, sprintID int, startDate, endDate string) error {
	return s.client.StartSprint(ctx, sprintID, startDate, endDate)
}

func (s *JiraService) MoveIssuesToSprint(ctx context.Context, sprintID int, issueKeys []string) error {
	return s.client.MoveIssuesToSprint(ctx, sprintID, issueKeys)
}

func (s *JiraService) CreateSubtask(ctx context.Context, parentKey string, input *domain.CreateIssueInput) (*domain.Issue, error) {
	return s.client.CreateSubtask(ctx, parentKey, input)
}

func (s *JiraService) LinkIssues(ctx context.Context, inwardKey, outwardKey, linkType string) error {
	return s.client.LinkIssues(ctx, inwardKey, outwardKey, linkType)
}

func (s *JiraService) GetLinkTypes(ctx context.Context) ([]domain.LinkType, error) {
	return s.client.GetLinkTypes(ctx)
}

func (s *JiraService) AddWorklog(ctx context.Context, issueKey, timeSpent, comment string) error {
	return s.client.AddWorklog(ctx, issueKey, timeSpent, comment)
}

func (s *JiraService) GetWorklogs(ctx context.Context, issueKey string) ([]domain.Worklog, error) {
	return s.client.GetWorklogs(ctx, issueKey)
}

func (s *JiraService) AddWatcher(ctx context.Context, issueKey, accountID string) error {
	return s.client.AddWatcher(ctx, issueKey, accountID)
}

func (s *JiraService) GetWatchers(ctx context.Context, issueKey string) ([]domain.User, error) {
	return s.client.GetWatchers(ctx, issueKey)
}

func (s *JiraService) AddLabel(ctx context.Context, issueKey, label string) error {
	return s.client.AddLabel(ctx, issueKey, label)
}

func (s *JiraService) GetProjects(ctx context.Context) ([]domain.Project, error) {
	return s.client.GetProjects(ctx)
}

func (s *JiraService) GetProject(ctx context.Context, key string) (*domain.ProjectDetail, error) {
	return s.client.GetProject(ctx, key)
}

func (s *JiraService) GetVersions(ctx context.Context, projectKey string) ([]domain.Version, error) {
	return s.client.GetVersions(ctx, projectKey)
}

func (s *JiraService) CreateVersion(ctx context.Context, projectKey, name, description string) (*domain.Version, error) {
	return s.client.CreateVersion(ctx, projectKey, name, description)
}

func (s *JiraService) ReleaseVersion(ctx context.Context, versionID string) error {
	return s.client.ReleaseVersion(ctx, versionID)
}

func (s *JiraService) GetAttachments(ctx context.Context, issueKey string) ([]domain.Attachment, error) {
	return s.client.GetAttachments(ctx, issueKey)
}

func (s *JiraService) GetComponents(ctx context.Context, projectKey string) ([]domain.Component, error) {
	return s.client.GetComponents(ctx, projectKey)
}

func (s *JiraService) GetFields(ctx context.Context) ([]domain.Field, error) {
	return s.client.GetFields(ctx)
}

// InvalidateBoardCache clears the cached board list, forcing a fresh fetch on next call.
func (s *JiraService) InvalidateBoardCache() {
	s.mu.Lock()
	s.boards = nil
	s.mu.Unlock()
}
