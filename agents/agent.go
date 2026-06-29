// Package agents provides the AI Agent Architecture layer.
//
// Agents are NOT part of the domain. They live in a separate layer and
// only communicate via System Events. Business domains never call agents
// directly.
//
// Architecture layers:
//
//	Agent Registry → Dispatcher → Planner → Coordinator → Executor → Tool Adapter
package agents

import "context"

// ---------------------------------------------------------------------------
// Core types
// ---------------------------------------------------------------------------

// Agent defines the interface for all AI agents.
// Every agent implements this interface and is selected by the Dispatcher
// based on the incoming System Event.
type Agent interface {
	// Name returns the agent's unique identifier.
	Name() string

	// Description returns a human-readable description.
	Description() string

	// EventTypes returns the system event names this agent handles.
	EventTypes() []string

	// Execute runs the agent with the given context and event data.
	Execute(ctx context.Context, req *Request) (*Result, error)
}

// Request carries the system event data for agent execution.
type Request struct {
	EventName     string
	Payload       map[string]interface{}
	WorkflowID    string
	CorrelationID string
	AgentID       string
}

// Result carries the agent's output after execution.
type Result struct {
	Success   bool
	Data      interface{}
	Error     string
	EventName string // Optional: new event to publish
	Actions   []Action
}

// Action is a single executable step produced by the planner/coordinator.
type Action struct {
	Tool   string
	Input  map[string]interface{}
	Output interface{}
	Error  string
}

// ---------------------------------------------------------------------------
// Dispatcher
// ---------------------------------------------------------------------------

// Dispatcher routes system events to the appropriate agent.
type Dispatcher struct {
	agents map[string]Agent
}

// NewDispatcher creates a new agent dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		agents: make(map[string]Agent),
	}
}

// Register adds an agent to the dispatcher.
func (d *Dispatcher) Register(a Agent) {
	for _, eventType := range a.EventTypes() {
		d.agents[eventType] = a
	}
}

// Dispatch routes a system event to the registered agent.
func (d *Dispatcher) Dispatch(ctx context.Context, req *Request) (*Result, error) {
	agent, ok := d.agents[req.EventName]
	if !ok {
		return nil, nil
	}
	return agent.Execute(ctx, req)
}

// ---------------------------------------------------------------------------
// Planner — decomposes a goal into a plan
// ---------------------------------------------------------------------------

// Planner decomposes a goal into a sequence of actions.
type Planner interface {
	Plan(ctx context.Context, goal string, context map[string]interface{}) ([]Action, error)
}

// ---------------------------------------------------------------------------
// Coordinator — orchestrates plan execution
// ---------------------------------------------------------------------------

// Coordinator orchestrates execution of a plan across tools.
type Coordinator interface {
	Execute(ctx context.Context, actions []Action, toolRegistry map[string]interface{}) (*Result, error)
}

// ---------------------------------------------------------------------------
// Executor — runs individual tool calls
// ---------------------------------------------------------------------------

// Executor runs a single tool call.
type Executor interface {
	Execute(ctx context.Context, tool string, input map[string]interface{}) (interface{}, error)
}
