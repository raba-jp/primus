package starlarkfn

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/file/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func Move(handler handlers.MoveHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params, err := parseMoveArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		params.Cwd = filepath.Dir(starlark.GetCurrentFilePath(thread))

		zap.L().Debug(
			"params",
			zap.String("source", params.Src),
			zap.String("destination", params.Dest),
		)
		ui.Infof("Coping file. Source: %s, Destination: %s\n", params.Src, params.Dest)
		if err := handler.Move(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseMoveArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.MoveParams, error) {
	a := &handlers.MoveParams{}
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "src", &a.Src, "dest", &a.Dest); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	return a, nil
}
