package agent

import "context"

// Tool is the interface for agent-accessible tools.
// All external integrations (Jira, GitLab, Slack, Lark, etc.) are tools.
// Business logic must NOT call tools directly — use System Events.
type Tool interface {
	// Name returns the tool's unique name.
	Name() string

	// Description returns a human-readable description for the agent.
	Description() string

	// Execute runs the tool with the given arguments.
	Execute(ctx context.Context, args map[string]interface{}) (interface{}, error)

	// Schema returns the JSON Schema for the tool's arguments.
	Schema() map[string]interface{}
}

// Registry maps tool names to Tool implementations.
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry creates an empty tool registry.
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry.
func (r *ToolRegistry) Register(t Tool) {
	r.tools[t.Name()] = t
}

// Get returns a tool by name.
func (r *ToolRegistry) Get(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// List returns all registered tools.
func (r *ToolRegistry) List() []Tool {
	result := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		result = append(result, t)
	}
	return result
}

// ToolFunc adapts a function to the Tool interface.
type ToolFunc struct {
	name        string
	description string
	fn          func(ctx context.Context, args map[string]interface{}) (interface{}, error)
	schema      map[string]interface{}
}

func NewToolFunc(name, description string, fn func(ctx context.Context, args map[string]interface{}) (interface{}, error), schema map[string]interface{}) *ToolFunc {
	return &ToolFunc{
		name:        name,
		description: description,
		fn:          fn,
		schema:      schema,
	}
}

func (t *ToolFunc) Name() string        { return t.name }
func (t *ToolFunc) Description() string { return t.description }
func (t *ToolFunc) Execute(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	return t.fn(ctx, args)
}
func (t *ToolFunc) Schema() map[string]interface{} { return t.schema }
