package starlarklib

import (
	"context"

	"go.starlark.net/starlark"
	"go.uber.org/zap"
)

var ctxKey = "context"

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
