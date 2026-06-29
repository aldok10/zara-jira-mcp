package peopleport

import (
	"context"

	jira "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
	memory "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
)

// TeamMetricRepository provides team metric data.
type TeamMetricRepository interface {
	FindByBoard(ctx context.Context, boardID int, limit int) ([]memory.TeamMetric, error)
	Save(ctx context.Context, m *memory.TeamMetric) error
}

// HealthRepository provides team health score data.
type HealthRepository interface {
	FindByBoard(ctx context.Context, boardID int, limit int) ([]memory.HealthScore, error)
	Save(ctx context.Context, h *memory.HealthScore) error
}

// BlockerRepository provides blocker data for people analysis.
type BlockerRepository interface {
	FindActive(ctx context.Context) ([]memory.Blocker, error)
	FindByBoard(ctx context.Context, boardID int) ([]memory.Blocker, error)
	Save(ctx context.Context, b *memory.Blocker) error
}

// JiraClient provides Jira data for team workload.
type JiraClient interface {
	GetBoards(ctx context.Context) ([]jira.Board, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]jira.Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]jira.Issue, error)
	SearchIssues(ctx context.Context, jql string, maxResults int, startAt int) (*jira.SearchResult, error)
	AddComment(ctx context.Context, issueKey, body string) error
}

// AIProvider provides AI analysis for people insights.
type AIProvider interface {
	Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}
