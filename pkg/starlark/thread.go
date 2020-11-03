package starlark

import (
	"context"

	"github.com/raba-jp/primus/pkg/ctxlib"

	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
)

const (
	ctxKey    = "context"
	dryRunKey = "dry_run"
	loggerKey = "logger"
)

type ThreadOption func(thread *lib.Thread)

func NewThread(name string, options ...ThreadOption) *lib.Thread {
	thread := &lib.Thread{
		Name: name,
	}

	// set default value
	thread.SetLocal(dryRunKey, false)
	thread.SetLocal(ctxKey, context.Background())
	thread.SetLocal(loggerKey, zap.L())

	for _, option := range options {
		option(thread)
	}
	return thread
}

func WithContext(ctx context.Context) ThreadOption {
	return func(thread *lib.Thread) {
		thread.SetLocal(ctxKey, ctx)
	}
}

func WithLogger(logger *zap.Logger) ThreadOption {
	return func(thread *lib.Thread) {
		thread.SetLocal(loggerKey, logger)
	}
}

func WithDryRunMode(dryrun bool) ThreadOption {
	return func(thread *lib.Thread) {
		thread.SetLocal(dryRunKey, dryrun)
	}
}

func WithLoad(loadFn LoadFn) ThreadOption {
	return func(thread *lib.Thread) {
		thread.Load = loadFn
	}
}

func withTakeOverParent(parent *lib.Thread) ThreadOption {
	return func(child *lib.Thread) {
		ctx := getCtx(parent)
		dryrun := getDryRunMode(parent)

		logger := getLogger(parent)
		logger.With(zap.Namespace(child.Name))

		for _, opt := range []ThreadOption{
			WithContext(ctx),
			WithDryRunMode(dryrun),
			WithLogger(logger),
		} {
			opt(child)
		}

		child.Load = parent.Load
	}
}

func GetCurrentFilePath(thread *lib.Thread) string {
	callframe := thread.CallFrame(thread.CallStackDepth() - 1)
	path := callframe.Pos.Filename()
	getLogger(thread).Debug(
		"callframe",
		zap.Int("depth", thread.CallStackDepth()),
		zap.String("filepath", path),
	)
	return path
}

func ToContext(thread *lib.Thread) context.Context {
	ctx := getCtx(thread)
	ctx = ctxlib.SetDryRun(ctx, getDryRunMode(thread))
	ctx = ctxlib.SetLogger(ctx, getLogger(thread))

	return ctx
}

func getCtx(thread *lib.Thread) context.Context {
	ctx, ok := thread.Local(ctxKey).(context.Context)
	if !ok {
		zap.L().Warn("assetion failed. return empty context.")
		return context.Background()
	}
	return ctx
}

func getDryRunMode(thread *lib.Thread) bool {
	dryrun, ok := thread.Local(dryRunKey).(bool)
	if !ok {
		return true
	}
	return dryrun
}

func getLogger(thread *lib.Thread) *zap.Logger {
	logger, ok := thread.Local(loggerKey).(*zap.Logger)
	if !ok {
		return zap.L()
	}
	return logger
}
