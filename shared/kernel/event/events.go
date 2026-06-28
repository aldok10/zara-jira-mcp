package event

// --- Concrete Event Types ---
// These are shared cross-module event contracts.
// Module-specific events live in the respective module's domain layer.

// HealthScoreComputed is published when a sprint health score is computed.
type HealthScoreComputed struct {
	BoardID    int
	SprintName string
	Score      int    // Overall score 0-100
	Rating     string // Healthy / Fair / At Risk / Critical
}

func (e HealthScoreComputed) EventName() string { return "health_score.computed" }

// AntiPatternDetected is published when anti-patterns are found.
type AntiPatternDetected struct {
	BoardID      int
	PatternCount int
	PatternNames []string
}

func (e AntiPatternDetected) EventName() string { return "antipattern.detected" }

// ForecastGenerated is published when a forecast is completed.
type ForecastGenerated struct {
	BoardID        int
	RemainingItems int
	Confidence50   float64
	Confidence85   float64
}

func (e ForecastGenerated) EventName() string { return "forecast.generated" }

// BlockerEscalated is published when a blocker exceeds threshold.
type BlockerEscalated struct {
	BoardID   int
	BlockerID int64
	IssueKey  string
	DaysOld   int
	Severity  string
}

func (e BlockerEscalated) EventName() string { return "blocker.escalated" }
