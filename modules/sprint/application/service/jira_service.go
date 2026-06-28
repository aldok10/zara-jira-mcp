// Package jira provides Jira CRUD and query use cases.
package service

import (
	"context"
	"fmt"

	jira "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/jiraport"
)

// JiraService defines standard Jira operations.
type JiraService interface {
	SearchIssues(ctx context.Context, jql string, maxResults int) (*jira.SearchResult, error)
	GetIssue(ctx context.Context, key string) (*jira.Issue, error)
	GetBoards(ctx context.Context) ([]jira.Board, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]jira.Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]jira.Issue, error)
	CreateIssue(ctx context.Context, input *jira.CreateIssueInput) (*jira.Issue, error)
	UpdateIssue(ctx context.Context, input *jira.UpdateIssueInput) error
	AddComment(ctx context.Context, issueKey, body string) error
	DeleteIssue(ctx context.Context, issueKey string) error
	CreateSprint(ctx context.Context, boardID int, name, goal string) (*jira.Sprint, error)
	CloseSprint(ctx context.Context, sprintID int) error
}

var _ JiraService = (*jiraService)(nil)

type jiraService struct {
	client jiraport.Client
	cache  jiraport.Cache
}

func NewJiraService(client jiraport.Client, cache jiraport.Cache) JiraService {
	return &jiraService{client: client, cache: cache}
}

// SearchIssues implements Service.SearchIssues.
func (j *jiraService) SearchIssues(ctx context.Context, jql string, maxResults int) (*jira.SearchResult, error) {
	// TODO: Implement search with caching
	return nil, fmt.Errorf("not implemented")
}

// GetIssue implements Service.GetIssue.
func (j *jiraService) GetIssue(ctx context.Context, key string) (*jira.Issue, error) {
	// TODO: Implement get issue
	return nil, fmt.Errorf("not implemented")
}

// GetBoards implements Service.GetBoards.
func (j *jiraService) GetBoards(ctx context.Context) ([]jira.Board, error) {
	// TODO: Implement get boards
	return nil, fmt.Errorf("not implemented")
}

// GetActiveSprints implements Service.GetActiveSprints.
func (j *jiraService) GetActiveSprints(ctx context.Context, boardID int) ([]jira.Sprint, error) {
	// TODO: Implement get active sprints
	return nil, fmt.Errorf("not implemented")
}

// GetSprintIssues implements Service.GetSprintIssues.
func (j *jiraService) GetSprintIssues(ctx context.Context, sprintID int) ([]jira.Issue, error) {
	// TODO: Implement get sprint issues
	return nil, fmt.Errorf("not implemented")
}

// CreateIssue implements Service.CreateIssue.
func (j *jiraService) CreateIssue(ctx context.Context, input *jira.CreateIssueInput) (*jira.Issue, error) {
	// TODO: Implement create issue
	return nil, fmt.Errorf("not implemented")
}

// UpdateIssue implements Service.UpdateIssue.
func (j *jiraService) UpdateIssue(ctx context.Context, input *jira.UpdateIssueInput) error {
	// TODO: Implement update issue
	return fmt.Errorf("not implemented")
}

// AddComment implements Service.AddComment.
func (j *jiraService) AddComment(ctx context.Context, issueKey, body string) error {
	// TODO: Implement add comment
	return fmt.Errorf("not implemented")
}

// DeleteIssue implements Service.DeleteIssue.
func (j *jiraService) DeleteIssue(ctx context.Context, issueKey string) error {
	// TODO: Implement delete issue
	return fmt.Errorf("not implemented")
}

// CreateSprint implements Service.CreateSprint.
func (j *jiraService) CreateSprint(ctx context.Context, boardID int, name, goal string) (*jira.Sprint, error) {
	// TODO: Implement create sprint
	return nil, fmt.Errorf("not implemented")
}

// CloseSprint implements Service.CloseSprint.
func (j *jiraService) CloseSprint(ctx context.Context, sprintID int) error {
	// TODO: Implement close sprint
	return fmt.Errorf("not implemented")
}
