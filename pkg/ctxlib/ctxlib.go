package ctxlib

import (
	"context"

	"go.uber.org/zap"
)

type key string

var (
	dryRunKey key = "dry-run-key"
	loggerKey key = "logger-key"
)

func SetDryRun(ctx context.Context, value bool) context.Context {
	return context.WithValue(ctx, dryRunKey, value)
}

func DryRun(ctx context.Context) bool {
	val, ok := ctx.Value(dryRunKey).(bool)
	if !ok {
		return false
	}
	return val
}

func SetLogger(ctx context.Context, value *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, value)
}

func Logger(ctx context.Context) (*zap.Logger, bool) {
	val, ok := ctx.Value(loggerKey).(*zap.Logger)
	return val, ok
}

func LoggerWithNamespace(ctx context.Context, namespace string) (context.Context, *zap.Logger) {
	logger, ok := Logger(ctx)
	if !ok {
		return ctx, zap.L()
	}
	logger = logger.With(zap.Namespace(namespace))
	newCtx := context.WithValue(ctx, loggerKey, logger)
	return newCtx, logger
}
