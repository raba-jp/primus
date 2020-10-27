package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/operations/systemd/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func EnableService(enableService handlers.EnableServiceHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)
		name, err := parseArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		if err := enableService.Run(ctx, dryrun, name); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func StartService(startService handlers.StartServiceHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)
		name, err := parseArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		if err := startService.Run(ctx, dryrun, name); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (string, error) {
	name := ""
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &name); err != nil {
		return "", xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	return name, nil
}
