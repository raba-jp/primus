package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/packages/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func DarwinPkgUninstall(handler handlers.DarwinPkgUninstallHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params := &handlers.DarwinPkgUninstallParams{}
		if err := lib.UnpackArgs(
			b.Name(), args, kwargs,
			"name", &params.Name,
			"cask?", &params.Cask,
			"cmd?", &params.Cmd,
		); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		ui.Printf("Uninstalling package. Name: %s\n", params.Name)
		if err := handler.Uninstall(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func ArchPkgUninstall(handler handlers.ArchPkgUninstallHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params := &handlers.ArchPkgUninstallParams{}
		if err := lib.UnpackArgs(
			b.Name(), args, kwargs,
			"name", &params.Name,
			"cmd?", &params.Cmd,
		); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		ui.Printf("Uninstalling package. Name: %s\n", params.Name)
		if err := handler.Uninstall(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}
