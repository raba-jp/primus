package functions

import (
	"context"

	"github.com/spf13/afero"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

// Symlink create symbolic link
// Example symlink(src string, dest string)
func Symlink(ctx context.Context, fs afero.Fs) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		src, dest, err := parseSymlinkFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		_, err = fs.Stat(dest)
		if err == nil {
			zap.L().Info("Already exists file")
			return starlark.True, nil
		}

		linker, ok := fs.(afero.Symlinker)
		if !ok {
			return starlark.False, xerrors.New("This filesystem does not support symlink")
		}
		if err := linker.SymlinkIfPossible(src, dest); err != nil {
			return starlark.False, xerrors.Errorf("Failed to create symbolic link: %w", err)
		}

		return starlark.True, nil
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
