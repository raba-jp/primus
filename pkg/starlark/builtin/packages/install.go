package packages

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func Install(chHandler handlers.CheckInstallHandler, inHandler handlers.InstallHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)
		params, err := parseArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		zap.L().Debug(
			"params",
			zap.String("name", params.Name),
			zap.String("option", params.Option),
		)
		ui.Infof("Installing package. Name: %s, Option: %s", params.Name, params.Option)
		if installed := chHandler.CheckInstall(ctx, params.Name); installed {
			return lib.None, nil
		}
		if err := inHandler.Install(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.InstallParams, error) {
	a := &handlers.InstallParams{}
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &a.Name, "option?", &a.Option); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	return a, nil
}
