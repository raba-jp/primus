package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/darwin/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func Uninstall(uninstall handlers.UninstallHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params := &handlers.UninstallParams{}
		if err := lib.UnpackArgs(
			b.Name(), args, kwargs,
			"name", &params.Name,
			"cask?", &params.Cask,
			"cmd?", &params.Cmd,
		); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		ui.Printf("Uninstalling package. Name: %s\n", params.Name)
		if err := uninstall.Run(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}
