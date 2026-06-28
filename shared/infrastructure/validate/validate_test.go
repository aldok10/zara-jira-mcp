// Package validate_test provides tests for the validate package.
// All tests follow the table-driven test pattern as recommended by the Go project.
package validate_test

import (
	"testing"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/validate"
)

func TestRequired(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want bool // true = error expected (invalid)
	}{
		{name: "empty string", arg: "", want: true},
		{name: "whitespace only", arg: "   ", want: true},
		{name: "valid string", arg: "hello", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Required("field", tt.arg)
			if tt.want && err == nil {
				t.Errorf("Required() expected error for arg %q", tt.arg)
			}
			if !tt.want && err != nil {
				t.Errorf("Required() unexpected error: %v", err)
			}
		})
	}
}

func TestIssueKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool // true = error expected (invalid)
	}{
		{name: "valid simple", key: "PROJ-123", want: false},
		{name: "valid with underscore", key: "PROJ_ABC-42", want: false},
		{name: "valid single digit", key: "ENG-1", want: false},
		{name: "empty key", key: "", want: true},
		{name: "lowercase prefix", key: "proj-123", want: true},
		{name: "no hyphen", key: "PROJ123", want: true},
		{name: "no digits", key: "PROJ-ABC", want: true},
		{name: "prefix with digits", key: "PROJ1-23", want: false},
		{name: "invalid characters", key: "PROJ-1@3", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.IssueKey(tt.key)
			if tt.want && err == nil {
				t.Errorf("IssueKey(%q) expected error", tt.key)
			}
			if !tt.want && err != nil {
				t.Errorf("IssueKey(%q) unexpected error: %v", tt.key, err)
			}
		})
	}
}

func TestJQL(t *testing.T) {
	tests := []struct {
		name string
		jql  string
		want bool // true = error expected (invalid)
	}{
		{name: "valid query", jql: "project = PROJ AND status = 'Done'", want: false},
		{name: "valid complex query", jql: "assignee = currentUser() AND sprint in openSprints() ORDER BY priority DESC", want: false},
		{name: "empty query", jql: "", want: true},
		{name: "delete keyword", jql: "DELETE FROM issues", want: true},
		{name: "drop keyword", jql: "drop table issues", want: true},
		{name: "truncate keyword", jql: "TRUNCATE issues", want: true},
		{name: "alter keyword", jql: "ALTER TABLE issues", want: true},
		{name: "sql injection attempt", jql: "'; DROP TABLE issues; --", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.JQL(tt.jql)
			if tt.want && err == nil {
				t.Errorf("JQL(%q) expected error", tt.jql)
			}
			if !tt.want && err != nil {
				t.Errorf("JQL(%q) unexpected error: %v", tt.jql, err)
			}
		})
	}
}

func TestProjectKey(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{name: "valid 2 char", key: "PR", want: false},
		{name: "valid 10 char", key: "ABCDEFGHIJ", want: false},
		{name: "with digits", key: "PROJ1", want: false},
		{name: "empty", key: "", want: true},
		{name: "too short", key: "P", want: true},
		{name: "too long", key: "ABCDEFGHIJK", want: true},
		{name: "with special chars", key: "PROJ-123", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.ProjectKey(tt.key)
			if tt.want && err == nil {
				t.Errorf("ProjectKey(%q) expected error", tt.key)
			}
			if !tt.want && err != nil {
				t.Errorf("ProjectKey(%q) unexpected error: %v", tt.key, err)
			}
		})
	}
}

func TestBoardID(t *testing.T) {
	tests := []struct {
		name string
		id   int
		want bool
	}{
		{name: "valid positive", id: 1, want: false},
		{name: "valid large", id: 1000, want: false},
		{name: "zero", id: 0, want: true},
		{name: "negative", id: -1, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.BoardID(tt.id)
			if tt.want && err == nil {
				t.Errorf("BoardID(%d) expected error", tt.id)
			}
			if !tt.want && err != nil {
				t.Errorf("BoardID(%d) unexpected error: %v", tt.id, err)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{name: "valid simple", email: "user@example.com", want: false},
		{name: "valid with dots", email: "first.last@company.co.uk", want: false},
		{name: "valid with plus", email: "user+tag@example.com", want: false},
		{name: "empty", email: "", want: true},
		{name: "no at sign", email: "userexample.com", want: true},
		{name: "no domain", email: "user@", want: true},
		{name: "no TLD", email: "user@example", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Email(tt.email)
			if tt.want && err == nil {
				t.Errorf("Email(%q) expected error", tt.email)
			}
			if !tt.want && err != nil {
				t.Errorf("Email(%q) unexpected error: %v", tt.email, err)
			}
		})
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{name: "valid https", url: "https://example.com", want: false},
		{name: "valid with path", url: "https://example.com/api/v1", want: false},
		{name: "valid http", url: "http://example.com", want: false},
		{name: "empty", url: "", want: true},
		{name: "no scheme", url: "example.com", want: true},
		{name: "too short", url: "http://a", want: true},
		{name: "ftp scheme", url: "ftp://example.com", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.URL(tt.url)
			if tt.want && err == nil {
				t.Errorf("URL(%q) expected error", tt.url)
			}
			if !tt.want && err != nil {
				t.Errorf("URL(%q) unexpected error: %v", tt.url, err)
			}
		})
	}
}

func TestPriority(t *testing.T) {
	tests := []struct {
		name     string
		priority string
		want     bool
	}{
		{name: "highest", priority: "Highest", want: false},
		{name: "high", priority: "High", want: false},
		{name: "medium", priority: "Medium", want: false},
		{name: "low", priority: "Low", want: false},
		{name: "lowest", priority: "Lowest", want: false},
		{name: "empty", priority: "", want: false}, // empty is allowed
		{name: "invalid", priority: "Critical", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Priority(tt.priority)
			if tt.want && err == nil {
				t.Errorf("Priority(%q) expected error", tt.priority)
			}
			if !tt.want && err != nil {
				t.Errorf("Priority(%q) unexpected error: %v", tt.priority, err)
			}
		})
	}
}

func TestLabels(t *testing.T) {
	tests := []struct {
		name   string
		labels string
		want   int  // expected number of labels for success, -1 for error
	}{
		{name: "empty", labels: "", want: 0},
		{name: "single", labels: "frontend", want: 1},
		{name: "multiple", labels: "frontend,backend,api", want: 3},
		{name: "with spaces", labels: "frontend, backend", want: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validate.Labels(tt.labels)
			if tt.want < 0 && err == nil {
				t.Errorf("Labels(%q) expected error", tt.labels)
			}
			if tt.want >= 0 && err != nil {
				t.Errorf("Labels(%q) unexpected error: %v", tt.labels, err)
			}
			if tt.want >= 0 && len(result) != tt.want {
				t.Errorf("Labels(%q) got %d labels, want %d", tt.labels, len(result), tt.want)
			}
		})
	}
}

func TestAccountID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{name: "valid 10 chars", id: "1234567890", want: false},
		{name: "valid 20 chars", id: "12345678901234567890", want: false},
		{name: "empty", id: "", want: true},
		{name: "too short", id: "123456789", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.AccountID(tt.id)
			if tt.want && err == nil {
				t.Errorf("AccountID(%q) expected error", tt.id)
			}
			if !tt.want && err != nil {
				t.Errorf("AccountID(%q) unexpected error: %v", tt.id, err)
			}
		})
	}
}
