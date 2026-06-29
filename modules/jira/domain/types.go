package domain

// Transition represents an available workflow transition.
type Transition struct {
	ID   string
	Name string
}

// Worklog represents a time log entry on an issue.
type Worklog struct {
	Author    string
	TimeSpent string
	Started   string
	Comment   string
}

// Attachment represents a file attached to an issue.
type Attachment struct {
	ID       string
	Filename string
	Size     int64
	MimeType string
	Author   string
	Created  string
	URL      string
}

// Version represents a project release version.
type Version struct {
	ID          string
	Name        string
	Description string
	Released    bool
	ReleaseDate string
}

// Component represents a project component.
type Component struct {
	ID   string
	Name string
	Lead string
}

// Field represents a Jira field definition.
type Field struct {
	ID     string
	Name   string
	Custom bool
	Type   string
}

// LinkType represents a Jira issue link type.
type LinkType struct {
	Name    string
	Inward  string
	Outward string
}
