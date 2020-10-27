package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/darwin/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func Install(install handlers.InstallHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params := &handlers.InstallParams{}
		if err := lib.UnpackArgs(
			b.Name(), args, kwargs,
			"name", &params.Name,
			"option?", &params.Option,
			"cask?", &params.Cask,
			"cmd?", &params.Cmd,
		); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		zap.L().Debug(
			"params",
			zap.String("name", params.Name),
			zap.String("option", params.Option),
			zap.Bool("cask", params.Cask),
			zap.String("cmd", params.Cmd),
		)
		ui.Infof("Installing package. Name: %s, Option: %s, Cask: %v, Command: %s\n", params.Name, params.Option, params.Cask, params.Cmd)
		if err := install.Run(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}
