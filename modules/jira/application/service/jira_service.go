// Package service implements application use cases for the jira module.
package service

import (
	"context"
	"fmt"

	"github.com/aldok10/zara-jira-mcp/modules/jira/application/port"
	"github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// Ensure service implements the interface at compile time.
var _ port.Inbound = (*JiraService)(nil)

// JiraService implements port.Inbound for Jira operations.
type JiraService struct {
	client port.Client
	cache  port.Cache
}

// NewJiraService creates a new JiraService with its dependencies.
func NewJiraService(client port.Client, cache port.Cache) *JiraService {
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

// Ensure Client interface compliance at compile time.
var _ port.Client = (*clientAdapter)(nil)

// clientAdapter wraps domain.Client to implement port.Client.
type clientAdapter struct {
	inner domain.Client
}

// NewClientAdapter adapts a domain.Client to the application's Client interface.
func NewClientAdapter(inner domain.Client) port.Client {
	return &clientAdapter{inner: inner}
}

func (a *clientAdapter) SearchIssues(ctx context.Context, jql string, maxResults int, startAt int) (*domain.SearchResult, error) {
	return a.inner.SearchIssues(ctx, jql, maxResults, startAt)
}

func (a *clientAdapter) GetIssue(ctx context.Context, key string) (*domain.Issue, error) {
	return a.inner.GetIssue(ctx, key)
}

func (a *clientAdapter) GetBoards(ctx context.Context) ([]domain.Board, error) {
	return a.inner.GetBoards(ctx)
}

func (a *clientAdapter) GetActiveSprints(ctx context.Context, boardID int) ([]domain.Sprint, error) {
	return a.inner.GetActiveSprints(ctx, boardID)
}

func (a *clientAdapter) GetSprintIssues(ctx context.Context, sprintID int) ([]domain.Issue, error) {
	return a.inner.GetSprintIssues(ctx, sprintID)
}

func (a *clientAdapter) CreateIssue(ctx context.Context, input *domain.CreateIssueInput) (*domain.Issue, error) {
	return a.inner.CreateIssue(ctx, input)
}

func (a *clientAdapter) UpdateIssue(ctx context.Context, input *domain.UpdateIssueInput) error {
	if input == nil {
		return fmt.Errorf("update input is nil")
	}
	return a.inner.UpdateIssue(ctx, input)
}

func (a *clientAdapter) AddComment(ctx context.Context, issueKey, body string) error {
	return a.inner.AddComment(ctx, issueKey, body)
}

func (a *clientAdapter) TransitionIssue(ctx context.Context, issueKey, transitionID string) error {
	return a.inner.TransitionIssue(ctx, issueKey, transitionID)
}

func (a *clientAdapter) AssignIssue(ctx context.Context, issueKey, accountID string) error {
	return a.inner.AssignIssue(ctx, issueKey, accountID)
}

func (a *clientAdapter) DeleteIssue(ctx context.Context, issueKey string) error {
	return a.inner.DeleteIssue(ctx, issueKey)
}

func (a *clientAdapter) CreateSubtask(ctx context.Context, parentKey string, input *domain.CreateIssueInput) (*domain.Issue, error) {
	return a.inner.CreateSubtask(ctx, parentKey, input)
}

func (a *clientAdapter) GetSprints(ctx context.Context, boardID int, state string) ([]domain.Sprint, error) {
	return a.inner.GetSprints(ctx, boardID, state)
}

func (a *clientAdapter) CreateSprint(ctx context.Context, boardID int, name, goal string) (*domain.Sprint, error) {
	return a.inner.CreateSprint(ctx, boardID, name, goal)
}

func (a *clientAdapter) StartSprint(ctx context.Context, sprintID int, startDate, endDate string) error {
	return a.inner.StartSprint(ctx, sprintID, startDate, endDate)
}

func (a *clientAdapter) CloseSprint(ctx context.Context, sprintID int) error {
	return a.inner.CloseSprint(ctx, sprintID)
}

func (a *clientAdapter) MoveIssuesToSprint(ctx context.Context, sprintID int, issueKeys []string) error {
	return a.inner.MoveIssuesToSprint(ctx, sprintID, issueKeys)
}

func (a *clientAdapter) LinkIssues(ctx context.Context, inwardKey, outwardKey, linkType string) error {
	return a.inner.LinkIssues(ctx, inwardKey, outwardKey, linkType)
}

func (a *clientAdapter) GetLinkTypes(ctx context.Context) ([]domain.LinkType, error) {
	return a.inner.GetLinkTypes(ctx)
}

func (a *clientAdapter) AddWorklog(ctx context.Context, issueKey, timeSpent, comment string) error {
	return a.inner.AddWorklog(ctx, issueKey, timeSpent, comment)
}

func (a *clientAdapter) GetWorklogs(ctx context.Context, issueKey string) ([]domain.Worklog, error) {
	return a.inner.GetWorklogs(ctx, issueKey)
}

func (a *clientAdapter) AddWatcher(ctx context.Context, issueKey, accountID string) error {
	return a.inner.AddWatcher(ctx, issueKey, accountID)
}

func (a *clientAdapter) GetWatchers(ctx context.Context, issueKey string) ([]domain.User, error) {
	return a.inner.GetWatchers(ctx, issueKey)
}

func (a *clientAdapter) AddLabel(ctx context.Context, issueKey, label string) error {
	return a.inner.AddLabel(ctx, issueKey, label)
}

func (a *clientAdapter) GetProjects(ctx context.Context) ([]domain.Project, error) {
	return a.inner.GetProjects(ctx)
}

func (a *clientAdapter) GetProject(ctx context.Context, key string) (*domain.ProjectDetail, error) {
	return a.inner.GetProject(ctx, key)
}

func (a *clientAdapter) RawRequest(ctx context.Context, method, path string, body []byte) ([]byte, int, error) {
	return a.inner.RawRequest(ctx, method, path, body)
}

func (a *clientAdapter) GetAttachments(ctx context.Context, issueKey string) ([]domain.Attachment, error) {
	return a.inner.GetAttachments(ctx, issueKey)
}

func (a *clientAdapter) GetVersions(ctx context.Context, projectKey string) ([]domain.Version, error) {
	return a.inner.GetVersions(ctx, projectKey)
}

func (a *clientAdapter) CreateVersion(ctx context.Context, projectKey, name, description string) (*domain.Version, error) {
	return a.inner.CreateVersion(ctx, projectKey, name, description)
}

func (a *clientAdapter) ReleaseVersion(ctx context.Context, versionID string) error {
	return a.inner.ReleaseVersion(ctx, versionID)
}

func (a *clientAdapter) GetComponents(ctx context.Context, projectKey string) ([]domain.Component, error) {
	return a.inner.GetComponents(ctx, projectKey)
}

func (a *clientAdapter) GetFields(ctx context.Context) ([]domain.Field, error) {
	return a.inner.GetFields(ctx)
}
