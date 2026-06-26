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

	// Daily progress (burndown)
	SaveDailyProgress(ctx context.Context, p *DailyProgress) error
	GetDailyProgress(ctx context.Context, boardID int, sprintName string) ([]DailyProgress, error)

	// Sprint goals
	SaveSprintGoal(ctx context.Context, g *SprintGoal) error
	UpdateSprintGoal(ctx context.Context, g *SprintGoal) error
	GetActiveGoals(ctx context.Context, boardID int) ([]SprintGoal, error)
	GetGoalHistory(ctx context.Context, boardID int, limit int) ([]SprintGoal, error)

	// Definition of Done
	SaveDoDItem(ctx context.Context, item *DoDItem) error
	GetDoD(ctx context.Context, project string) ([]DoDItem, error)
	DeleteDoDItem(ctx context.Context, id int64) error

	// Escalations
	SaveEscalation(ctx context.Context, e *Escalation) error
	GetRecentEscalations(ctx context.Context, limit int) ([]Escalation, error)
	AcknowledgeEscalation(ctx context.Context, id int64) error

	// Team Pulse
	SaveTeamPulse(ctx context.Context, p *TeamPulse) error
	GetTeamPulseHistory(ctx context.Context, limit int) ([]TeamPulse, error)

	// Meeting Effectiveness
	SaveMeetingEffectiveness(ctx context.Context, m *MeetingEffectiveness) error
	GetMeetingEffectivenessHistory(ctx context.Context, ceremony string, limit int) ([]MeetingEffectiveness, error)

	// Team Radar
	SaveTeamRadar(ctx context.Context, r *TeamRadar) error
	GetTeamRadarHistory(ctx context.Context, limit int) ([]TeamRadar, error)

	// Raw database access for ad-hoc queries
	DB() RawDB
}

// RawDB provides raw database access for custom queries.
type RawDB interface {
	Exec(query string, args ...any) (any, error)
	Query(query string, args ...any) (Rows, error)
	QueryRow(query string, args ...any) Row
}

// Rows is a minimal interface for sql.Rows.
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
}

// Row is a minimal interface for sql.Row.
type Row interface {
	Scan(dest ...any) error
}
