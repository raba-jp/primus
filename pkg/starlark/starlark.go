package starlark

import (
	"context"

	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type Exec = func(ctx context.Context, path string, predeclared lib.StringDict) error

func NewExecFn(fs afero.Fs) Exec {
	return func(ctx context.Context, path string, predeclared lib.StringDict) error {
		data, err := afero.ReadFile(fs, path)
		if err != nil {
			return xerrors.Errorf("Read file failed: %s: %w", path, err)
		}

		thread := NewThread(
			"main",
			WithLoad(Load(fs, predeclared)),
			WithContext(ctx),
		)
		if _, err := lib.ExecFile(thread, path, data, predeclared); err != nil {
			return xerrors.Errorf("Failed exec file: %w", err)
		}

		return nil
	}
}

type Fn = func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kargs []lib.Tuple) (lib.Value, error)

func ExecForTest(name string, data string, fn Fn) (lib.StringDict, error) {
	predeclared := lib.StringDict{
		name: lib.NewBuiltin(name, fn),
	}
	return lib.ExecFile(NewThread("test"), "/sym/test.star", data, predeclared)
}
