package file

import (
	"os"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func Copy(handler handlers.FileCopyHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)
		path := starlark.GetCurrentFilePath(thread)

		params, err := parseCopyArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		// TODO: paramsにCwdを追加してhandlerでやるようにする
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
		ui.Infof("Coping file. Source: %s, Destination: %s, Permission: %v\n", params.Src, params.Dest, params.Permission.String())
		if err := handler.FileCopy(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseCopyArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.FileCopyParams, error) {
	a := &handlers.FileCopyParams{}

	var perm = 0o777
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "src", &a.Src, "dest", &a.Dest, "permission?", &perm); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	a.Permission = os.FileMode(perm)
	return a, nil
}
