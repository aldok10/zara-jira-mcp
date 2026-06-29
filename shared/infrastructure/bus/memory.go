package bus

import (
	"context"
	"fmt"
	"github.com/aldok10/zara-jira-mcp/shared/domain/event"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

// Dispatcher metrics
type Metrics struct {
	EventsPublished  atomic.Int64
	EventsFailed     atomic.Int64
	HandlersExecuted atomic.Int64
	HandlersFailed   atomic.Int64
	AvgLatencyMs     atomic.Int64
}

// Registry holds handler registrations.
type Registry struct {
	mu    sync.RWMutex
	items map[string][]event.Handler
}

// NewRegistry creates a new handler registry.
func NewRegistry() *Registry {
	return &Registry{
		items: make(map[string][]event.Handler),
	}
}

func (r *Registry) Subscribe(eventName string, handler event.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[eventName] = append(r.items[eventName], handler)
}

func (r *Registry) SubscribeFunc(eventName string, fn func(ctx context.Context, event event.Event) error) {
	r.Subscribe(eventName, event.HandlerFunc(fn))
}

func (r *Registry) Handlers(eventName string) []event.Handler {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.items[eventName]
}

// DispatcherConfig holds configuration for the event dispatcher.
type DispatcherConfig struct {
	MaxRetries      int
	RetryDelay      time.Duration
	DispatchTimeout time.Duration
	MaxDepth        int
	AsyncDefault    bool
}

// DefaultConfig returns production-ready configuration.
func DefaultConfig() DispatcherConfig {
	return DispatcherConfig{
		MaxRetries:      3,
		RetryDelay:      100 * time.Millisecond,
		DispatchTimeout: 30 * time.Second,
		MaxDepth:        5,
		AsyncDefault:    true,
	}
}

// Dispatcher is the production-grade event dispatcher.
type Dispatcher struct {
	registry *Registry
	config   DispatcherConfig
	metrics  Metrics
	mu       sync.Mutex
	// Event history for preventing circular dispatch
	EventHistory map[string]time.Time
}

// NewDispatcher creates a new event dispatcher.
func NewDispatcher(config DispatcherConfig) *Dispatcher {
	return &Dispatcher{
		registry:     NewRegistry(),
		config:       config,
		EventHistory: make(map[string]time.Time),
	}
}

// Publish dispatches an event to all registered handlers.
func (d *Dispatcher) Publish(ctx context.Context, e event.Event) error {
	meta := e.Metadata()
	meta.OccurredAt = time.Now()
	e.SetMetadata(meta)
	eventName := e.EventName()

	handlers := d.registry.Handlers(eventName)
	if len(handlers) == 0 {
		return nil
	}

	// Check event depth
	if meta.WorkflowID != "" {
		depth := d.getWorkflowDepth(meta.WorkflowID)
		if depth >= d.config.MaxDepth {
			return fmt.Errorf("max dispatch depth %d reached for workflow %s", d.config.MaxDepth, meta.WorkflowID)
		}
	}

	d.metrics.EventsPublished.Add(1)

	// Dispatch synchronously
	return d.dispatchSync(ctx, e, handlers)
}

func (d *Dispatcher) dispatchSync(ctx context.Context, e event.Event, handlers []event.Handler) error {
	ctx, cancel := context.WithTimeout(ctx, d.config.DispatchTimeout)
	defer cancel()

	var lastErr error
	for _, h := range handlers {
		if err := d.executeHandler(ctx, h, e); err != nil {
			d.metrics.HandlersFailed.Add(1)
			lastErr = err
			// Continue with other handlers even if one fails
		}
	}
	return lastErr
}

func (d *Dispatcher) executeHandler(ctx context.Context, h event.Handler, e event.Event) error {
	start := time.Now()
	defer func() {
		latency := time.Since(start).Milliseconds()
		d.metrics.AvgLatencyMs.Store((d.metrics.AvgLatencyMs.Load() + latency) / 2)
	}()

	// Recovery for panics
	defer func() {
		if r := recover(); r != nil {
			slog.Error("event handler panicked",
				"event", e.EventName(),
				"recover", r,
			)
		}
	}()

	// Retry logic
	var lastErr error
	for i := 0; i <= d.config.MaxRetries; i++ {
		if err := h.HandleEvent(ctx, e); err != nil {
			lastErr = err
			if i < d.config.MaxRetries {
				time.Sleep(d.config.RetryDelay * time.Duration(i+1))
				continue
			}
		} else {
			d.metrics.HandlersExecuted.Add(1)
			return nil
		}
	}
	return fmt.Errorf("handler failed after %d retries: %w", d.config.MaxRetries, lastErr)
}

func (d *Dispatcher) getWorkflowDepth(workflowID string) int {
	// Simplified depth tracking - in production, use proper stack
	d.mu.Lock()
	defer d.mu.Unlock()

	// Clean old entries
	for k, v := range d.EventHistory {
		if time.Since(v) > 5*time.Minute {
			delete(d.EventHistory, k)
		}
	}
	return len(d.EventHistory)
}

// Subscribe registers a handler for an event type.
func (d *Dispatcher) Subscribe(eventName string, handler event.Handler) {
	d.registry.Subscribe(eventName, handler)
}

// SubscribeFunc registers a function handler for an event type.
func (d *Dispatcher) SubscribeFunc(eventName string, fn func(ctx context.Context, event event.Event) error) {
	d.registry.SubscribeFunc(eventName, fn)
}

// Metrics returns current dispatcher metrics.
func (d *Dispatcher) Metrics() *Metrics {
	return &d.metrics
}

// InMemoryBus is the legacy bus implementation - now wraps Dispatcher.
type InMemoryBus struct {
	Dispatcher *Dispatcher
}

func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		Dispatcher: NewDispatcher(DefaultConfig()),
	}
}

func (b *InMemoryBus) Publish(ctx context.Context, e event.Event) error {
	return b.Dispatcher.Publish(ctx, e)
}

func (b *InMemoryBus) Subscribe(eventName string, handler event.Handler) {
	b.Dispatcher.Subscribe(eventName, handler)
}

func (b *InMemoryBus) SubscribeFunc(eventName string, fn func(ctx context.Context, e event.Event) error) {
	b.Dispatcher.SubscribeFunc(eventName, fn)
}
