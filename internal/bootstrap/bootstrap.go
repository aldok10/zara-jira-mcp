package bootstrap

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"

	"go.uber.org/fx"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/aldok10/zara-jira-mcp/config"
	"github.com/aldok10/zara-jira-mcp/internal/ai"
	"github.com/aldok10/zara-jira-mcp/internal/jira"
	"github.com/aldok10/zara-jira-mcp/internal/lark"
	"github.com/aldok10/zara-jira-mcp/internal/memory"
	islack "github.com/aldok10/zara-jira-mcp/internal/slack"
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
		provideMemory,
		provideHandlers,
		transport.NewMCPServer,
	),
)

func provideMemory() (*memory.SQLiteStore, error) {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".zara-jira-mcp")
	os.MkdirAll(dir, 0o755)
	return memory.NewSQLiteStore(filepath.Join(dir, "pm.db"))
}

func provideHandlers(j *jira.RestClient, a *ai.OpenAIClient, l *lark.WebhookClient, s *islack.Client, m *memory.SQLiteStore) *tools.Handlers {
	return &tools.Handlers{
		Jira:   j,
		AI:     a,
		Lark:   l,
		Slack:  s,
		Memory: m,
	}
}

type LifecycleParams struct {
	fx.In
	Lifecycle  fx.Lifecycle
	Shutdowner fx.Shutdowner
	Server     *transport.MCPServer
}

func Invoke(p LifecycleParams) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting zara-jira-mcp server via stdio")
			go func() {
				stdio := mcpserver.NewStdioServer(p.Server.Server())
				if err := stdio.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
					logger.Info("server stopped", "reason", err.Error())
				}
				p.Shutdowner.Shutdown()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down zara-jira-mcp")
			return nil
		},
	})
}
