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
	DueDate     *time.Time // standard Jira field
	SprintName  string
	StoryPoints float64                // from customfield (0 = unestimated)
	Custom      map[string]interface{} // all custom field values for discovery
}

// dueKeywords are status substrings indicating completion.
var dueKeywords = []string{"done", "closed", "resolved", "complete", "finish", "merged"}

// blockedKeywords are status substrings indicating blocked.
var blockedKeywords = []string{"blocked", "block", "impediment", "waiting"}

// progressKeywords are status substrings indicating active work.
var progressKeywords = []string{"progress", "review", "testing", "dev", "development", "implement", "working", "in-progress", "in progress"}

// IsDone returns true if the issue status indicates completion.
func (i *Issue) IsDone() bool {
	lower := strings.ToLower(i.Status)
	for _, kw := range dueKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// IsBlocked returns true if the issue status indicates blocked.
func (i *Issue) IsBlocked() bool {
	lower := strings.ToLower(i.Status)
	for _, kw := range blockedKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// IsInProgress returns true if the issue is actively being worked.
func (i *Issue) IsInProgress() bool {
	if i.IsBlocked() {
		return false
	}
	lower := strings.ToLower(i.Status)
	for _, kw := range progressKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
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

// CustomString returns a custom field value as string, if present.
func (i *Issue) CustomString(fieldID string) (string, bool) {
	if i.Custom == nil {
		return "", false
	}
	v, ok := i.Custom[fieldID]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// CustomNumber returns a custom field value as float64, if present.
func (i *Issue) CustomNumber(fieldID string) (float64, bool) {
	if i.Custom == nil {
		return 0, false
	}
	v, ok := i.Custom[fieldID]
	if !ok {
		return 0, false
	}
	n, ok := v.(float64)
	return n, ok
}

// CustomTime returns a custom field value as time.Time, if present.
func (i *Issue) CustomTime(fieldID string) (time.Time, bool) {
	if i.Custom == nil {
		return time.Time{}, false
	}
	v, ok := i.Custom[fieldID]
	if !ok {
		return time.Time{}, false
	}
	t, ok := v.(time.Time)
	return t, ok
}

// HasCustomField returns true if the issue has a value for the given custom field ID.
func (i *Issue) HasCustomField(fieldID string) bool {
	if i.Custom == nil {
		return false
	}
	_, ok := i.Custom[fieldID]
	return ok
}
