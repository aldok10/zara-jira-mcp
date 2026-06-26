package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/fx"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/aldok10/zara-jira-mcp/config"
	"github.com/aldok10/zara-jira-mcp/internal/ai"
	"github.com/aldok10/zara-jira-mcp/internal/cache"
	icalendar "github.com/aldok10/zara-jira-mcp/internal/calendar"
	"github.com/aldok10/zara-jira-mcp/internal/clockify"
	"github.com/aldok10/zara-jira-mcp/internal/confluence"
	"github.com/aldok10/zara-jira-mcp/internal/dashboard"
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
	"github.com/aldok10/zara-jira-mcp/internal/webhook"
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
		provideCache,
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

func provideCache(cfg *config.Config) *cache.Client {
	ttl := 5 * time.Minute
	if cfg.Redis.TTL != "" {
		if d, err := time.ParseDuration(cfg.Redis.TTL); err == nil {
			ttl = d
		}
	}
	return cache.NewClient(cfg.Redis.URL, ttl)
}

func provideHandlers(
	cfg *config.Config,
	j *jira.RestClient, a *ai.OpenAIClient, l *lark.WebhookClient,
	okr *lark.OKRClient,
	s *islack.Client, d *idiscord.Client, t *itelegram.Client,
	te *iteams.Client, e *iemail.Client, c *confluence.Client,
	cal *icalendar.Client, gh *igithub.Client, gl *igitlab.Client, n *inotion.Client,
	lin *linear.Client, pd *pagerduty.Client, cl *clockify.Client, sh *sheets.Client,
	m *memory.SQLiteStore, rc *cache.Client,
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
		Cache:      rc,
		OKR:        okr,
	}
}

type LifecycleParams struct {
	fx.In
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Server     *transport.MCPServer
	LarkClient *lark.WebhookClient
	AI         *ai.OpenAIClient
	Memory     *memory.SQLiteStore
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
			if p.Config.Lark.BotEnabled {
				bot := lark.NewBotHandler(p.LarkClient, p.AI, p.Memory, p.Config.Lark.VerificationToken, logger)
				mux := http.NewServeMux()
				mux.Handle("/webhook/event", bot)
				logger.Info("starting lark bot webhook", "port", p.Config.Lark.BotPort)
				go func() {
					if err := http.ListenAndServe(":"+p.Config.Lark.BotPort, mux); err != nil {
						logger.Error("lark bot server error", "err", err)
					}
				}()
			}
			if p.Config.Webhook.Enabled && p.Config.Server.Transport != "stdio" {
				wh := webhook.NewHandler(p.Config.Webhook.Secret, p.Memory, logger)
				mux := http.NewServeMux()
				mux.Handle("/webhook/jira", wh)
				logger.Info("starting jira webhook receiver", "port", p.Config.Webhook.Port)
				go func() {
					if err := http.ListenAndServe(":"+p.Config.Webhook.Port, mux); err != nil {
						logger.Error("webhook server error", "err", err)
					}
				}()
			}
			if p.Config.Server.DashboardEnabled {
				dash := dashboard.New(p.Config, logger)
				if err := dash.Start(); err != nil {
					logger.Error("dashboard error", "err", err)
				}
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down zara-jira-mcp")
			return nil
		},
	})
}
