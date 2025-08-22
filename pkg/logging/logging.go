package logging

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/getsentry/sentry-go"
)

type Logger struct {
	slogger *slog.Logger
}

func NewLogger() *Logger {
	return &Logger{
		slogger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}
}

func (l *Logger) With(args ...any) *Logger {
	l.slogger = l.slogger.With(args...)

	return l
}

func (l *Logger) Info(msg string, args ...any) {
	sentry.CaptureMessage(fmt.Sprintf(msg, args...))

	l.slogger.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	sentry.CaptureException(fmt.Errorf(msg, args...))

	l.slogger.Error(msg, args...)
}
