package memory

import "time"

// SprintSnapshot captures the state of a sprint at a point in time.
type SprintSnapshot struct {
	ID             int64
	SprintName     string
	BoardID        int
	SnapshotDate   time.Time
	TotalIssues    int
	Done           int
	InProgress     int
	Todo           int
	Blocked        int
	Carryover      int // issues carried from previous sprint
	Velocity       int // story points completed
	CompletionRate float64
	Notes          string
}

// Risk tracks identified project risks.
type Risk struct {
	ID          int64
	Title       string
	Description string
	Severity    string // critical, high, medium, low
	Status      string // open, mitigating, resolved, accepted
	Owner       string
	Mitigation  string
	IdentifiedAt time.Time
	ResolvedAt  *time.Time
	SprintName  string
}

// Decision records a project decision with context.
type Decision struct {
	ID        int64
	Title     string
	Context   string // what situation led to this decision
	Decision  string // what was decided
	Rationale string // why this over alternatives
	Outcome   string // result after implementation (filled later)
	MadeBy    string
	MadeAt    time.Time
	Tags      string // comma-separated
}

// Blocker tracks impediments and their resolution.
type Blocker struct {
	ID           int64
	IssueKey     string // Jira issue key if linked
	Description  string
	BlockedSince time.Time
	ResolvedAt   *time.Time
	Resolution   string
	Owner        string
	DaysBlocked  int
}

// TeamMetric captures individual team member workload signals.
type TeamMetric struct {
	ID           int64
	MemberName   string
	SprintName   string
	RecordedAt   time.Time
	IssuesAssigned int
	IssuesDone   int
	BlockerCount int
	CarryoverCount int
	Notes        string
}

// Retrospective stores retro outcomes and action items.
type Retrospective struct {
	ID          int64
	SprintName  string
	Date        time.Time
	WentWell    string
	Improvements string
	ActionItems string // JSON array or newline-separated
	Status      string // open, closed
}

// ActionItem from a retrospective.
type ActionItem struct {
	ID          int64
	RetroID     int64
	Description string
	Owner       string
	DueDate     *time.Time
	Status      string // pending, done, cancelled
}
