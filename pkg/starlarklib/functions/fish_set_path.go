package functions

import (
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func FishSetPath(handler handlers.FishSetPathHandler) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		dryrun := starlarklib.GetDryRun(thread)

		fishArgs, err := arguments.NewFishSetPathArguments(b, args, kwargs)
		if err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}

		if err := handler.FishSetPath(ctx, dryrun, &handlers.FishSetPathParams{Values: fishArgs.Values}); err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}
		return retValue, nil
	}
}
