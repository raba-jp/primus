package functions

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func CreateDirectory(handler handlers.CreateDirectoryHandler) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		dryrun := starlarklib.GetDryRun(thread)
		mkdirArgs, err := arguments.NewCreateDirectoryArguments(b, args, kwargs)
		if err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}
		zap.L().Debug(
			"params",
			zap.String("path", mkdirArgs.Path),
			zap.String("permission", mkdirArgs.Permission.String()),
		)

		ui.Infof("Creating directories: %s", mkdirArgs.Path)
		if err := handler.CreateDirectory(ctx, dryrun, &handlers.CreateDirectoryParams{
			Path:       mkdirArgs.Path,
			Permission: mkdirArgs.Permission,
		}); err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}
		return retValue, nil
	}
}
