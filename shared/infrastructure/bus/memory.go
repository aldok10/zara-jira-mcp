package bus

import (
	"context"
	"log/slog"
	"sync"

	"github.com/aldok10/zara-jira-mcp/shared/kernel/event"
)

// InMemoryBus is an async in-memory event bus with per-subscriber goroutines.
type InMemoryBus struct {
	mu       sync.RWMutex
	handlers map[string][]event.Handler
}

// NewInMemoryBus creates a new InMemoryBus.
func NewInMemoryBus() *InMemoryBus {
	return &InMemoryBus{
		handlers: make(map[string][]event.Handler),
	}
}

func (b *InMemoryBus) Publish(ctx context.Context, e event.Event) error {
	b.mu.RLock()
	handlers := b.handlers[e.EventName()]
	b.mu.RUnlock()

	if len(handlers) == 0 {
		return nil
	}

	// Fire-and-forget async using background context.
	// Caller's ctx may be cancelled before handlers finish.
	bg := context.Background()
	for _, h := range handlers {
		h := h
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("event handler panicked",
						"event", e.EventName(),
						"recover", r,
					)
				}
			}()
			if err := h.HandleEvent(bg, e); err != nil {
				slog.Error("event handler failed", "event", e.EventName(), "error", err)
			}
		}()
	}
	return nil
}

func (b *InMemoryBus) Subscribe(eventName string, handler event.Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}

func (b *InMemoryBus) SubscribeFunc(eventName string, fn func(ctx context.Context, e event.Event) error) {
	b.Subscribe(eventName, event.HandlerFunc(fn))
}
