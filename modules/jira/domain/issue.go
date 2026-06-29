// Package domain provides Jira domain entities and interfaces.
package domain

import (
	"strings"
	"time"
)

// Issue represents a Jira issue with relevant fields for PM intelligence.
type Issue struct {
	Key         string
	Summary     string
	Description string
	Status      string
	Priority    string
	Type        string
	Assignee    string
	Reporter    string
	Labels      []string
	Created     time.Time
	Updated     time.Time
	SprintName  string
	StoryPoints float64 // from customfield (0 = unestimated)
}

// IsDone returns true if the issue status indicates completion.
func (i *Issue) IsDone() bool {
	lower := strings.ToLower(i.Status)
	switch {
	case strings.Contains(lower, "done"),
		strings.Contains(lower, "closed"),
		strings.Contains(lower, "resolved"),
		strings.Contains(lower, "complete"):
		return true
	default:
		return false
	}
}

// IsBlocked returns true if the issue status indicates blocked.
func (i *Issue) IsBlocked() bool {
	return strings.Contains(strings.ToLower(i.Status), "blocked")
}

// IsInProgress returns true if the issue is actively being worked.
func (i *Issue) IsInProgress() bool {
	if i.IsBlocked() {
		return false
	}
	lower := strings.ToLower(i.Status)
	switch {
	case strings.Contains(lower, "progress"),
		strings.Contains(lower, "review"),
		strings.Contains(lower, "testing"),
		strings.Contains(lower, "dev"):
		return true
	default:
		return false
	}
}

// Classification returns the general status class: "done", "progress", "blocked", "todo".
func (i *Issue) Classification() string {
	switch {
	case i.IsDone():
		return "done"
	case i.IsBlocked():
		return "blocked"
	case i.IsInProgress():
		return "progress"
	default:
		return "todo"
	}
}

// IsEstimated returns true if the issue has story points.
func (i *Issue) IsEstimated() bool {
	return i.StoryPoints > 0
}
