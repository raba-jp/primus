package modules

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetGlobalLogger(level string) {
	if level == "" {
		logger := zap.NewNop()
		zap.ReplaceGlobals(logger)
		return
	}
	lv := getLevel(level)
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(lv)
	cfg.Encoding = "console"
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
	logger.Info("Logging", zap.String("level", lv.String()))
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.DebugLevel
	}
}
