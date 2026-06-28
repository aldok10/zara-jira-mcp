package domain

import "fmt"

// ErrJiraAPI indicates a Jira API error with status code.
type ErrJiraAPI struct {
	StatusCode int
	Message    string
	Endpoint   string
}

func (e *ErrJiraAPI) Error() string {
	return fmt.Sprintf("jira api error %d on %s: %s", e.StatusCode, e.Endpoint, e.Message)
}
