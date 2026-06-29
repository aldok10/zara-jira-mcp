// Package event provides the domain event infrastructure for decoupled
// cross-cutting concerns across all three event layers:
//
//	Domain Events      — aggregate-level business events
//	Application Events — internal workflow orchestration
//	System Events      — automation, CI/CD, AI agents
//
// This is the CANONICAL event package.
package event

import (
	"context"
	"time"
)

// EventNameProvider provides the routing name for an event.
type EventNameProvider interface {
	EventName() string
}

// Metadata carries tracing and observability data on every event.
// All events MUST carry metadata for end-to-end traceability.
type Metadata struct {
	EventID       string    `json:"event_id"`
	CorrelationID string    `json:"correlation_id"`
	TraceID       string    `json:"trace_id"`
	WorkflowID    string    `json:"workflow_id"`
	RequestID     string    `json:"request_id"`
	SessionID     string    `json:"session_id"`
	AgentID       string    `json:"agent_id"`
	ToolID        string    `json:"tool_id"`
	WorkspaceID   string    `json:"workspace_id"`
	Module        string    `json:"module"`
	Version       string    `json:"version"`
	OccurredAt    time.Time `json:"occurred_at,omitempty"` // ISO8601 format
}

// Event is a domain event that something meaningful happened.
type Event interface {
	EventNameProvider
	Metadata() Metadata
	SetMetadata(m Metadata)
}

// BaseEvent provides the default Metadata implementation.
type BaseEvent struct {
	meta Metadata
}

func (b *BaseEvent) Metadata() Metadata     { return b.meta }
func (b *BaseEvent) SetMetadata(m Metadata) { b.meta = m }

// Handler processes a published event.
type Handler interface {
	HandleEvent(ctx context.Context, event Event) error
}

// HandlerFunc adapts a plain function to the Handler interface.
type HandlerFunc func(ctx context.Context, event Event) error

func (f HandlerFunc) HandleEvent(ctx context.Context, event Event) error {
	return f(ctx, event)
}

// Bus publishes events to registered handlers.
type Bus interface {
	// Publish emits an event to all subscribed handlers.
	Publish(ctx context.Context, event Event) error
	// Subscribe registers a handler for an event type.
	Subscribe(eventName string, handler Handler)
	// SubscribeFunc registers a function handler for an event type.
	SubscribeFunc(eventName string, fn func(ctx context.Context, event Event) error)
}
