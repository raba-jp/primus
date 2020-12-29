package starlark

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	lib "go.starlark.net/starlark"
)

const (
	ctxKey    = "context"
	loggerKey = "logger"
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

func WithLogger(logger *zerolog.Logger) ThreadOption {
	return func(thread *lib.Thread) {
		thread.SetLocal(loggerKey, logger)
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

func withTakeOverParent(parent *lib.Thread) ThreadOption {
	return func(child *lib.Thread) {
		ctx := ToContext(parent)
		logger := getLogger(parent)

		for _, opt := range []ThreadOption{
			WithContext(ctx),
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
	getLogger(thread).Debug().
		Int("depth", thread.CallStackDepth()).
		Str("filepath", path).
		Msg("callframe")
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

func getLogger(thread *lib.Thread) *zerolog.Logger {
	logger, ok := thread.Local(loggerKey).(*zerolog.Logger)
	if !ok {
		l := log.With().Logger()
		return &l
	}
	return logger
}
