package functions

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func HTTPRequest(handler handlers.HTTPRequestHandler) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		dryrun := starlarklib.GetDryRun(thread)
		reqArgs, err := arguments.NewHTTPRequestArguments(b, args, kwargs)
		if err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}

		zap.L().Debug(
			"params",
			zap.String("url", reqArgs.URL),
			zap.String("path", reqArgs.Path),
		)
		ui.Infof("HTTP requesting. URL: %s, Path: %s", reqArgs.URL, reqArgs.Path)
		if err := handler.HTTPRequest(ctx, dryrun, &handlers.HTTPRequestParams{URL: reqArgs.URL, Path: reqArgs.Path}); err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}
		return retValue, nil
	}
}
