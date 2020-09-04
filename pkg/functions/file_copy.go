package functions

import (
	"context"
	"os"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func FileCopy(ctx context.Context, exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		path := starlarklib.GetCurrentFilePath(thread)

		params, err := parseFileCopyFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		if !filepath.IsAbs(params.Src) {
			params.Src = filepath.Join(filepath.Dir(path), params.Src)
		}
		if !filepath.IsAbs(params.Dest) {
			params.Dest = filepath.Join(filepath.Dir(path), params.Dest)
		}

		zap.L().Debug(
			"params",
			zap.String("source", params.Src),
			zap.String("destination", params.Dest),
			zap.String("permission", params.Permission.String()),
		)
		ret, err := exc.FileCopy(ctx, params)
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
) (*executor.FileCopyParams, error) {
	var src string
	var dest string
	var perm int
	var err error
	if err = starlark.UnpackArgs(b.Name(), args, kargs, "src", &src, "dest", &dest, "permission?", &perm); err != nil {
		return nil, xerrors.Errorf("Failed to parse file_copy function args: %w", err)
	}

	return &executor.FileCopyParams{
		Src:        src,
		Dest:       dest,
		Permission: os.FileMode(perm),
	}, nil
}
