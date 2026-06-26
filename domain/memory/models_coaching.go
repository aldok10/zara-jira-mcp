package memory

import "time"

// TeamPulse records a team health pulse survey entry.
type TeamPulse struct {
	ID         int64
	SprintName string
	Member     string
	Score      int
	Notes      string
	CreatedAt  time.Time
}

// MeetingEffectiveness records how effective a ceremony was.
type MeetingEffectiveness struct {
	ID              int64
	Ceremony        string
	DurationMinutes int
	Score           int
	Notes           string
	SprintName      string
	CreatedAt       time.Time
}

// TeamRadar records a single dimension of a team health radar.
type TeamRadar struct {
	ID         int64
	SprintName string
	Dimension  string
	Score      int
	CreatedAt  time.Time
}
