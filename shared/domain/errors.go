package domain

import "fmt"

// ErrNotFound indicates a requested resource doesn't exist.
type ErrNotFound struct {
	Resource string
	ID       string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Resource, e.ID)
}

// ErrValidation indicates invalid input.
type ErrValidation struct {
	Field   string
	Message string
}

func (e *ErrValidation) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

// ErrJiraAPI indicates a Jira API error with status code.
type ErrJiraAPI struct {
	StatusCode int
	Message    string
	Endpoint   string
}

func (e *ErrJiraAPI) Error() string {
	return fmt.Sprintf("jira api error %d on %s: %s", e.StatusCode, e.Endpoint, e.Message)
}
