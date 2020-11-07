package filesystem

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type ExistsFileRunner func(ctx context.Context, path string) (exists bool)

func NewExistsFileFunction(runner ExistsFileRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "function/exists_file")

		path := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "path", &path); err != nil {
			return lib.False, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		logger.Debug("Params", zap.String("path", path))

		ui.Infof("Check existence file. Path: %s\n", path)
		return starlark.ToBool(runner(ctx, path)), nil
	}
}

func ExistsFile(fs afero.Fs) ExistsFileRunner {
	return func(ctx context.Context, path string) bool {
		_, logger := ctxlib.LoggerWithNamespace(ctx, "exists_file")
		_, err := fs.Stat(path)
		if err == nil {
			logger.Info("Already exists file")
			return true
		}
		return false
	}
}
