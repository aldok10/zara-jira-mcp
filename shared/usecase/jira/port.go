package jira

import (
	"context"

	"github.com/aldok10/zara-jira-mcp/shared/domain/jira"
)

// Client provides all Jira API operations.
type Client interface {
	SearchIssues(ctx context.Context, jql string, maxResults int, startAt int) (*jira.SearchResult, error)
	GetIssue(ctx context.Context, key string) (*jira.Issue, error)
	GetBoards(ctx context.Context) ([]jira.Board, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]jira.Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]jira.Issue, error)
	CreateIssue(ctx context.Context, input *jira.CreateIssueInput) (*jira.Issue, error)
	UpdateIssue(ctx context.Context, input *jira.UpdateIssueInput) error
	AddComment(ctx context.Context, issueKey, body string) error
	TransitionIssue(ctx context.Context, issueKey, transitionID string) error
	AssignIssue(ctx context.Context, issueKey, accountID string) error
	DeleteIssue(ctx context.Context, issueKey string) error
	CreateSubtask(ctx context.Context, parentKey string, input *jira.CreateIssueInput) (*jira.Issue, error)
	GetSprints(ctx context.Context, boardID int, state string) ([]jira.Sprint, error)
	CreateSprint(ctx context.Context, boardID int, name, goal string) (*jira.Sprint, error)
	StartSprint(ctx context.Context, sprintID int, startDate, endDate string) error
	CloseSprint(ctx context.Context, sprintID int) error
	MoveIssuesToSprint(ctx context.Context, sprintID int, issueKeys []string) error
	LinkIssues(ctx context.Context, inwardKey, outwardKey, linkType string) error
	GetLinkTypes(ctx context.Context) ([]jira.LinkType, error)
	AddWorklog(ctx context.Context, issueKey, timeSpent, comment string) error
	GetWorklogs(ctx context.Context, issueKey string) ([]jira.Worklog, error)
	AddWatcher(ctx context.Context, issueKey, accountID string) error
	GetWatchers(ctx context.Context, issueKey string) ([]jira.User, error)
	AddLabel(ctx context.Context, issueKey, label string) error
	GetProjects(ctx context.Context) ([]jira.Project, error)
	GetProject(ctx context.Context, key string) (*jira.ProjectDetail, error)
	RawRequest(ctx context.Context, method, path string, body []byte) ([]byte, int, error)
	// Attachments
	GetAttachments(ctx context.Context, issueKey string) ([]jira.Attachment, error)
	// Versions
	GetVersions(ctx context.Context, projectKey string) ([]jira.Version, error)
	CreateVersion(ctx context.Context, projectKey, name, description string) (*jira.Version, error)
	ReleaseVersion(ctx context.Context, versionID string) error
	// Components
	GetComponents(ctx context.Context, projectKey string) ([]jira.Component, error)
	// Fields
	GetFields(ctx context.Context) ([]jira.Field, error)
}

// Cache provides caching for Jira data.
type Cache interface {
	Get(key string) ([]byte, bool)
	Set(key string, data []byte)
	TTL() int
}
