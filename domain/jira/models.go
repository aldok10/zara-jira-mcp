package jira

import "time"

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
}

// SearchResult holds paginated search results.
type SearchResult struct {
	Issues     []Issue
	Total      int
	MaxResults int
}

// Board represents a Jira board.
type Board struct {
	ID   int
	Name string
	Type string
}

// Sprint represents a Jira sprint.
type Sprint struct {
	ID    int
	Name  string
	State string
	Goal  string
}
