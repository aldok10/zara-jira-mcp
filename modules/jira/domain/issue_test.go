package domain

import (
	"testing"
)

func TestIssue(t *testing.T) {
	// Test that issue type is correctly initialized
	i := Issue{}
	if i.Key != "" {
		t.Errorf("new Issue should have empty Key, got %q", i.Key)
	}
	if i.StoryPoints != 0 {
		t.Errorf("new Issue should have zero StoryPoints, got %f", i.StoryPoints)
	}
}

func TestCreateIssueInput(t *testing.T) {
	input := CreateIssueInput{
		Project:     "PROJ",
		Summary:     "Test issue",
		IssueType:   "Task",
		Description: "Test description",
		Priority:    "High",
		Labels:      []string{"bug", "critical"},
	}

	if input.Project != "PROJ" {
		t.Errorf("expected Project=PROJ, got %q", input.Project)
	}
	if input.Summary != "Test issue" {
		t.Errorf("expected Summary='Test issue', got %q", input.Summary)
	}
	if len(input.Labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(input.Labels))
	}
}
