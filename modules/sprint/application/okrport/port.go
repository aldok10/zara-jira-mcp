package okrport

import (
	"context"

	jira "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
	memory "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
)

// Repository provides OKR signal persistence.
type Repository interface {
	FindByBoard(ctx context.Context, boardID int) ([]memory.OKRSignal, error)
	Save(ctx context.Context, o *memory.OKRSignal) error
	UpdateProgress(ctx context.Context, id int64, progress float64) error
}

// SprintRepository provides sprint data for OKR alignment.
type SprintRepository interface {
	FindByBoard(ctx context.Context, boardID int, limit int) ([]memory.SprintSnapshot, error)
	FindLatest(ctx context.Context, boardID int) (*memory.SprintSnapshot, error)
}

// JiraClient provides Jira data for OKR progress calculation.
type JiraClient interface {
	SearchIssues(ctx context.Context, jql string, maxResults int, startAt int) (*jira.SearchResult, error)
	GetSprints(ctx context.Context, boardID int, state string) ([]jira.Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]jira.Issue, error)
}

// AIProvider provides AI analysis for OKR suggestions.
type AIProvider interface {
	Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}
