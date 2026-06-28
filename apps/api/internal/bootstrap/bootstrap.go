// Package bootstrap wires up the MCP server with all tools.
package bootstrap

import (
	"context"
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/aldok10/zara-jira-mcp/apps/api/internal/mcp"
	jira_mcp "github.com/aldok10/zara-jira-mcp/modules/jira/interfaces/mcp"
	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/port"
	sprint_mcp "github.com/aldok10/zara-jira-mcp/modules/sprint/interfaces/mcp"
)

// Module is a placeholder for future DI.
var Module = struct{}{}

// NewServer creates a fully-configured MCP server with all tools registered.
func NewServer(
	jira *jira_mcp.Handlers,
	sprint *sprint_mcp.Handlers,
) *server.MCPServer {
	s := server.NewMCPServer(
		"zara-jira-mcp",
		"0.4.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Register all tool groups
	mcp.RegisterJiraTools(s, jira)
	mcp.RegisterSprintTools(s, sprint)
	// mcp.RegisterInfraTools(s, calHandler, ghHandler) // TODO: Wire in later when needed

	slog.Info("MCP server initialized with all tools",
		"version", "0.4.0",
	)

	return s
}

// placeholderSprintService provides a no-op sprint service for the initial setup.
type placeholderSprintService struct{}

func (p *placeholderSprintService) CalculateHealth(ctx context.Context, boardID int) (*port.HealthResult, error) {
	return &port.HealthResult{Score: 100, Rating: "Healthy"}, nil
}

func (p *placeholderSprintService) Forecast(ctx context.Context, boardID int, remaining int) (*port.ForecastResult, error) {
	return &port.ForecastResult{}, nil
}

func (p *placeholderSprintService) DetectAntiPatterns(ctx context.Context, boardID int) ([]port.AntiPattern, error) {
	return nil, nil
}

func (p *placeholderSprintService) VelocityTrend(ctx context.Context, boardID int) (string, error) {
	return "stabil", nil
}

// Run starts the MCP server with stdio transport.
func Run() error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))

	// Build module-specific handlers with available dependencies.
	// In production, use proper DI.
	jiraHandler := jira_mcp.NewHandlers(nil)
	sprintHandler := sprint_mcp.NewHandlers(
		nil,                            // memory store
		&placeholderSprintService{},   // sprint service
		nil,                            // AI provider
		nil,                            // config
		nil,                            // cache
	)

	srv := NewServer(jiraHandler, sprintHandler)

	// Use stdio transport for MCP
	return server.ServeStdio(srv)
}
