package main

import (
	"log/slog"
	"os"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/aldok10/zara-jira-mcp/internal/bootstrap"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))

	app := fx.New(
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.SlogLogger{Logger: logger}
		}),
		bootstrap.Module,
		fx.Invoke(bootstrap.Invoke),
	)

	app.Run()
}
