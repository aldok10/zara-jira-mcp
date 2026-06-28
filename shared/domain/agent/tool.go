package agent

import "context"

// ToolParam describes a single parameter a tool accepts.
type ToolParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // string, number, boolean, array, object
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// ToolDefinition describes a tool the agent can call.
// Modeled after Hermes tool registry + OpenClaw skill concept.
type ToolDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  []ToolParam `json:"parameters"`
	Execute     ToolFunc    `json:"-"` // not serialized
}

// ToolFunc is the actual function that runs a tool.
type ToolFunc func(ctx context.Context, params map[string]any) (string, error)

// ToolCall represents the AI's decision to invoke a tool.
type ToolCall struct {
	Name   string         `json:"tool"` // "tool" in JSON to match the prompt format
	Params map[string]any `json:"params"`
}

// ToolResult is the outcome of executing a tool call.
type ToolResult struct {
	Name    string `json:"name"`
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}
