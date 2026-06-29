// Package port defines the inbound boundaries for the jira module.
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
	GetDefaultBoardID(ctx context.Context) (int, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]domain.Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]domain.Issue, error)
	CreateIssue(ctx context.Context, input *domain.CreateIssueInput) (*domain.Issue, error)
	UpdateIssue(ctx context.Context, input *domain.UpdateIssueInput) error
	AddComment(ctx context.Context, issueKey, body string) error
	DeleteIssue(ctx context.Context, issueKey string) error
	CreateSprint(ctx context.Context, boardID int, name, goal string) (*domain.Sprint, error)
	CloseSprint(ctx context.Context, sprintID int) error

	// Transitions
	TransitionIssue(ctx context.Context, issueKey, transitionID string) error
	GetTransitions(ctx context.Context, issueKey string) ([]domain.Transition, error)

	// Assignments
	AssignIssue(ctx context.Context, issueKey, accountID string) error
	FindUser(ctx context.Context, query string) ([]domain.User, error)

	// Sprint management (beyond active)
	GetSprints(ctx context.Context, boardID int, state string) ([]domain.Sprint, error)
	StartSprint(ctx context.Context, sprintID int, startDate, endDate string) error
	MoveIssuesToSprint(ctx context.Context, sprintID int, issueKeys []string) error

	// Subtasks
	CreateSubtask(ctx context.Context, parentKey string, input *domain.CreateIssueInput) (*domain.Issue, error)

	// Links
	LinkIssues(ctx context.Context, inwardKey, outwardKey, linkType string) error
	GetLinkTypes(ctx context.Context) ([]domain.LinkType, error)

	// Worklogs
	AddWorklog(ctx context.Context, issueKey, timeSpent, comment string) error
	GetWorklogs(ctx context.Context, issueKey string) ([]domain.Worklog, error)

	// Watchers
	AddWatcher(ctx context.Context, issueKey, accountID string) error
	GetWatchers(ctx context.Context, issueKey string) ([]domain.User, error)

	// Labels
	AddLabel(ctx context.Context, issueKey, label string) error

	// Projects
	GetProjects(ctx context.Context) ([]domain.Project, error)
	GetProject(ctx context.Context, key string) (*domain.ProjectDetail, error)

	// Versions
	GetVersions(ctx context.Context, projectKey string) ([]domain.Version, error)
	CreateVersion(ctx context.Context, projectKey, name, description string) (*domain.Version, error)
	ReleaseVersion(ctx context.Context, versionID string) error

	// Attachments
	GetAttachments(ctx context.Context, issueKey string) ([]domain.Attachment, error)

	// Components
	GetComponents(ctx context.Context, projectKey string) ([]domain.Component, error)

	// Fields
	GetFields(ctx context.Context) ([]domain.Field, error)
}
