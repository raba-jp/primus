package starlark_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
)

func TestThreadOptions(t *testing.T) {
	ctx := context.Background()
	thread := starlark.NewThread(
		"test",
		starlark.WithContext(ctx),
		starlark.WithDryRunMode(true),
		starlark.WithLoad(func(thread *lib.Thread, module string) (lib.StringDict, error) {
			return nil, nil
		}),
	)
	if got := starlark.GetCtx(thread); ctx != got {
		t.Error("different context")
	}
	if dryrun := starlark.GetDryRunMode(thread); !dryrun {
		t.Errorf("want: true, got: %v", dryrun)
	}
	if load := thread.Load; load == nil {
		t.Error("Load is not set")
	}
}
