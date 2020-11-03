package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/operations/fish/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func SetVariable(setVariable handlers.SetVariableHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params, err := parseSetVariableArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		if err := setVariable.Run(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseSetVariableArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.SetVariableParams, error) {
	a := &handlers.SetVariableParams{}

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
		a.Scope = handlers.UniversalScope
	case "global":
		a.Scope = handlers.GlobalScope
	case "local":
		a.Scope = handlers.LocalScope
	default:
		return nil, xerrors.Errorf("Unexpected scope: %s", scope)
	}

	return a, nil
}
