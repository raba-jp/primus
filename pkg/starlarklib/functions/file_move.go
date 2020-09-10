package functions

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func FileMove(exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		path := starlarklib.GetCurrentFilePath(thread)

		moveArgs, err := arguments.NewFileMoveArguments(b, args, kwargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		src := moveArgs.Src
		dest := moveArgs.Dest

		if !filepath.IsAbs(src) {
			src = filepath.Join(filepath.Dir(path), src)
		}
		if !filepath.IsAbs(dest) {
			dest = filepath.Join(filepath.Dir(path), dest)
		}

		zap.L().Debug(
			"params",
			zap.String("source", src),
			zap.String("destination", dest),
		)
		ui.Infof("Coping file. Source: %s, Destination: %s", src, dest)
		ret, err := exc.FileMove(ctx, &executor.FileMoveParams{Src: src, Dest: dest})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}
