// Package agents provides the AI Agent Architecture layer.
package agents

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aldok10/zara-jira-mcp/shared/domain/event"
)

// BusBridge connects the Event Bus to the Agent Dispatcher.
// It subscribes to bus events and routes them to registered agents.
type BusBridge struct {
	dispatcher *Dispatcher
}

// NewBusBridge creates a bridge that subscribes to bus events
// and forwards them to the agent dispatcher.
func NewBusBridge(dispatcher *Dispatcher) *BusBridge {
	return &BusBridge{dispatcher: dispatcher}
}

// HandleEvent implements event.Handler. It converts a system event
// to an agent Request and dispatches it to the registered agent.
func (b *BusBridge) HandleEvent(ctx context.Context, e event.Event) error {
	eventName := e.EventName()
	meta := e.Metadata()

	req := &Request{
		EventName:     eventName,
		CorrelationID: meta.CorrelationID,
		WorkflowID:    meta.WorkflowID,
		AgentID:       meta.AgentID,
		Payload:       map[string]interface{}{},
	}

	// Convert event fields to payload map
	payload := make(map[string]interface{})
	switch v := e.(type) {
	case *event.HealthScoreComputed:
		payload["board_id"] = v.BoardID
		payload["sprint_name"] = v.SprintName
		payload["score"] = v.Score
		payload["rating"] = v.Rating
	case *event.AntiPatternDetected:
		payload["board_id"] = v.BoardID
		payload["pattern_count"] = v.PatternCount
		payload["pattern_names"] = v.PatternNames
	case *event.BlockerEscalated:
		payload["board_id"] = v.BoardID
		payload["blocker_id"] = v.BlockerID
		payload["issue_key"] = v.IssueKey
		payload["days_old"] = v.DaysOld
		payload["severity"] = v.Severity
	case *event.SprintCompleted:
		payload["board_id"] = v.BoardID
		payload["sprint_id"] = v.SprintID
		payload["sprint_name"] = v.SprintName
		payload["board_name"] = v.BoardName
	case *event.SprintCreated:
		payload["board_id"] = v.BoardID
		payload["sprint_id"] = v.SprintID
		payload["sprint_name"] = v.SprintName
		payload["goal"] = v.Goal
	case *event.SprintUpdated:
		payload["board_id"] = v.BoardID
		payload["sprint_id"] = v.SprintID
	case *event.RiskDetected:
		payload["title"] = v.Title
		payload["severity"] = v.Severity
		payload["description"] = v.Description
		payload["owner"] = v.Owner
		payload["source"] = v.Source
	case *event.BlockerDetected:
		payload["issue_key"] = v.IssueKey
		payload["description"] = v.Description
		payload["owner"] = v.Owner
		payload["source"] = v.Source
	case *event.DecisionRecorded:
		payload["title"] = v.Title
		payload["decision"] = v.Decision
		payload["context"] = v.Context
		payload["rationale"] = v.Rationale
		payload["made_by"] = v.MadeBy
		payload["tags"] = v.Tags
	}
	req.Payload = payload

	result, err := b.dispatcher.Dispatch(ctx, req)
	if err != nil {
		slog.Warn("agent dispatch failed", "event", eventName, "error", err)
		return fmt.Errorf("dispatch event %s: %w", eventName, err)
	}

	// If agent produced a follow-up event, publish it
	if result != nil && result.EventName != "" {
		slog.Debug("agent produced follow-up event", "event", result.EventName)
	}

	return nil
}

// SubscribeTo subscribes the bridge to all events the dispatcher's agents handle.
// This connects every agent's EventTypes() to the bus.
func (b *BusBridge) SubscribeTo(bus event.Bus) {
	for _, et := range b.dispatcher.RegisteredEventTypes() {
		bus.Subscribe(et, b)
		slog.Debug("agent bridge subscribed", "event", et)
	}
}
