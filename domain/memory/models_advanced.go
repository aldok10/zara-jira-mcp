package memory

import "time"

// Dependency tracks cross-issue or cross-team dependencies.
type Dependency struct {
	ID             int64
	FromIssueKey   string // the issue that depends on something
	ToIssueKey     string // what it depends on (issue key or external)
	DependencyType string // blocks, blocked_by, relates_to, external
	Description    string
	Status         string // open, resolved
	CreatedAt      time.Time
	ResolvedAt     *time.Time
}

// MeetingNote records outcomes from ceremonies (standup, planning, etc).
type MeetingNote struct {
	ID          int64
	MeetingType string // standup, planning, retro, grooming, adhoc
	Date        time.Time
	Attendees   string // comma-separated
	Notes       string
	Decisions   string // decisions made during meeting
	ActionItems string // follow-ups
	SprintName  string
}

// HealthScore captures a computed sprint health score.
type HealthScore struct {
	ID            int64
	SprintName    string
	BoardID       int
	ComputedAt    time.Time
	OverallScore  int    // 0-100
	VelocityScore int    // 0-25
	BlockerScore  int    // 0-25
	ScopeScore    int    // 0-25
	TeamScore     int    // 0-25
	Details       string // JSON breakdown
}

// DailyProgress tracks daily sprint progress for burndown charts.
type DailyProgress struct {
	ID          int64
	SprintName  string
	BoardID     int
	Date        time.Time
	TotalIssues int
	Done        int
	InProgress  int
	Todo        int
	Blocked     int
	PointsDone  int
	PointsTotal int
}

// SprintGoal tracks explicit sprint goals and their outcomes.
type SprintGoal struct {
	ID          int64
	SprintName  string
	BoardID     int
	Goal        string
	KeyResults  string // newline-separated measurable outcomes
	Status      string // active, achieved, partially_achieved, missed
	Outcome     string // what actually happened
	CreatedAt   time.Time
	ClosedAt    *time.Time
}

// DoDItem represents a Definition of Done checklist item.
type DoDItem struct {
	ID        int64
	Project   string // project key (or "*" for global)
	Item      string // checklist item description
	Category  string // code, testing, docs, review, deploy
	OrderNum  int
	Active    bool
}

// Escalation tracks auto-escalated items sent to Lark.
type Escalation struct {
	ID          int64
	Type        string // risk, blocker, stale
	ReferenceID int64  // ID of the risk/blocker that triggered
	Title       string
	Severity    string
	EscalatedAt time.Time
	Channel     string // lark, manual
	Acknowledged bool
}
