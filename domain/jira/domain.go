package jira

import (
	"context"
	"time"
)

// Issue represents a Jira issue with relevant fields for PM intelligence.
type Issue struct {
	Key         string
	Summary     string
	Description string
	Status      string
	Priority    string
	Type        string
	Assignee    string
	Reporter    string
	Labels      []string
	Created     time.Time
	Updated     time.Time
	SprintName  string
}

// SearchResult holds paginated search results.
type SearchResult struct {
	Issues     []Issue
	Total      int
	StartAt    int
	MaxResults int
	HasMore    bool
}

// Board represents a Jira board.
type Board struct {
	ID   int
	Name string
	Type string
}

// Sprint represents a Jira sprint.
type Sprint struct {
	ID        int
	Name      string
	State     string
	Goal      string
	StartDate string
	EndDate   string
}

// CreateIssueInput holds parameters for creating an issue.
type CreateIssueInput struct {
	Project     string
	Summary     string
	Description string
	IssueType   string // Task, Bug, Story
	Priority    string // optional
	Assignee    string // optional, account ID
	Labels      []string
}

// UpdateIssueInput holds parameters for updating an issue.
type UpdateIssueInput struct {
	Key         string
	Summary     string
	Description string
	Priority    string
	Assignee    string   // account ID, empty = no change
	Labels      []string // nil = no change, empty slice = clear
}

// Client defines the interface for Jira API operations.
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
}

// Transition represents an available workflow transition.
type Transition struct {
	ID   string
	Name string
}

// User represents a Jira user.
type User struct {
	AccountID   string
	DisplayName string
	Email       string
}

// Project represents a Jira project summary.
type Project struct {
	Key  string
	Name string
	Lead string
	Type string
}

// ProjectDetail represents full project info.
type ProjectDetail struct {
	Key         string
	Name        string
	Lead        string
	Type        string
	Description string
	Components  []string
	Versions    []string
}

// LinkType represents a Jira issue link type.
type LinkType struct {
	Name    string
	Inward  string
	Outward string
}

// Worklog represents a time log entry on an issue.
type Worklog struct {
	Author    string
	TimeSpent string
	Started   string
	Comment   string
}

// Attachment represents a file attached to an issue.
type Attachment struct {
	ID       string
	Filename string
	Size     int64
	MimeType string
	Author   string
	Created  string
	URL      string
}

// Version represents a project release version.
type Version struct {
	ID          string
	Name        string
	Description string
	Released    bool
	ReleaseDate string
}

// Component represents a project component.
type Component struct {
	ID   string
	Name string
	Lead string
}

// Field represents a Jira field definition.
type Field struct {
	ID     string
	Name   string
	Custom bool
	Type   string
}
