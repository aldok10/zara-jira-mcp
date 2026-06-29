package event

import "time"

// ---------------------------------------------------------------------------
// Event name prefixes for the three layers:
//   domain.      — business event (aggregate changed)
//   application. — internal workflow event
//   system.      — automation / AI agent / infrastructure event
// ---------------------------------------------------------------------------

// Domain Events

// HealthScoreComputed is published when a sprint health score is computed.
type HealthScoreComputed struct {
	BaseEvent
	BoardID    int
	SprintName string
	Score      int    // 0-100
	Rating     string // Healthy / Fair / At Risk / Critical
}

func (e HealthScoreComputed) EventName() string { return "domain.health_score.computed" }

// AntiPatternDetected is published when anti-patterns are found.
type AntiPatternDetected struct {
	BaseEvent
	BoardID      int
	PatternCount int
	PatternNames []string
}

func (e AntiPatternDetected) EventName() string { return "domain.antipattern.detected" }

// ForecastGenerated is published when a forecast is completed.
type ForecastGenerated struct {
	BaseEvent
	BoardID        int
	RemainingItems int
	Confidence50   float64
	Confidence85   float64
}

func (e ForecastGenerated) EventName() string { return "domain.forecast.generated" }

// BlockerEscalated is published when a blocker exceeds threshold.
type BlockerEscalated struct {
	BaseEvent
	BoardID   int
	BlockerID int64
	IssueKey  string
	DaysOld   int
	Severity  string
}

func (e BlockerEscalated) EventName() string { return "domain.blocker.escalated" }

// ---------------------------------------------------------------------------
// Application Events — internal workflow orchestration
// ---------------------------------------------------------------------------

// SendWelcomeEmail is dispatched when a new user needs onboarding.
type SendWelcomeEmail struct {
	BaseEvent
	UserID   string
	Email    string
	UserName string
}

func (e SendWelcomeEmail) EventName() string { return "application.send_welcome_email" }

// RefreshCache signals that cached data should be invalidated/refreshed.
type RefreshCache struct {
	BaseEvent
	CacheKey string
}

func (e RefreshCache) EventName() string { return "application.refresh_cache" }

// GenerateInvoice triggers invoice generation for an order.
type GenerateInvoice struct {
	BaseEvent
	OrderID string
	UserID  string
}

func (e GenerateInvoice) EventName() string { return "application.generate_invoice" }

// ---------------------------------------------------------------------------
// System Events — automation, CI/CD, AI agents, infrastructure
// ---------------------------------------------------------------------------

// PipelineFailed is published when a CI/CD pipeline fails.
type PipelineFailed struct {
	BaseEvent
	PipelineID string
	Project    string
	Branch     string
	FailedAt   string
}

func (e PipelineFailed) EventName() string { return "system.pipeline.failed" }

// MergeRequestOpened is published when a merge/PR is created.
type MergeRequestOpened struct {
	BaseEvent
	MRID      int
	Project   string
	Author    string
	SourceRef string
	TargetRef string
	URL       string
}

func (e MergeRequestOpened) EventName() string { return "system.merge_request.opened" }

// RiskDetected is published when the AI detects a project risk.
type RiskDetected struct {
	BaseEvent
	Title       string
	Severity    string // critical, high, medium, low
	Description string
	Owner       string
	Source      string
}

func (e RiskDetected) EventName() string { return "system.risk.detected" }

// SprintCompleted is published when a sprint ends.
type SprintCompleted struct {
	BaseEvent
	BoardID    int
	SprintID   int
	SprintName string
	BoardName  string
}

func (e SprintCompleted) EventName() string { return "system.sprint.completed" }

// SprintCreated is published when a new sprint starts.
type SprintCreated struct {
	BaseEvent
	BoardID    int
	SprintID   int
	SprintName string
	Goal       string
}

func (e SprintCreated) EventName() string { return "system.sprint.created" }

// SprintUpdated is published when the sprint backlog changes.
type SprintUpdated struct {
	BaseEvent
	BoardID  int
	SprintID int
}

func (e SprintUpdated) EventName() string { return "system.sprint.updated" }

// SprintBacklogChanged is published when issues move in/out of the sprint.
type SprintBacklogChanged struct {
	BaseEvent
	BoardID    int
	SprintID   int
	ChangeType string // added, removed, reprioritized
	IssueKeys  []string
}

func (e SprintBacklogChanged) EventName() string { return "system.sprint.backlog_changed" }

// BlockerDetected is published when a blocker is identified.
type BlockerDetected struct {
	BaseEvent
	IssueKey    string
	Description string
	Owner       string
	Source      string
}

func (e BlockerDetected) EventName() string { return "system.blocker.detected" }

// DecisionRecorded is published when an architectural/technical decision is made.
type DecisionRecorded struct {
	BaseEvent
	Title     string
	Decision  string
	Context   string
	Rationale string
	MadeBy    string
	Tags      []string
}

func (e DecisionRecorded) EventName() string { return "system.decision.recorded" }

// WeeklyRiskScanRequested is published by the scheduler for periodic risk scanning.
type WeeklyRiskScanRequested struct {
	BaseEvent
	BoardID int
}

func (e WeeklyRiskScanRequested) EventName() string { return "system.weekly_risk_scan.requested" }

// StandupRequested is published daily for standup preparation.
type StandupRequested struct {
	BaseEvent
	BoardID int
}

func (e StandupRequested) EventName() string { return "system.standup.requested" }

// CodeAnalyzed is published when code analysis completes.
type CodeAnalyzed struct {
	BaseEvent
	Project      string
	Branch       string
	IssuesFound  int
	IssueDetails []string
}

func (e CodeAnalyzed) EventName() string { return "system.code.analyzed" }

// ---------------------------------------------------------------------------
// Metadata helpers
// ---------------------------------------------------------------------------

// NewMetadata creates a Metadata with OccurredAt set to now.
func NewMetadata() Metadata {
	return Metadata{
		OccurredAt: time.Now(),
	}
}

// AttachMetadata sets metadata on an event and returns it.
func AttachMetadata(e Event, m Metadata) Event {
	e.SetMetadata(m)
	return e
}
