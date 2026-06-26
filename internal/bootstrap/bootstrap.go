package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"go.uber.org/fx"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/aldok10/zara-jira-mcp/config"
	"github.com/aldok10/zara-jira-mcp/internal/ai"
	icalendar "github.com/aldok10/zara-jira-mcp/internal/calendar"
	"github.com/aldok10/zara-jira-mcp/internal/clockify"
	"github.com/aldok10/zara-jira-mcp/internal/confluence"
	idiscord "github.com/aldok10/zara-jira-mcp/internal/discord"
	iemail "github.com/aldok10/zara-jira-mcp/internal/email"
	igithub "github.com/aldok10/zara-jira-mcp/internal/github"
	igitlab "github.com/aldok10/zara-jira-mcp/internal/gitlab"
	"github.com/aldok10/zara-jira-mcp/internal/jira"
	"github.com/aldok10/zara-jira-mcp/internal/lark"
	"github.com/aldok10/zara-jira-mcp/internal/linear"
	"github.com/aldok10/zara-jira-mcp/internal/memory"
	inotion "github.com/aldok10/zara-jira-mcp/internal/notion"
	"github.com/aldok10/zara-jira-mcp/internal/pagerduty"
	"github.com/aldok10/zara-jira-mcp/internal/sheets"
	islack "github.com/aldok10/zara-jira-mcp/internal/slack"
	iteams "github.com/aldok10/zara-jira-mcp/internal/teams"
	itelegram "github.com/aldok10/zara-jira-mcp/internal/telegram"
	"github.com/aldok10/zara-jira-mcp/transport"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

var Module = fx.Module("bootstrap",
	fx.Provide(
		config.Load,
		jira.NewRestClient,
		ai.NewOpenAIClient,
		lark.NewWebhookClient,
		islack.NewClient,
		idiscord.NewClient,
		itelegram.NewClient,
		iteams.NewClient,
		iemail.NewClient,
		confluence.NewClient,
		icalendar.NewClient,
		igithub.NewClient,
		igitlab.NewClient,
		inotion.NewClient,
		linear.NewClient,
		pagerduty.NewClient,
		clockify.NewClient,
		sheets.NewClient,
		provideMemory,
		provideHandlers,
		transport.NewMCPServer,
	),
)

func provideMemory() (*memory.SQLiteStore, error) {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".zara-jira-mcp")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	return memory.NewSQLiteStore(filepath.Join(dir, "pm.db"))
}

func provideHandlers(
	cfg *config.Config,
	j *jira.RestClient, a *ai.OpenAIClient, l *lark.WebhookClient,
	s *islack.Client, d *idiscord.Client, t *itelegram.Client,
	te *iteams.Client, e *iemail.Client, c *confluence.Client,
	cal *icalendar.Client, gh *igithub.Client, gl *igitlab.Client, n *inotion.Client,
	lin *linear.Client, pd *pagerduty.Client, cl *clockify.Client, sh *sheets.Client,
	m *memory.SQLiteStore,
) *tools.Handlers {
	return &tools.Handlers{
		Config:     cfg,
		Jira:       j,
		AI:         a,
		Lark:       l,
		Slack:      s,
		Discord:    d,
		Telegram:   t,
		Teams:      te,
		Email:      e,
		Confluence: c,
		Calendar:   cal,
		GitHub:     gh,
		GitLab:     gl,
		Notion:     n,
		Linear:     lin,
		PagerDuty:  pd,
		Clockify:   cl,
		Sheets:     sh,
		Memory:     m,
	}
}

type LifecycleParams struct {
	fx.In
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Server     *transport.MCPServer
	Config     *config.Config
}

func Invoke(p LifecycleParams) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			switch p.Config.Server.Transport {
			case "sse":
				logger.Info("starting zara-jira-mcp via SSE", "port", p.Config.Server.Port)
				go func() {
					sseServer := mcpserver.NewSSEServer(p.Server.Server())
					if err := sseServer.Start(":" + p.Config.Server.Port); err != nil {
						logger.Error("SSE server error", "err", err)
					}
					p.Shutdowner.Shutdown() //nolint:errcheck
				}()
			case "http":
				logger.Info("starting zara-jira-mcp via StreamableHTTP", "port", p.Config.Server.Port)
				go func() {
					httpServer := mcpserver.NewStreamableHTTPServer(p.Server.Server())
					if err := httpServer.Start(":" + p.Config.Server.Port); err != nil {
						logger.Error("StreamableHTTP server error", "err", err)
					}
					p.Shutdowner.Shutdown() //nolint:errcheck
				}()
			default:
				logger.Info("starting zara-jira-mcp via stdio")
				go func() {
					stdio := mcpserver.NewStdioServer(p.Server.Server())
					if err := stdio.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
						logger.Info("server stopped", "reason", err.Error())
					}
					p.Shutdowner.Shutdown() //nolint:errcheck
				}()
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down zara-jira-mcp")
			return nil
		},
	})
}
