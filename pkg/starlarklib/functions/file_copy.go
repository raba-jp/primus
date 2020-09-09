package functions

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func FileCopy(exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		path := starlarklib.GetCurrentFilePath(thread)

		copyArgs, err := arguments.NewFileCopyArguments(b, args, kwargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		src := copyArgs.Src
		dest := copyArgs.Dest
		perm := copyArgs.Perm

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
			zap.String("permission", perm.String()),
		)
		ret, err := exc.FileCopy(ctx, &executor.FileCopyParams{
			Src:        src,
			Dest:       dest,
			Permission: perm,
		})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}
