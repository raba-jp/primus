package functions

import (
	"context"

	"github.com/raba-jp/primus/pkg/executor"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type StarlarkFn = func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error)

func toStarlarkBool(v bool) starlark.Value {
	if v {
		return starlark.True
	}
	return starlark.False
}

func ExecStarlarkFile(ctx context.Context, exc executor.Executor, path string) error {
	predeclared := starlark.StringDict{
		"execute":      starlark.NewBuiltin("execute", Command(ctx, exc)),
		"symlink":      starlark.NewBuiltin("symlink", Symlink(ctx, exc)),
		"http_request": starlark.NewBuiltin("http_request", HTTPRequest(ctx, exc)),
		"package":      starlark.NewBuiltin("package", Package(ctx, exc)),
		"file_copy":    starlark.NewBuiltin("file_copy", FileCopy(ctx, exc)),
		"file_move":    starlark.NewBuiltin("file_move", FileMove(ctx, exc)),
	}

	fs := afero.NewOsFs()
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return xerrors.Errorf("Failed to read file: %w", err)
	}
	thread := &starlark.Thread{
		Name: "main",
	}
	_, err = starlark.ExecFile(thread, path, data, predeclared)
	if err != nil {
		return xerrors.Errorf("Failed exec file: %w", err)
	}

	return nil
}
