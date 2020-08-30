package functions

import (
	"context"

	"github.com/spf13/afero"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func FileMove(ctx context.Context, fs afero.Fs) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		src, dest, err := parseFileMoveFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		if err := fs.Rename(src, dest); err != nil {
			return starlark.False, xerrors.Errorf("Failed to move %s to %s file: %w", src, dest, err)
		}
		return starlark.True, nil
	}
}

func parseFileMoveFnArgs(
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
