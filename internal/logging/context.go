package logging

import (
	"context"
	"log/slog"
)

type contextKey string

const loggerContextKey contextKey = "logger"

func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	if logger == nil {
		logger = slog.Default()
	}
	return context.WithValue(ctx, loggerContextKey, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerContextKey).(*slog.Logger)
	if !ok || logger == nil {
		return slog.Default()
	}
	return logger
}
