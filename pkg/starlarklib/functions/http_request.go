package functions

import (
	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func HTTPRequest(exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		reqArgs, err := arguments.NewHTTPRequestArguments(b, args, kwargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		zap.L().Debug(
			"params",
			zap.String("url", reqArgs.URL),
			zap.String("path", reqArgs.Path),
		)
		ret, err := exc.HTTPRequest(ctx, &executor.HTTPRequestParams{URL: reqArgs.URL, Path: reqArgs.Path})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}
