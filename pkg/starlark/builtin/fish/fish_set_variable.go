package fish

import (
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func SetVariable(handler handlers.FishSetVariableHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params, err := parseSetVariableArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		if err := handler.FishSetVariable(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseSetVariableArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.FishSetVariableParams, error) {
	a := &handlers.FishSetVariableParams{}

	var scope string
	if err := lib.UnpackArgs(
		b.Name(),
		args,
		kwargs,
		"name", &a.Name,
		"value", &a.Value,
		"scope", &scope,
		"export", &a.Export,
	); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	switch scope {
	case "universal":
		a.Scope = handlers.FishVariableUniversalScope
	case "global":
		a.Scope = handlers.FishVariableGlobalScope
	case "local":
		a.Scope = handlers.FishVariableLocalScope
	default:
		return nil, xerrors.Errorf("Unexpected scope: %s", scope)
	}

	return a, nil
}
