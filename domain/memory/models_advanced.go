package memory

import "time"

// Dependency tracks cross-issue or cross-team dependencies.
type Dependency struct {
	ID            int64
	FromIssueKey  string // the issue that depends on something
	ToIssueKey    string // what it depends on (issue key or external)
	DependencyType string // blocks, blocked_by, relates_to, external
	Description   string
	Status        string // open, resolved
	CreatedAt     time.Time
	ResolvedAt    *time.Time
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
	ID           int64
	SprintName   string
	BoardID      int
	ComputedAt   time.Time
	OverallScore int    // 0-100
	VelocityScore int   // 0-25
	BlockerScore  int   // 0-25
	ScopeScore    int   // 0-25
	TeamScore     int   // 0-25
	Details       string // JSON breakdown
}
