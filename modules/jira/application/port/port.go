// Package port defines the inbound and outbound boundaries for the jira module.
package port

import (
	"context"

	"github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// Inbound defines the Jira use cases exposed by this module.
type Inbound interface {
	SearchIssues(ctx context.Context, jql string, maxResults int) (*domain.SearchResult, error)
	GetIssue(ctx context.Context, key string) (*domain.Issue, error)
	GetBoards(ctx context.Context) ([]domain.Board, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]domain.Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]domain.Issue, error)
	CreateIssue(ctx context.Context, input *domain.CreateIssueInput) (*domain.Issue, error)
	UpdateIssue(ctx context.Context, input *domain.UpdateIssueInput) error
	AddComment(ctx context.Context, issueKey, body string) error
	DeleteIssue(ctx context.Context, issueKey string) error
	CreateSprint(ctx context.Context, boardID int, name, goal string) (*domain.Sprint, error)
	CloseSprint(ctx context.Context, sprintID int) error
}

// Outbound defines the ports this module requires from external systems.
type Outbound interface {
	Client
	Cache
}

// Client provides all Jira API operations.
type Client interface {
	SearchIssues(ctx context.Context, jql string, maxResults int, startAt int) (*domain.SearchResult, error)
	GetIssue(ctx context.Context, key string) (*domain.Issue, error)
	GetBoards(ctx context.Context) ([]domain.Board, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]domain.Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]domain.Issue, error)
	CreateIssue(ctx context.Context, input *domain.CreateIssueInput) (*domain.Issue, error)
	UpdateIssue(ctx context.Context, input *domain.UpdateIssueInput) error
	AddComment(ctx context.Context, issueKey, body string) error
	TransitionIssue(ctx context.Context, issueKey, transitionID string) error
	AssignIssue(ctx context.Context, issueKey, accountID string) error
	DeleteIssue(ctx context.Context, issueKey string) error
	CreateSubtask(ctx context.Context, parentKey string, input *domain.CreateIssueInput) (*domain.Issue, error)
	GetSprints(ctx context.Context, boardID int, state string) ([]domain.Sprint, error)
	CreateSprint(ctx context.Context, boardID int, name, goal string) (*domain.Sprint, error)
	StartSprint(ctx context.Context, sprintID int, startDate, endDate string) error
	CloseSprint(ctx context.Context, sprintID int) error
	MoveIssuesToSprint(ctx context.Context, sprintID int, issueKeys []string) error
	LinkIssues(ctx context.Context, inwardKey, outwardKey, linkType string) error
	GetLinkTypes(ctx context.Context) ([]domain.LinkType, error)
	AddWorklog(ctx context.Context, issueKey, timeSpent, comment string) error
	GetWorklogs(ctx context.Context, issueKey string) ([]domain.Worklog, error)
	AddWatcher(ctx context.Context, issueKey, accountID string) error
	GetWatchers(ctx context.Context, issueKey string) ([]domain.User, error)
	AddLabel(ctx context.Context, issueKey, label string) error
	GetProjects(ctx context.Context) ([]domain.Project, error)
	GetProject(ctx context.Context, key string) (*domain.ProjectDetail, error)
	RawRequest(ctx context.Context, method, path string, body []byte) ([]byte, int, error)
	GetAttachments(ctx context.Context, issueKey string) ([]domain.Attachment, error)
	GetVersions(ctx context.Context, projectKey string) ([]domain.Version, error)
	CreateVersion(ctx context.Context, projectKey, name, description string) (*domain.Version, error)
	ReleaseVersion(ctx context.Context, versionID string) error
	GetComponents(ctx context.Context, projectKey string) ([]domain.Component, error)
	GetFields(ctx context.Context) ([]domain.Field, error)
}

// Cache provides caching for Jira data.
type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, data []byte)
	TTL() int
}
