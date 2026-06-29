// Package listener provides the Event Listener pattern for single-responsibility listeners.
//
// Each listener handles exactly ONE concern. For example:
//   - SendWelcomeEmailListener — sends welcome email
//   - AuditListener — logs audit trail
//   - PublishKafkaListener — publishes to Kafka
//   - RiskRecorderListener — records risk to PM memory
//
// Rule: One listener = one responsibility. Never create a listener that does
// multiple unrelated things.
package listener

import (
	"context"
	"fmt"
	"github.com/aldok10/zara-jira-mcp/shared/domain/event"
	"log/slog"
	"time"
)

// Listener defines the single-responsibility listener contract.
//
// Each listener should handle exactly one type of event and perform
// exactly one concern. If you need multiple things done in response
// to an event, register multiple listeners.
type Listener interface {
	// Handle processes the event. Returns error if the action failed.
	Handle(ctx context.Context, e event.Event) error

	// Name returns a human-readable name for observability.
	Name() string
}

// ListenerFunc adapts a function to the Listener interface.
type ListenerFunc struct {
	fn   func(ctx context.Context, e event.Event) error
	name string
}

func (l ListenerFunc) Handle(ctx context.Context, e event.Event) error {
	return l.fn(ctx, e)
}

func (l ListenerFunc) Name() string { return l.name }

// NewListenerFunc creates a Listener from a function.
func NewListenerFunc(name string, fn func(ctx context.Context, e event.Event) error) ListenerFunc {
	return ListenerFunc{fn: fn, name: name}
}

// Registry maps event names to their listeners.
type Registry struct {
	entries map[string][]Listener
}

// NewRegistry creates a new listener registry.
func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string][]Listener),
	}
}

// Register adds a listener for the given event name.
func (r *Registry) Register(eventName string, listener Listener) {
	r.entries[eventName] = append(r.entries[eventName], listener)
}

// Listeners returns all listeners registered for an event name.
func (r *Registry) Listeners(eventName string) []Listener {
	return r.entries[eventName]
}

// Adapter wraps a Listener into an event.Handler for use with the dispatcher.
type Adapter struct {
	listener Listener
}

// NewAdapter creates an event.Handler from a Listener.
func NewAdapter(l Listener) *Adapter {
	return &Adapter{listener: l}
}

func (a *Adapter) Name() string { return a.listener.Name() }

// HandleEvent implements event.Handler.
func (a *Adapter) HandleEvent(ctx context.Context, e event.Event) error {
	start := time.Now()
	err := a.listener.Handle(ctx, e)
	latency := time.Since(start)

	if err != nil {
		slog.Error("listener failed",
			"listener", a.listener.Name(),
			"event", e.EventName(),
			"latency", latency,
			"error", err,
		)
		return fmt.Errorf("%s failed on %s: %w", a.listener.Name(), e.EventName(), err)
	}

	slog.Debug("listener completed",
		"listener", a.listener.Name(),
		"event", e.EventName(),
		"latency", latency,
	)
	return nil
}

// RegisterAll registers all listeners for an event into the given bus.
// Convenience helper for bootstrap wiring.
func RegisterAll(bus event.Bus, registry *Registry) {
	for eventName, listeners := range registry.entries {
		for _, l := range listeners {
			adapter := NewAdapter(l)
			bus.Subscribe(eventName, adapter)
		}
	}
}
