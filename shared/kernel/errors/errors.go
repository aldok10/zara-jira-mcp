// Package errors provides foundational error types shared across all domains.
// Each domain may define additional domain-specific error types.
package errors

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
