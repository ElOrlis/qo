package qo

import (
	"context"
	"log/slog"
)

type Logger interface {
	Info(ctx context.Context, format string, args ...any)
	Debug(ctx context.Context, format string, args ...any)
	Error(ctx context.Context, format string, args ...any)
}

func newDefaultLogger() defaultLogger {
	return defaultLogger{slog.Default()}
}

type defaultLogger struct {
	*slog.Logger
}

func (d defaultLogger) Info(ctx context.Context, format string, args ...any) {
	d.Logger.InfoContext(ctx, format, args...)
}

func (d defaultLogger) Debug(ctx context.Context, format string, args ...any) {
	d.Logger.DebugContext(ctx, format, args...)
}

func (d defaultLogger) Error(ctx context.Context, format string, args ...any) {
	d.Logger.ErrorContext(ctx, format, args...)
}
