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
	"github.com/aldok10/zara-jira-mcp/agents"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/bus"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/github"
	ghmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/github/mcp"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/gitlab"
	glmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/gitlab/mcp"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/confluence"
	cmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/confluence/mcp"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/linear"
	lmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/linear/mcp"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/notion"
	nmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/notion/mcp"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/pagerduty"
	pdmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/pagerduty/mcp"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/calendar"
	calmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/calendar/mcp"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/clockify"
	clmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/clockify/mcp"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/sheets"
	shmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/sheets/mcp"
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
	// jiraHandler created later after memStore is available

	// --- Sprint/PM Module ---
	// memStore is shared between Jira and Sprint modules
	var memStore memory.Store
	var sprintService sprintPort.Inbound

	if cfg.Memory.DBPath != "" {
		sqliteStore, err := persistence.NewSQLiteStore(cfg.Memory.DBPath)
		if err != nil {
			return fmt.Errorf("create sqlite store: %w", err)
		}
		memStore = sqliteStore

		// --- Event Bus & Agent Layer ---
		eventBus := bus.NewInMemoryBus()
		agentDispatcher := agents.NewDispatcher()
		// Bridge connects event bus → agent dispatcher for registered agents
		agentBridge := agents.NewBusBridge(agentDispatcher)
		agentBridge.SubscribeTo(eventBus)
		if len(agentDispatcher.RegisteredAgents()) > 0 {
			slog.Info("agent layer initialized", "agents", len(agentDispatcher.RegisteredAgents()))
		} else {
			slog.Debug("agent layer initialized — no agents registered yet")
		}

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
			eventBus,
		)
		slog.Info("sprint service initialized with AI provider and event bus")
	} else {
		slog.Warn("sprint memory not configured — PM tools require PM_MEMORY_DB_PATH")
	}

	sprintHandler := sprint_mcp.NewHandlers(memStore, sprintService, nil, cfg, nil)
	slog.Info("sprint module initialized")

	// Jira handler needs memStore — created after sprint module
	jiraHandler := jira_mcp.NewHandlers(jiraSvc, memStore)
	slog.Info("jira module initialized")

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

	// --- GitLab Module ---
	glClient := gitlab.NewClient(cfg)
	glHandler := glmcp.NewHandlers(glClient)
	if glClient.Available() {
		slog.Info("gitlab module initialized")
	} else {
		slog.Warn("gitlab not configured — set GITLAB_TOKEN and GITLAB_PROJECT_ID")
	}

	// --- PagerDuty Module ---
	pdClient := pagerduty.NewClient(cfg)
	pdHandler := pdmcp.NewHandlers(pdClient)
	if pdClient.Available() {
		slog.Info("pagerduty module initialized")
	} else {
		slog.Warn("pagerduty not configured — set PAGERDUTY_API_KEY")
	}

	// --- Confluence Module ---
	cfClient := confluence.NewClient(cfg)
	cfHandler := cmcp.NewHandlers(cfClient)
	if cfClient.Available() {
		slog.Info("confluence module initialized")
	} else {
		slog.Warn("confluence not configured — set CONFLUENCE_BASE_URL and CONFLUENCE_API_TOKEN")
	}

	// --- Linear Module ---
	lnClient := linear.NewClient(cfg)
	lnHandler := lmcp.NewHandlers(lnClient)
	if lnClient.Available() {
		slog.Info("linear module initialized")
	} else {
		slog.Warn("linear not configured — set LINEAR_API_KEY")
	}

	// --- Notion Module ---
	nnClient := notion.NewClient(cfg)
	nnHandler := nmcp.NewHandlers(nnClient)
	if nnClient.Available() {
		slog.Info("notion module initialized")
	} else {
		slog.Warn("notion not configured — set NOTION_API_KEY")
	}

	// --- Calendar Module (Lark) ---
	calClient := calendar.NewClient(cfg)
	calHandler := calmcp.NewHandlers(calClient)
	if calClient.Available() {
		slog.Info("calendar module initialized")
	} else {
		slog.Warn("calendar not configured — set LARK_APP_ID and LARK_APP_SECRET")
	}

	// --- Clockify Module ---
	clClient := clockify.NewClient(cfg)
	clHandler := clmcp.NewHandlers(clClient)
	if clClient.Available() {
		slog.Info("clockify module initialized")
	} else {
		slog.Warn("clockify not configured — set CLOCKIFY_API_KEY and CLOCKIFY_WORKSPACE_ID")
	}

	// --- Sheets Module ---
	shClient := sheets.NewClient(cfg)
	shHandler := shmcp.NewHandlers(shClient)
	if shClient.Available() {
		slog.Info("sheets module initialized")
	} else {
		slog.Warn("sheets not configured — set GOOGLE_SHEETS_API_KEY")
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

	mcp.RegisterGitLabTools(s, glHandler)
	slog.Info("gitlab tools registered")

	mcp.RegisterPagerDutyTools(s, pdHandler)
	slog.Info("pagerduty tools registered")

	mcp.RegisterConfluenceTools(s, cfHandler)
	slog.Info("confluence tools registered")

	mcp.RegisterLinearTools(s, lnHandler)
	slog.Info("linear tools registered")

	mcp.RegisterNotionTools(s, nnHandler)
	slog.Info("notion tools registered")

	mcp.RegisterCalendarTools(s, calHandler)
	slog.Info("calendar tools registered")

	mcp.RegisterClockifyTools(s, clHandler)
	slog.Info("clockify tools registered")

	mcp.RegisterSheetsTools(s, shHandler)
	slog.Info("sheets tools registered")

	// Register PM backup tool (self-contained)
	if cfg.Memory.DBPath != "" {
		mcp.RegisterBackupTools(s, cfg.Memory.DBPath)
		slog.Info("pm_backup registered", "db_path", cfg.Memory.DBPath)
	}

	// Count registered tools for onboard wizard
	// Jira: 25, Sprint: 17, Notification: 5, GitHub: 10, GitLab: 9, PagerDuty: 2,
	// Confluence: 3, Linear: 3, Notion: 3, Calendar: 3, Clockify: 2, Sheets: 1,
	// Backup: 1, Onboard: 1
	toolsCount := 25 + 17 + 5 + 10 + 9 + 2 + 3 + 3 + 3 + 3 + 2 + 1 + 1 + 1 // 85 total
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
