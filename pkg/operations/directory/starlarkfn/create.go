package starlarkfn

import (
	"os"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/directory/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func Create(handler handlers.CreateHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params, err := parseCreateArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		params.Cwd = filepath.Dir(starlark.GetCurrentFilePath(thread))

		zap.L().Debug(
			"params",
			zap.String("path", params.Path),
			zap.String("permission", params.Permission.String()),
			zap.String("cwd", params.Cwd),
		)

		ui.Infof("Creating directories: %s\n", params.Path)
		if err := handler.Create(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseCreateArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.CreateParams, error) {
	a := &handlers.CreateParams{}
	var perm = 0o644
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "path", &a.Path, "permission?", &perm); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	a.Permission = os.FileMode(perm)
	return a, nil
}
