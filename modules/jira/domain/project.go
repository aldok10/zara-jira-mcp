package domain

// Project represents a Jira project summary.
type Project struct {
	Key  string
	Name string
	Lead string
	Type string
}

// ProjectDetail represents full project info.
type ProjectDetail struct {
	Key         string
	Name        string
	Lead        string
	Type        string
	Description string
	Components  []string
	Versions    []string
}
