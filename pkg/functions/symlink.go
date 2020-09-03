package functions

import (
	"context"

	"github.com/raba-jp/primus/pkg/executor"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

// Symlink create symbolic link
// Example symlink(src string, dest string)
func Symlink(ctx context.Context, exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		src, dest, err := parseSymlinkFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		ret, err := exc.Symlink(ctx, &executor.SymlinkParams{Src: src, Dest: dest})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}

func parseSymlinkFnArgs(
	b *starlark.Builtin,
	args starlark.Tuple,
	kargs []starlark.Tuple,
) (src string, dest string, err error) {
	err = starlark.UnpackArgs(b.Name(), args, kargs, "src", &src, "dest", &dest)
	if err != nil {
		return "", "", xerrors.Errorf("Failed to parse symlink function args: %w", err)
	}

	return src, dest, nil
}
