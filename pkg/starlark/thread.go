package starlark

import (
	"context"

	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
)

const (
	ctxKey    = "context"
	dryrunKey = "dry_run"
)

type ThreadOption func(thread *lib.Thread)

func NewThread(name string, options ...ThreadOption) *lib.Thread {
	thread := &lib.Thread{
		Name: name,
	}

	// set default value
	thread.SetLocal(dryrunKey, false)
	thread.SetLocal(ctxKey, context.Background())

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

func WithDryRunMode(dryrun bool) ThreadOption {
	return func(thread *lib.Thread) {
		thread.SetLocal(dryrunKey, dryrun)
	}
}

func WithLoad(loadFn StarlarkLoadFn) ThreadOption {
	return func(thread *lib.Thread) {
		thread.Load = loadFn
	}
}

func withTakeOverParent(parent *lib.Thread) ThreadOption {
	return func(child *lib.Thread) {
		ctx := GetCtx(parent)
		dryrun := GetDryRunMode(parent)
		child.SetLocal(ctxKey, ctx)
		child.SetLocal(dryrunKey, dryrun)
		child.Load = parent.Load
	}
}

func GetCtx(thread *lib.Thread) context.Context {
	ctx, ok := thread.Local(ctxKey).(context.Context)
	if !ok {
		zap.L().Warn("assetion failed. return empty context.")
		return context.Background()
	}
	return ctx
}

func GetDryRunMode(thread *lib.Thread) bool {
	dryrun, ok := thread.Local(dryrunKey).(bool)
	if !ok {
		return true
	}
	return dryrun
}

func GetCurrentFilePath(thread *lib.Thread) string {
	callframe := thread.CallFrame(thread.CallStackDepth() - 1)
	path := callframe.Pos.Filename()
	zap.L().Debug(
		"callframe",
		zap.Int("depth", thread.CallStackDepth()),
		zap.String("filepath", path),
	)
	return path
}
