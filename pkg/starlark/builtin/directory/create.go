package directory

import (
	"os"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/raba-jp/primus/pkg/starlark/builtin"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func Create(handler handlers.CreateDirectoryHandler) builtin.StarlarkFn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params, err := parseCreateArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		// TODO: paramsにCwdを追加してhandlerでやるようにする
		if !filepath.IsAbs(params.Path) {
			current := starlark.GetCurrentFilePath(thread)
			params.Path = filepath.Join(filepath.Dir(current), params.Path)
		}

		zap.L().Debug(
			"params",
			zap.String("path", params.Path),
			zap.String("permission", params.Permission.String()),
		)

		ui.Infof("Creating directories: %s", params.Path)
		if err := handler.CreateDirectory(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseCreateArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.CreateDirectoryParams, error) {
	a := &handlers.CreateDirectoryParams{}
	var perm = 0o644
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "path", &a.Path, "permission?", &perm); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	a.Permission = os.FileMode(perm)
	return a, nil
}
