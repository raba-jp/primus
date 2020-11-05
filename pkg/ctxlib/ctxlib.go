package ctxlib

import (
	"context"

	"go.uber.org/zap"
)

type key string

var (
	dryRunKey     key = "dry-run-key"
	loggerKey     key = "logger-key"
	privilegedKey key = "previleged-key"
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

func Logger(ctx context.Context) *zap.Logger {
	val, ok := ctx.Value(loggerKey).(*zap.Logger)
	if !ok {
		return zap.L()
	}
	return val
}

func LoggerWithNamespace(ctx context.Context, namespace string) (context.Context, *zap.Logger) {
	logger := Logger(ctx)
	logger = logger.With(zap.Namespace(namespace))
	newCtx := context.WithValue(ctx, loggerKey, logger)
	return newCtx, logger
}

func SetPrevilegedAccessKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, privilegedKey, key)
}

func PrevilegedAccessKey(ctx context.Context) string {
	val, ok := ctx.Value(privilegedKey).(string)
	if !ok {
		return ""
	}
	return val
}
