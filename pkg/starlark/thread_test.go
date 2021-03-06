package starlark_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
)

func TestThreadOptions(t *testing.T) {
	ctx := context.Background()
	thread := starlark.NewThread(
		"test",
		starlark.WithContext(ctx),
		starlark.WithLoad(func(thread *lib.Thread, module string) (lib.StringDict, error) {
			return nil, nil
		}),
	)
	got := starlark.ToContext(thread)
	assert.Equalf(t, ctx, got, "different context")

	load := thread.Load
	assert.NotNil(t, load)
}
