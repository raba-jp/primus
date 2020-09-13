package functions

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/internal/backend"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func Package(be backend.Backend) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		pkgArgs, err := arguments.NewPackageArguments(b, args, kwargs)
		if err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}

		zap.L().Debug(
			"params",
			zap.String("name", pkgArgs.Name),
		)
		ui.Infof("Installing package. Name: %s", pkgArgs.Name)
		if installed := be.CheckInstall(ctx, pkgArgs.Name); installed {
			return retValue, nil
		}
		if err := be.Install(ctx, &backend.InstallParams{Name: pkgArgs.Name}); err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}
		return retValue, nil
	}
}
