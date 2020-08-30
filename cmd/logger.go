package cmd

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func enableDebugLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func enableLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}
