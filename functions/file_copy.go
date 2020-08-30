package functions

import (
	"context"
	"io"
	"os"

	"github.com/spf13/afero"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func FileCopy(ctx context.Context, fs afero.Fs) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		src, dest, err := parseFileCopyFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		srcFile, err := fs.Open(src)
		if err != nil {
			return starlark.False, xerrors.Errorf("Failed to open src file: %w", err)
		}
		destFile, err := fs.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return starlark.False, xerrors.Errorf("Failed to open dest file: %w", err)
		}
		if _, err := io.Copy(destFile, srcFile); err != nil {
			return starlark.False, xerrors.Errorf("Failed to copy src to dest: %w", err)
		}
		return starlark.True, nil
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
