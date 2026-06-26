package jira

import "context"

// Client defines the interface for Jira API operations.
type Client interface {
	SearchIssues(ctx context.Context, jql string, maxResults int) (*SearchResult, error)
	GetIssue(ctx context.Context, key string) (*Issue, error)
	GetBoards(ctx context.Context) ([]Board, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]Issue, error)
}
