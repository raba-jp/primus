package functions

import (
	"context"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func FileMove(ctx context.Context, exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		path := starlarklib.GetCurrentFilePath(thread)

		src, dest, err := parseFileMoveFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		if !filepath.IsAbs(src) {
			src = filepath.Join(filepath.Dir(path), src)
		}
		if !filepath.IsAbs(dest) {
			src = filepath.Join(filepath.Dir(path), dest)
		}

		ret, err := exc.FileMove(ctx, &executor.FileMoveParams{Src: src, Dest: dest})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
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
