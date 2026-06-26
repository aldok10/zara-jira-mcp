package memory

import "context"

// Store defines the persistence interface for PM memory.
type Store interface {
	// Sprint snapshots
	SaveSprintSnapshot(ctx context.Context, s *SprintSnapshot) error
	GetSprintSnapshots(ctx context.Context, boardID int, limit int) ([]SprintSnapshot, error)
	GetLatestSnapshot(ctx context.Context, boardID int) (*SprintSnapshot, error)

	// Risks
	SaveRisk(ctx context.Context, r *Risk) error
	UpdateRisk(ctx context.Context, r *Risk) error
	GetOpenRisks(ctx context.Context) ([]Risk, error)
	GetAllRisks(ctx context.Context, limit int) ([]Risk, error)

	// Decisions
	SaveDecision(ctx context.Context, d *Decision) error
	GetDecisions(ctx context.Context, limit int) ([]Decision, error)
	SearchDecisions(ctx context.Context, query string) ([]Decision, error)

	// Blockers
	SaveBlocker(ctx context.Context, b *Blocker) error
	ResolveBlocker(ctx context.Context, id int64, resolution string) error
	GetActiveBlockers(ctx context.Context) ([]Blocker, error)
	GetBlockerHistory(ctx context.Context, limit int) ([]Blocker, error)

	// Team metrics
	SaveTeamMetric(ctx context.Context, m *TeamMetric) error
	GetTeamMetrics(ctx context.Context, memberName string, limit int) ([]TeamMetric, error)
	GetTeamOverview(ctx context.Context, sprintName string) ([]TeamMetric, error)

	// Retrospectives
	SaveRetrospective(ctx context.Context, r *Retrospective) error
	GetRetrospectives(ctx context.Context, limit int) ([]Retrospective, error)
	SaveActionItem(ctx context.Context, a *ActionItem) error
	GetPendingActionItems(ctx context.Context) ([]ActionItem, error)
	CompleteActionItem(ctx context.Context, id int64) error

	// Dependencies
	SaveDependency(ctx context.Context, d *Dependency) error
	ResolveDependency(ctx context.Context, id int64) error
	GetOpenDependencies(ctx context.Context) ([]Dependency, error)
	GetDependenciesForIssue(ctx context.Context, issueKey string) ([]Dependency, error)

	// Meeting notes
	SaveMeetingNote(ctx context.Context, m *MeetingNote) error
	GetMeetingNotes(ctx context.Context, meetingType string, limit int) ([]MeetingNote, error)

	// Health scores
	SaveHealthScore(ctx context.Context, h *HealthScore) error
	GetHealthScores(ctx context.Context, boardID int, limit int) ([]HealthScore, error)
}
