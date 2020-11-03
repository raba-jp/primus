package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/operations/network/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func HTTPRequest(httpRequest handlers.HTTPRequestHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "http_request")
		params, err := parseHTTPRequestArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		logger.Debug(
			"http_request Params",
			zap.String("url", params.URL),
			zap.String("path", params.Path),
		)
		if err := httpRequest.Run(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		logger.Info(
			"Finish HTTP request",
			zap.String("url", params.URL),
			zap.String("path", params.Path),
		)
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
