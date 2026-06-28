// Package bootstrap wires up the MCP server with all tools.
package bootstrap

import (
	"log/slog"
	"os"
	"sync"

	"github.com/mark3labs/mcp-go/server"

	"github.com/aldok10/zara-jira-mcp/apps/api/internal/mcp"
	jira_mcp "github.com/aldok10/zara-jira-mcp/modules/jira/interfaces/mcp"
	"github.com/aldok10/zara-jira-mcp/modules/jira/application/service"
	"github.com/aldok10/zara-jira-mcp/modules/jira/infrastructure/client"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
)

// simpleCache implements domain.Cache for Jira service.
type simpleCache struct {
	mu   sync.RWMutex
	data map[string][]byte
	ttl  int
}

func newSimpleCache(ttl int) *simpleCache {
	return &simpleCache{
		data: make(map[string][]byte),
		ttl:  ttl,
	}
}

func (c *simpleCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[key]
	return v, ok
}

func (c *simpleCache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = data
}

func (c *simpleCache) TTL() int {
	return c.ttl
}

// Run starts the MCP server with stdio transport.
func Run() error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})))
	slog.Info("starting zara-jira-mcp server")

	// Load config from environment
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// Build Jira REST client
	restClient, err := client.NewRestClient(cfg)
	if err != nil {
		return err
	}

	// Build cache
	cache := newSimpleCache(60)

	// Build Jira service
	jiraSvc := service.NewJiraService(restClient, cache)

	// Build Jira handler
	jiraHandler := jira_mcp.NewHandlers(jiraSvc)

	// Create MCP server with all tools
	s := server.NewMCPServer(
		"zara-jira-mcp",
		"0.4.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Register Jira tools
	mcp.RegisterJiraTools(s, jiraHandler)

	slog.Info("server ready, waiting for MCP connections",
		"version", "0.4.0",
	)

	// Use stdio transport for MCP
	return server.ServeStdio(s)
}
