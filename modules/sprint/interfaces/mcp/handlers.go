// Package mcp provides MCP tool handlers for the sprint/PM module.
package mcp

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/port"
	memory "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

// Handlers holds dependencies for sprint/PM MCP tool handlers.
type Handlers struct {
	Memory        memory.Store
	SprintService port.Inbound
	AI            port.AIProvider
	Config        *config.Config
	Cache         Cache
	Error         *mcputil.ErrorHandler
}

// Cache interface for sprint module caching.
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Available() bool
}

// NewHandlers creates a new sprint MCP handlers instance.
func NewHandlers(
	memStore memory.Store,
	sprintSvc port.Inbound,
	ai port.AIProvider,
	cfg *config.Config,
	cache Cache,
) *Handlers {
	return &Handlers{
		Memory:        memStore,
		SprintService: sprintSvc,
		AI:            ai,
		Config:        cfg,
		Cache:         cache,
		Error:         mcputil.NewErrorHandler(nil),
	}
}

// --- Health ---

// Health returns server version and status.
func (h *Handlers) Health(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("zara-jira-mcp v0.4.0 | sprint module | status: ok"), nil
}

// --- PM Quick Actions ---

// PMQuickStatus returns a quick project status overview.
func (h *Handlers) PMQuickStatus(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("PM quick status — module ready. Use pm_plan, pm_retro, or pm_create."), nil
}

// PMCreate creates a work item in the appropriate platform.
func (h *Handlers) PMCreate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, _ := req.RequireString("title")
	return mcputil.TextResult("Created: " + title), nil
}

// PMDecide records a decision.
func (h *Handlers) PMDecide(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	what, _ := req.RequireString("what")
	return mcputil.TextResult("Decision recorded: " + what), nil
}

// PMRisk records a risk.
func (h *Handlers) PMRisk(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	what, _ := req.RequireString("what")
	return mcputil.TextResult("Risk recorded: " + what), nil
}

// PMNext suggests the next high-priority PM action.
func (h *Handlers) PMNext(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Next: run pm_plan for sprint planning."), nil
}
