package functions

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

// Symlink create symbolic link
// Example symlink(src string, dest string)
func Symlink(exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		lnArgs, err := arguments.NewSymlinkArguments(b, args, kwargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		zap.L().Debug(
			"params",
			zap.String("source", lnArgs.Src),
			zap.String("destination", lnArgs.Dest),
		)
		ui.Infof("Creating symbolic link. Source: %s, Destination: %s", lnArgs.Src, lnArgs.Dest)
		ret, err := exc.Symlink(ctx, &executor.SymlinkParams{Src: lnArgs.Src, Dest: lnArgs.Dest})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}
