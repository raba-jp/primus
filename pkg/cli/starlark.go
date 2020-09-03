package cli

import (
	"context"

	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/functions"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func ExecStarlarkFile(ctx context.Context, exc executor.Executor, path string) error {
	predeclared := starlark.StringDict{
		"execute":      starlark.NewBuiltin("execute", functions.Command(ctx, exc)),
		"symlink":      starlark.NewBuiltin("symlink", functions.Symlink(ctx, exc)),
		"http_request": starlark.NewBuiltin("http_request", functions.HTTPRequest(ctx, exc)),
		"package":      starlark.NewBuiltin("package", functions.Package(ctx, exc)),
		"file_copy":    starlark.NewBuiltin("file_copy", functions.FileCopy(ctx, exc)),
		"file_move":    starlark.NewBuiltin("file_move", functions.FileMove(ctx, exc)),
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
		zap.L().Error("Failed to exec", zap.Error(err))
		return xerrors.Errorf("Failed exec file: %w", err)
	}

	return nil
}
