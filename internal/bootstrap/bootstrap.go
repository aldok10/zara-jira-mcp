package bootstrap

import (
	"context"
	"log/slog"
	"os"

	"go.uber.org/fx"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/aldok10/zara-jira-mcp/config"
	"github.com/aldok10/zara-jira-mcp/internal/ai"
	"github.com/aldok10/zara-jira-mcp/internal/jira"
	"github.com/aldok10/zara-jira-mcp/internal/lark"
	"github.com/aldok10/zara-jira-mcp/transport"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

var Module = fx.Module("bootstrap",
	fx.Provide(
		config.Load,
		jira.NewRestClient,
		ai.NewOpenAIClient,
		lark.NewWebhookClient,
		provideHandlers,
		transport.NewMCPServer,
	),
)

func provideHandlers(j *jira.RestClient, a *ai.OpenAIClient, l *lark.WebhookClient) *tools.Handlers {
	return &tools.Handlers{
		Jira: j,
		AI:   a,
		Lark: l,
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
