package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/network/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func HTTPRequest(handler handlers.HTTPRequestHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)
		params, err := parseHTTPRequestArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		zap.L().Debug(
			"params",
			zap.String("url", params.URL),
			zap.String("path", params.Path),
		)
		ui.Infof("HTTP requesting. URL: %s, Path: %s", params.URL, params.Path)
		if err := handler.HTTPRequest(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseHTTPRequestArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.HTTPRequestParams, error) {
	a := &handlers.HTTPRequestParams{}
	err := lib.UnpackArgs(b.Name(), args, kwargs, "url", &a.URL, "path", &a.Path)
	if err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	return a, nil
}
