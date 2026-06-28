// Package event provides domain event infrastructure for decoupled cross-cutting concerns.
package event

import "context"

// Event is a domain event that something meaningful happened.
type Event interface {
	// EventName returns the event type name for routing.
	EventName() string
}

// Handler processes a published event.
type Handler interface {
	HandleEvent(ctx context.Context, event Event) error
}

// HandlerFunc adapts a function to the Handler interface.
type HandlerFunc func(ctx context.Context, event Event) error

func (f HandlerFunc) HandleEvent(ctx context.Context, event Event) error {
	return f(ctx, event)
}

// Bus publishes events to registered handlers asynchronously.
type Bus interface {
	// Publish emits an event to all subscribed handlers.
	Publish(ctx context.Context, event Event) error
	// Subscribe registers a handler for an event type.
	Subscribe(eventName string, handler Handler)
	// SubscribeFunc registers a function handler for an event type.
	SubscribeFunc(eventName string, fn func(ctx context.Context, event Event) error)
}
