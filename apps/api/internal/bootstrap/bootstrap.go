// Package bootstrap wires up the MCP server with all tools.
package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/server"

	"github.com/aldok10/zara-jira-mcp/apps/api/internal/mcp"
	"github.com/aldok10/zara-jira-mcp/modules/jira/application/service"
	"github.com/aldok10/zara-jira-mcp/modules/jira/infrastructure/client"
	jira_mcp "github.com/aldok10/zara-jira-mcp/modules/jira/interfaces/mcp"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
)

// secureFilePermissions ensures config directory and DB file have safe permissions.
// Config directory: 0700 (owner-only access).
// DB file: 0600 (owner-only read/write).
func secureFilePermissions(dbPath string) error {
	if dbPath == "" {
		return nil
	}

	dir := filepath.Dir(dbPath)
	// Create directory with 0700 if it doesn't exist
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create config dir %s: %w", dir, err)
	}
	// Enforce 0700 on config directory
	if err := os.Chmod(dir, 0700); err != nil {
		return fmt.Errorf("chmod config dir %s: %w", dir, err)
	}

	// Enforce 0600 on DB file if it exists
	if _, err := os.Stat(dbPath); err == nil {
		if err := os.Chmod(dbPath, 0600); err != nil {
			return fmt.Errorf("chmod db file %s: %w", dbPath, err)
		}
	}

	return nil
}

// Run starts the MCP server with stdio transport.
// It respects context cancellation for graceful shutdown.
func Run(ctx context.Context) error {
	slog.Info("starting zara-jira-mcp server")

	// Fail fast if context already cancelled (e.g. signal arrived during startup).
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context cancelled during startup: %w", err)
	}

	// Load config from environment
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Secure file permissions before using any files
	if err := secureFilePermissions(cfg.Memory.DBPath); err != nil {
		slog.Warn("file permission setup", "error", err)
	}

	// Build Jira REST client
	restClient, err := client.NewRestClient(cfg)
	if err != nil {
		return fmt.Errorf("create rest client: %w", err)
	}

	// Build Jira service
	jiraSvc := service.NewJiraService(restClient)

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

	// Register PM backup tool (self-contained, no sprint module needed)
	if cfg.Memory.DBPath != "" {
		mcp.RegisterBackupTools(s, cfg.Memory.DBPath)
		slog.Info("pm_backup registered", "db_path", cfg.Memory.DBPath)
	}

	// Register onboard wizard (always available)
	toolsCount := 5 // 4 jira + 1 onboard
	if cfg.Memory.DBPath != "" {
		toolsCount++ // +1 for pm_backup
	}
	mcp.RegisterOnboardTool(s, cfg, toolsCount)
	slog.Info("pm_onboard registered", "tools_total", toolsCount)

	slog.Info("server ready, waiting for MCP connections",
		"version", "0.4.0",
	)

	// ServeStdio blocks, so run it in a goroutine and wait for
	// either the server to exit or the context to be cancelled.
	done := make(chan error, 1)
	go func() {
		done <- server.ServeStdio(s)
	}()

	select {
	case err := <-done:
		return fmt.Errorf("serve stdio: %w", err)
	case <-ctx.Done():
		slog.Info("shutdown signal received")
		return fmt.Errorf("shutdown: %w", ctx.Err())
	}
}
