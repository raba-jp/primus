package functions

import (
	"context"

	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
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
		"execute":      starlark.NewBuiltin("execute", Command(exc)),
		"symlink":      starlark.NewBuiltin("symlink", Symlink(exc)),
		"http_request": starlark.NewBuiltin("http_request", HTTPRequest(exc)),
		"package":      starlark.NewBuiltin("package", Package(exc)),
		"file_copy":    starlark.NewBuiltin("file_copy", FileCopy(exc)),
		"file_move":    starlark.NewBuiltin("file_move", FileMove(exc)),
	}
	fs := afero.NewOsFs()
	data, err := afero.ReadFile(fs, path)
	if err != nil {
		return xerrors.Errorf("Failed to read file: %w", err)
	}
	thread := &starlark.Thread{
		Name: "main",
		Load: Load(fs, predeclared),
	}
	starlarklib.SetCtx(ctx, thread)

	_, err = starlark.ExecFile(thread, path, data, predeclared)
	if err != nil {
		return xerrors.Errorf("Failed exec file: %w", err)
	}

	return nil
}
