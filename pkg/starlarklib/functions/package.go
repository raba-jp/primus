package functions

import (
	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func Package(exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		pkgArgs, err := arguments.NewPackageArguments(b, args, kwargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		zap.L().Debug(
			"params",
			zap.String("name", pkgArgs.Name),
		)
		ret, err := exc.Package(ctx, &executor.PackageParams{Name: pkgArgs.Name})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}
