package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/modules/notification/domain"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

type Handlers struct {
	Notifiers map[string]domain.Notifier
	Router    domain.Router
	Error     *mcputil.ErrorHandler
}

func NewHandlers(notifiers map[string]domain.Notifier, router domain.Router) *Handlers {
	return &Handlers{
		Notifiers: notifiers,
		Router:    router,
		Error:     mcputil.NewErrorHandler(nil),
	}
}

func (h *Handlers) NotifyLark(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Notify Lark - NOT YET IMPLEMENTED"), nil
}

func (h *Handlers) NotifyRouted(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Notify routed - NOT YET IMPLEMENTED"), nil
}
