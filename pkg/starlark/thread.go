package starlark

import (
	"context"

	"github.com/rs/zerolog/log"
	lib "go.starlark.net/starlark"
)

const (
	ctxKey = "context"
)

type ThreadOption func(thread *lib.Thread)

func NewThread(name string, options ...ThreadOption) *lib.Thread {
	thread := &lib.Thread{
		Name: name,
	}

	// set default value
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

func WithLoad(loadFn LoadFn) ThreadOption {
	return func(thread *lib.Thread) {
		thread.Load = loadFn
	}
}

func SetContext(ctx context.Context, thread *lib.Thread) {
	thread.SetLocal(ctxKey, ctx)
}

func GetCurrentFilePath(thread *lib.Thread) string {
	callframe := thread.CallFrame(thread.CallStackDepth() - 1)
	path := callframe.Pos.Filename()
	log.Debug().Int("depth", thread.CallStackDepth()).Str("filepath", path).Msg("callframe")
	return path
}

func ToContext(thread *lib.Thread) context.Context {
	ctx, ok := thread.Local(ctxKey).(context.Context)
	if !ok {
		log.Warn().Msg("assetion failed. return empty context.")
		return context.Background()
	}
	return ctx
}
