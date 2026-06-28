package analysis

import (
	"context"

	"github.com/aldok10/zara-jira-mcp/shared/domain/jira"
	memory "github.com/aldok10/zara-jira-mcp/shared/domain/memory"
)

// SnapshotRepository provides sprint snapshot data for analysis.
type SnapshotRepository interface {
	FindByBoard(ctx context.Context, boardID int, limit int) ([]memory.SprintSnapshot, error)
}

// MeetingRepository provides meeting effectiveness data.
type MeetingRepository interface {
	FindByCeremony(ctx context.Context, ceremony string, limit int) ([]memory.MeetingEffectiveness, error)
	Save(ctx context.Context, m *memory.MeetingEffectiveness) error
}

// ActionItemRepository provides action item data.
type ActionItemRepository interface {
	FindPending(ctx context.Context) ([]memory.ActionItem, error)
}

// RetroRepository provides retrospective data.
type RetroRepository interface {
	FindRecent(ctx context.Context, limit int) ([]memory.Retrospective, error)
}

// JiraClient provides Jira data for analysis.
type JiraClient interface {
	GetBoards(ctx context.Context) ([]jira.Board, error)
	GetActiveSprints(ctx context.Context, boardID int) ([]jira.Sprint, error)
	GetSprintIssues(ctx context.Context, sprintID int) ([]jira.Issue, error)
	SearchIssues(ctx context.Context, jql string, maxResults int, startAt int) (*jira.SearchResult, error)
	GetSprints(ctx context.Context, boardID int, state string) ([]jira.Sprint, error)
}

// AIProvider provides AI analysis.
type AIProvider interface {
	Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}
