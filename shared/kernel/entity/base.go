// Package entity provides base types for domain entities following DDD tactical patterns.
package entity

import "time"

// Entity is the base interface all domain entities must implement.
type Entity interface {
	// Identity returns the unique identifier for this entity.
	Identity() string
}

// AggregateRoot marks an entity as the root of an aggregate.
// All external references to an aggregate must go through the root.
type AggregateRoot interface {
	Entity
}

// BaseEntity provides common fields for entities that need audit trails.
type BaseEntity struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Identity implements Entity.Identity.
func (b BaseEntity) Identity() string {
	return b.ID
}
