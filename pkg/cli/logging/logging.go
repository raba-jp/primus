package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/xerrors"
)

func EnableLogger(level string) error {
	if level == "" {
		logger := zap.NewNop()
		zap.ReplaceGlobals(logger)
		return nil
	}
	lv, err := getLevel(level)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(lv)
	cfg.Encoding = "console"
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
	return nil
}

func getLevel(level string) (zapcore.Level, error) {
	switch level {
	case "debug":
		return zapcore.DebugLevel, nil
	case "info":
		return zapcore.InfoLevel, nil
	case "warn":
		return zapcore.WarnLevel, nil
	case "error":
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.DebugLevel, xerrors.Errorf("unknown level flag: %s", level)
	}
}
