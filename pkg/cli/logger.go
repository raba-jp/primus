package cli

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func EnableDebugLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func EnableLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}
