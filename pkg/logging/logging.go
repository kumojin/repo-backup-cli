package logging

import (
	"context"
	"log/slog"
	"os"

	sentrySlog "github.com/getsentry/sentry-go/slog"
	slogmulti "github.com/samber/slog-multi"
)

func NewLogger(ctx context.Context) *slog.Logger {
	handler := slogmulti.Fanout(slog.NewJSONHandler(os.Stdout, nil), sentrySlog.Option{
		LogLevel:  []slog.Level{slog.LevelWarn, slog.LevelInfo},
		AddSource: true,
	}.NewSentryHandler(ctx))

	return slog.New(handler)
}
