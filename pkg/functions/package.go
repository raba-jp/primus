package functions

import (
	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func Package(exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		name, err := parsePackageFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		ret, err := exc.Package(ctx, &executor.PackageParams{Name: name})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}

func parsePackageFnArgs(b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (name string, err error) {
	err = starlark.UnpackArgs(b.Name(), args, kargs, "name", &name)
	if err != nil {
		return "", xerrors.Errorf("Failed to parse package function args: %w", err)
	}

	return
}
