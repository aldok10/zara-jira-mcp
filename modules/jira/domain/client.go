package domain

import "context"

// Client defines the interface for Jira API operations.
// Implementations are in the infrastructure layer.
type Client interface {
	SearchIssues(ctx context.Context, jql string, maxResults int, startAt int) (*SearchResult, error)
	GetIssue(ctx context.Context, key string) (*Issue, error)
	GetBoards(ctx context.Context) ([]Board, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]Issue, error)
	CreateIssue(ctx context.Context, input *CreateIssueInput) (*Issue, error)
	UpdateIssue(ctx context.Context, input *UpdateIssueInput) error
	AddComment(ctx context.Context, issueKey, body string) error
	TransitionIssue(ctx context.Context, issueKey, transitionID string) error
	GetTransitions(ctx context.Context, issueKey string) ([]Transition, error)
	AssignIssue(ctx context.Context, issueKey, accountID string) error
	DeleteIssue(ctx context.Context, issueKey string) error
	CreateSubtask(ctx context.Context, parentKey string, input *CreateIssueInput) (*Issue, error)
	FindUser(ctx context.Context, query string) ([]User, error)
	SetEpicLink(ctx context.Context, issueKey, epicKey string) error
	RemoveEpicLink(ctx context.Context, issueKey string) error
	GetSprints(ctx context.Context, boardID int, state string) ([]Sprint, error)
	CreateSprint(ctx context.Context, boardID int, name, goal string) (*Sprint, error)
	StartSprint(ctx context.Context, sprintID int, startDate, endDate string) error
	CloseSprint(ctx context.Context, sprintID int) error
	MoveIssuesToSprint(ctx context.Context, sprintID int, issueKeys []string) error
	LinkIssues(ctx context.Context, inwardKey, outwardKey, linkType string) error
	GetLinkTypes(ctx context.Context) ([]LinkType, error)
	AddWorklog(ctx context.Context, issueKey, timeSpent, comment string) error
	GetWorklogs(ctx context.Context, issueKey string) ([]Worklog, error)
	AddWatcher(ctx context.Context, issueKey, accountID string) error
	GetWatchers(ctx context.Context, issueKey string) ([]User, error)
	AddLabel(ctx context.Context, issueKey, label string) error
	GetProjects(ctx context.Context) ([]Project, error)
	GetProject(ctx context.Context, key string) (*ProjectDetail, error)
	RawRequest(ctx context.Context, method, path string, body []byte) ([]byte, int, error)
	// Attachments
	GetAttachments(ctx context.Context, issueKey string) ([]Attachment, error)
	// Versions
	GetVersions(ctx context.Context, projectKey string) ([]Version, error)
	CreateVersion(ctx context.Context, projectKey, name, description string) (*Version, error)
	ReleaseVersion(ctx context.Context, versionID string) error
	// Components
	GetComponents(ctx context.Context, projectKey string) ([]Component, error)
	// Fields
	GetFields(ctx context.Context) ([]Field, error)
	// Board configuration
	GetBoardConfiguration(ctx context.Context, boardID int) (*BoardConfiguration, error)
}
