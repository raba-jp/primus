package starlarklib

import (
	"context"

	"go.starlark.net/starlark"
	"go.uber.org/zap"
)

var ctxKey = "context"
var dryRunKey = "dry_run"

func SetCtx(ctx context.Context, thread *starlark.Thread) {
	thread.SetLocal(ctxKey, ctx)
}

func GetCtx(thread *starlark.Thread) context.Context {
	ctx, ok := thread.Local(ctxKey).(context.Context)
	if !ok {
		zap.L().Error("assetion failed. return empty context.")
		return context.Background()
	}
	return ctx
}

func SetDryRun(thread *starlark.Thread, dryrun bool) {
	thread.SetLocal(dryRunKey, dryrun)
}

func GetDryRun(thread *starlark.Thread) bool {
	dryrun, ok := thread.Local(dryRunKey).(bool)
	if !ok {
		zap.L().Error("assetion failed. return empty context.")
		return true
	}
	return dryrun
}
