package functions

import (
	"context"

	"github.com/raba-jp/primus/pkg/internal/backend"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type StarlarkFn = func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error)

const retValue = starlark.None

func ExecStarlarkFile(ctx context.Context, be backend.Backend, path string) error {
	predeclared := starlark.StringDict{
		"execute":      starlark.NewBuiltin("execute", Command(be)),
		"symlink":      starlark.NewBuiltin("symlink", Symlink(be)),
		"http_request": starlark.NewBuiltin("http_request", HTTPRequest(be)),
		"package":      starlark.NewBuiltin("package", Package(be)),
		"file_copy":    starlark.NewBuiltin("file_copy", FileCopy(be)),
		"file_move":    starlark.NewBuiltin("file_move", FileMove(be)),
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
