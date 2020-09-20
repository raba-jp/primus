package functions

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func FileCopy(handler handlers.FileCopyHandler) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		dryrun := starlarklib.GetDryRun(thread)
		path := starlarklib.GetCurrentFilePath(thread)

		copyArgs, err := arguments.NewFileCopyArguments(b, args, kwargs)
		if err != nil {
			return retValue, xerrors.Errorf(": %w", err)
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
		ui.Infof("Coping file. Source: %s, Destination: %s, Permission: %v", src, dest, perm)
		if err := handler.FileCopy(ctx, dryrun, &handlers.FileCopyParams{
			Src:        src,
			Dest:       dest,
			Permission: perm,
		}); err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}
		return retValue, nil
	}
}
