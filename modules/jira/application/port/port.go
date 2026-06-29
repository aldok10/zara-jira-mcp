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
}
