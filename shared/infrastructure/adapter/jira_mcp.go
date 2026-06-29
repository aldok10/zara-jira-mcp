package adapter

import (
	"context"
	"fmt"
	"github.com/aldok10/zara-jira-mcp/shared/domain/event"
	"log/slog"
)

// JiraMCPAdapter implements the tool adapter pattern for Jira MCP.
type JiraMCPAdapter struct {
	eventBus event.Bus
}

func NewJiraMCPAdapter(eventBus event.Bus) *JiraMCPAdapter {
	return &JiraMCPAdapter{eventBus: eventBus}
}

func (a *JiraMCPAdapter) Name() string { return "jira_mcp" }
func (a *JiraMCPAdapter) Description() string {
	return "Jira MCP integration with event-driven architecture"
}

func (a *JiraMCPAdapter) Execute(ctx context.Context, action string, input map[string]interface{}) (interface{}, error) {
	// Based on action, publish appropriate System Event
	switch action {
	case "get_sprint_status":
		return a.publishSystemEvent(ctx, "system.sprint.status.requested", input)
	case "get_blockers":
		return a.publishSystemEvent(ctx, "system.blocker.requested", input)
	case "get_risks":
		return a.publishSystemEvent(ctx, "system.risk.requested", input)
	case "record_blocker":
		return a.publishSystemEvent(ctx, "system.blocker.detected", input)
	case "record_decision":
		return a.publishSystemEvent(ctx, "system.decision.recorded", input)
	default:
		return nil, fmt.Errorf("unknown Jira MCP action: %s", action)
	}
}

func (a *JiraMCPAdapter) publishSystemEvent(ctx context.Context, eventName string, input map[string]interface{}) (interface{}, error) {
	// Create System Event
	eventMap := convertMap(input)
	// In real implementation, we would create actual System event structs from events.go
	// For now, use a generic event
	event := NewGenericEvent(eventName, eventMap)

	// Publish to event bus
	if err := a.eventBus.Publish(ctx, event); err != nil {
		slog.Error("failed to publish Jira MCP event", "event", eventName, "error", err)
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// GenericEvent is a simple implementation of event.Event for system events.
type GenericEvent struct {
	meta    event.Metadata
	name    string
	payload map[string]interface{}
}

// NewGenericEvent creates a new GenericEvent.
func NewGenericEvent(name string, data map[string]interface{}) *GenericEvent {
	return &GenericEvent{name: name, payload: data}
}

func (g *GenericEvent) EventName() string            { return g.name }
func (g *GenericEvent) Metadata() event.Metadata     { return g.meta }
func (g *GenericEvent) SetMetadata(m event.Metadata) { g.meta = m }

func convertMap(input map[string]interface{}) map[string]interface{} {
	temp := make(map[string]interface{})
	for k, v := range input {
		temp[k] = v
	}
	return temp
}
