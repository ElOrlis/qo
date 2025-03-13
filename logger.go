package qo

import (
	"context"
	"log/slog"
)

type Logger interface {
	Info(ctx context.Context, format string, args ...interface{})
	Debug(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, format string, args ...interface{})
}

type defaultLogger struct {
	*slog.Logger
}

func (d defaultLogger) Info(ctx context.Context, format string, args ...interface{}) {
	d.Logger.InfoContext(ctx, format, args...)
}

func (d defaultLogger) Debug(ctx context.Context, format string, args ...interface{}) {
	d.Logger.DebugContext(ctx, format, args...)
}

func (d defaultLogger) Error(ctx context.Context, format string, args ...interface{}) {
	d.Logger.ErrorContext(ctx, format, args...)
}
