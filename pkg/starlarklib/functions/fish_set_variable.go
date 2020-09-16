package functions

import (
	"github.com/raba-jp/primus/pkg/internal/handlers"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func FishSetVariable(handler handlers.FishSetVariableHandler) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		dryrun := starlarklib.GetDryRun(thread)

		fishArgs, err := arguments.NewFishSetVariableArguments(b, args, kwargs)
		if err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}

		var scope handlers.FishVariableScope
		switch fishArgs.Scope {
		case arguments.FishVariableUniversalScope:
			scope = handlers.FishVariableUniversalScope
		case arguments.FishVariableGlobalScope:
			scope = handlers.FishVariableGlobalScope
		case arguments.FishVariableLocalScope:
			scope = handlers.FishVariableLocalScope
		}

		if err := handler.FishSetVariable(ctx, dryrun, &handlers.FishSetVariableParams{
			Name:   fishArgs.Name,
			Value:  fishArgs.Value,
			Scope:  scope,
			Export: fishArgs.Export,
		}); err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}
		return retValue, nil
	}
}
