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
	MaxResults int
}

// Board represents a Jira board.
type Board struct {
	ID   int
	Name string
	Type string
}

// Sprint represents a Jira sprint.
type Sprint struct {
	ID    int
	Name  string
	State string
	Goal  string
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

// Client defines the interface for Jira API operations.
type Client interface {
	SearchIssues(ctx context.Context, jql string, maxResults int) (*SearchResult, error)
	GetIssue(ctx context.Context, key string) (*Issue, error)
	GetBoards(ctx context.Context) ([]Board, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]Issue, error)
	CreateIssue(ctx context.Context, input *CreateIssueInput) (*Issue, error)
	AddComment(ctx context.Context, issueKey, body string) error
	TransitionIssue(ctx context.Context, issueKey, transitionID string) error
	GetTransitions(ctx context.Context, issueKey string) ([]Transition, error)
}

// Transition represents an available workflow transition.
type Transition struct {
	ID   string
	Name string
}
