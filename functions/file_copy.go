package functions

import (
	"context"

	"github.com/raba-jp/primus/executor"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func FileCopy(ctx context.Context, exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		src, dest, err := parseFileCopyFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		ret, err := exc.FileCopy(ctx, &executor.FileCopyParams{Src: src, Dest: dest})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}

func parseFileCopyFnArgs(
	b *starlark.Builtin,
	args starlark.Tuple,
	kargs []starlark.Tuple,
) (src string, dest string, err error) {
	err = starlark.UnpackArgs(b.Name(), args, kargs, "src", &src, "dest", &dest)
	if err != nil {
		return "", "", xerrors.Errorf("Failed to parse file_copy function args: %w", err)
	}

	return
}
