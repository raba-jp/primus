package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/operations/vscode/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func InstallExtension(handler handlers.InstallExtensionHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params := &handlers.InstallExtensionParams{}
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &params.Name, "version?", &params.Version); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		if err := handler.InstallExtension(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}
