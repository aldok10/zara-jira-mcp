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
	notifDomain "github.com/aldok10/zara-jira-mcp/modules/notification/domain"
	"github.com/aldok10/zara-jira-mcp/modules/notification/infrastructure/discord"
	"github.com/aldok10/zara-jira-mcp/modules/notification/infrastructure/lark"
	"github.com/aldok10/zara-jira-mcp/modules/notification/infrastructure/slack"
	"github.com/aldok10/zara-jira-mcp/modules/notification/infrastructure/telegram"
	notif_mcp "github.com/aldok10/zara-jira-mcp/modules/notification/interfaces/mcp"
	sprintPort "github.com/aldok10/zara-jira-mcp/modules/sprint/application/port"
	sprintSvc "github.com/aldok10/zara-jira-mcp/modules/sprint/application/service"
	"github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
	"github.com/aldok10/zara-jira-mcp/modules/sprint/infrastructure/persistence"
	"github.com/aldok10/zara-jira-mcp/modules/sprint/infrastructure/sprintstore"
	sprint_mcp "github.com/aldok10/zara-jira-mcp/modules/sprint/interfaces/mcp"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/ai"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/github"
	ghmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/github/mcp"
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

	// --- Jira Module ---
	restClient, err := client.NewRestClient(cfg)
	if err != nil {
		return fmt.Errorf("create rest client: %w", err)
	}
	jiraSvc := service.NewJiraService(restClient)
	jiraHandler := jira_mcp.NewHandlers(jiraSvc)
	slog.Info("jira module initialized")

	// --- Sprint/PM Module ---
	var memStore memory.Store
	var sprintService sprintPort.Inbound

	if cfg.Memory.DBPath != "" {
		sqliteStore, err := persistence.NewSQLiteStore(cfg.Memory.DBPath)
		if err != nil {
			return fmt.Errorf("create sqlite store: %w", err)
		}
		memStore = sqliteStore

		// Build sprint service with real dependencies
		aiProvider := ai.NewOpenAIClient(cfg)
		sprintService = sprintSvc.NewSprintService(
			sprintstore.NewSnapshotRepository(sqliteStore),
			sprintstore.NewHealthRepository(sqliteStore),
			sprintstore.NewRiskRepository(sqliteStore),
			sprintstore.NewBlockerRepository(sqliteStore),
			sprintstore.NewGoalRepository(sqliteStore),
			restClient, // jira domain.Client satisfies sprint port.JiraClient
			aiProvider,
			&sprintstore.NoopEventBus{},
		)
		slog.Info("sprint service initialized with AI provider")
	} else {
		slog.Warn("sprint memory not configured — PM tools require PM_MEMORY_DB_PATH")
	}

	sprintHandler := sprint_mcp.NewHandlers(memStore, sprintService, nil, cfg, nil)
	slog.Info("sprint module initialized")

	// --- Notification Module ---
	notifiers := make(map[string]notifDomain.Notifier)

	if sl := slack.NewClient(cfg); sl.Available() {
		notifiers["slack"] = slack.NewNotifierAdapter(sl)
		slog.Info("slack notifier registered")
	}
	if dc := discord.NewClient(cfg); dc.Available() {
		notifiers["discord"] = discord.NewNotifierAdapter(dc)
		slog.Info("discord notifier registered")
	}
	if tg := telegram.NewClient(cfg); tg.Available() {
		notifiers["telegram"] = telegram.NewNotifierAdapter(tg)
		slog.Info("telegram notifier registered")
	}
	lkW := lark.NewWebhookClient(cfg)
	if cfg.Lark.WebhookURL != "" || (cfg.Lark.AppID != "" && cfg.Lark.AppSecret != "") {
		notifiers["lark"] = lark.NewNotifierAdapter(lkW)
		slog.Info("lark notifier registered")
	}

	notifHandler := notif_mcp.NewHandlers(notifiers, nil)
	slog.Info("notification module initialized", "channels", len(notifiers))

	// --- GitHub Module ---
	ghClient := github.NewClient(cfg)
	ghHandler := ghmcp.NewHandlers(ghClient)
	if ghClient.Available() {
		slog.Info("github module initialized")
	} else {
		slog.Warn("github not configured — set GITHUB_TOKEN, GITHUB_OWNER, GITHUB_REPO")
	}

	// --- Create MCP Server ---
	s := server.NewMCPServer(
		"zara-jira-mcp",
		"0.4.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Register all tools
	mcp.RegisterJiraTools(s, jiraHandler)
	slog.Info("jira tools registered")

	mcp.RegisterSprintTools(s, sprintHandler)
	slog.Info("sprint/pm tools registered")

	mcp.RegisterNotificationTools(s, notifHandler)
	slog.Info("notification tools registered")

	mcp.RegisterGitHubTools(s, ghHandler)
	slog.Info("github tools registered")

	// Register PM backup tool (self-contained)
	if cfg.Memory.DBPath != "" {
		mcp.RegisterBackupTools(s, cfg.Memory.DBPath)
		slog.Info("pm_backup registered", "db_path", cfg.Memory.DBPath)
	}

	// Count registered tools for onboard wizard
	// Jira: 25, Sprint: 17, Notification: 5, GitHub: 10, Backup: 1, Onboard: 1
	toolsCount := 25 + 17 + 5 + 10 + 1 + 1 // 59 total
	mcp.RegisterOnboardTool(s, cfg, toolsCount)
	slog.Info("pm_onboard registered", "tools_total", toolsCount)

	slog.Info("server ready, waiting for MCP connections",
		"version", "0.4.0",
		"tools", toolsCount,
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
