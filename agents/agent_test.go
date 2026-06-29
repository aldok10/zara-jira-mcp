package agents

import (
	"context"
	"testing"

	"github.com/aldok10/zara-jira-mcp/shared/domain/event"
)

// mockAgent implements Agent for testing.
type mockAgent struct {
	name       string
	desc       string
	eventTypes []string
	execFn     func(ctx context.Context, req *Request) (*Result, error)
}

func (m *mockAgent) Name() string         { return m.name }
func (m *mockAgent) Description() string  { return m.desc }
func (m *mockAgent) EventTypes() []string { return m.eventTypes }
func (m *mockAgent) Execute(ctx context.Context, req *Request) (*Result, error) {
	if m.execFn != nil {
		return m.execFn(ctx, req)
	}
	return &Result{Success: true}, nil
}

func TestDispatcher_RegisterAndDispatch(t *testing.T) {
	d := NewDispatcher()

	executed := false
	agent := &mockAgent{
		name:       "test-agent",
		desc:       "test",
		eventTypes: []string{"test.event"},
		execFn: func(ctx context.Context, req *Request) (*Result, error) {
			executed = true
			if req.EventName != "test.event" {
				t.Errorf("expected test.event, got %s", req.EventName)
			}
			return &Result{Success: true}, nil
		},
	}

	d.Register(agent)

	// Dispatch to registered event
	result, err := d.Dispatch(context.Background(), &Request{EventName: "test.event"})
	if err != nil {
		t.Fatalf("dispatch error: %v", err)
	}
	if !executed {
		t.Fatal("agent was not executed")
	}
	if !result.Success {
		t.Fatal("expected success")
	}

	// Dispatch to unregistered event — should return nil, nil
	result, err = d.Dispatch(context.Background(), &Request{EventName: "unknown.event"})
	if err != nil {
		t.Fatalf("dispatch error for unknown event: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result for unregistered event")
	}
}

func TestDispatcher_RegisteredAgents(t *testing.T) {
	d := NewDispatcher()

	d.Register(&mockAgent{name: "agent-a", eventTypes: []string{"event.a"}})
	d.Register(&mockAgent{name: "agent-b", eventTypes: []string{"event.b", "event.c"}})

	agents := d.RegisteredAgents()
	if len(agents) != 2 {
		t.Fatalf("expected 2 agents, got %d: %v", len(agents), agents)
	}
}

func TestDispatcher_RegisteredEventTypes(t *testing.T) {
	d := NewDispatcher()

	d.Register(&mockAgent{name: "agent-a", eventTypes: []string{"event.a", "event.b"}})
	d.Register(&mockAgent{name: "agent-b", eventTypes: []string{"event.b", "event.c"}})

	types := d.RegisteredEventTypes()
	if len(types) != 3 {
		t.Fatalf("expected 3 unique event types, got %d: %v", len(types), types)
	}
}

func TestBusBridge_HandleEvent(t *testing.T) {
	d := NewDispatcher()

	executed := false
	d.Register(&mockAgent{
		name:       "test-agent",
		desc:       "test",
		eventTypes: []string{"system.sprint.completed"},
		execFn: func(ctx context.Context, req *Request) (*Result, error) {
			executed = true
			if req.EventName != "system.sprint.completed" {
				t.Errorf("unexpected event: %s", req.EventName)
			}
			if req.Payload["board_id"] != 42 {
				t.Errorf("expected board_id 42, got %v (type: %T)", req.Payload["board_id"], req.Payload["board_id"])
			}
			return &Result{Success: true}, nil
		},
	})

	bridge := NewBusBridge(d)

	// Simulate event handler call
	err := bridge.HandleEvent(context.Background(), &event.SprintCompleted{
		BoardID:    42,
		SprintID:   100,
		SprintName: "Sprint 10",
		BoardName:  "Dev Board",
	})
	if err != nil {
		t.Fatalf("handle event error: %v", err)
	}
	if !executed {
		t.Fatal("agent was not executed via bridge")
	}
}

func TestBusBridge_UnregisteredEvent(t *testing.T) {
	d := NewDispatcher()
	bridge := NewBusBridge(d)

	// No agents registered — should not error
	err := bridge.HandleEvent(context.Background(), &event.HealthScoreComputed{
		BoardID:    1,
		SprintName: "Sprint 1",
		Score:      85,
		Rating:     "Healthy",
	})
	if err != nil {
		t.Fatalf("handle unregistered event should not error: %v", err)
	}
}

func TestBusBridge_SubscribeTo(t *testing.T) {
	d := NewDispatcher()
	d.Register(&mockAgent{
		name:       "test-agent",
		eventTypes: []string{"domain.health_score.computed", "system.sprint.completed"},
	})

	bus := &mockBus{subscriptions: make(map[string]int)}
	bridge := NewBusBridge(d)
	bridge.SubscribeTo(bus)

	if bus.subscriptions["domain.health_score.computed"] != 1 {
		t.Errorf("expected 1 subscription for domain.health_score.computed, got %d", bus.subscriptions["domain.health_score.computed"])
	}
	if bus.subscriptions["system.sprint.completed"] != 1 {
		t.Errorf("expected 1 subscription for system.sprint.completed, got %d", bus.subscriptions["system.sprint.completed"])
	}
	if len(bus.subscriptions) != 2 {
		t.Errorf("expected 2 subscriptions total, got %d", len(bus.subscriptions))
	}
}

// mockBus implements event.Bus for testing.
type mockBus struct {
	subscriptions map[string]int
}

func (m *mockBus) Publish(ctx context.Context, e event.Event) error { return nil }
func (m *mockBus) Subscribe(eventName string, handler event.Handler) {
	m.subscriptions[eventName]++
}
func (m *mockBus) SubscribeFunc(eventName string, fn func(ctx context.Context, e event.Event) error) {
	m.subscriptions[eventName]++
}
