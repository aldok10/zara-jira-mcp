package domain

// CreateIssueInput holds parameters for creating an issue.
type CreateIssueInput struct {
	Project     string
	Summary     string
	Description string
	IssueType   string // Task, Bug, Story
	Priority    string // optional
	Assignee    string // optional, account ID
	Labels      []string
}

// UpdateIssueInput holds parameters for updating an issue.
type UpdateIssueInput struct {
	Key         string
	Summary     string
	Description string
	Priority    string
	Assignee    string   // account ID, empty = no change
	Labels      []string // nil = no change, empty slice = clear
}
