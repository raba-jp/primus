package starlarklib

import (
	"context"

	"go.starlark.net/starlark"
)

type ThreadOption func(*starlark.Thread)

func NewThread(name string, options ...ThreadOption) *starlark.Thread {
	thread := &starlark.Thread{
		Name: name,
	}
	for _, option := range options {
		option(thread)
	}
	return thread
}

func WithContext(ctx context.Context) ThreadOption {
	return func(thread *starlark.Thread) {
		SetCtx(ctx, thread)
	}
}

func WithLoad(loadFn func(thread *starlark.Thread, module string) (starlark.StringDict, error)) ThreadOption {
	return func(thread *starlark.Thread) {
		thread.Load = loadFn
	}
}
