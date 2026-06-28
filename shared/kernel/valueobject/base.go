// Package valueobject provides base types for immutable value objects.
package valueobject

// ValueObject is the base interface all value objects must implement.
// Value objects are immutable, have no identity, and are compared by their properties.
type ValueObject interface {
	// Equals checks structural equality with another value object.
	Equals(other ValueObject) bool
}
