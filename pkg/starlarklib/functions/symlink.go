package functions

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/internal/handlers"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

// Symlink create symbolic link
// Example symlink(src string, dest string)
func Symlink(handler handlers.SymlinkHandler) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		dryrun := starlarklib.GetDryRun(thread)
		lnArgs, err := arguments.NewSymlinkArguments(b, args, kwargs)
		if err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}

		zap.L().Debug(
			"params",
			zap.String("source", lnArgs.Src),
			zap.String("destination", lnArgs.Dest),
		)
		ui.Infof("Creating symbolic link. Source: %s, Destination: %s", lnArgs.Src, lnArgs.Dest)
		if err := handler.Symlink(ctx, dryrun, &handlers.SymlinkParams{Src: lnArgs.Src, Dest: lnArgs.Dest}); err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}
		return retValue, nil
	}
}
