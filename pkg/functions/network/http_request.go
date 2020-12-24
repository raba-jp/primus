package network

import (
	"context"
	"net/http"
	"time"

	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const timeout = 10 * time.Minute

type HTTPRequestParams struct {
	URL  string
	Path string
}

type HTTPRequestRunner func(ctx context.Context, p *HTTPRequestParams) error

func NewHTTPRequestFunction(runner HTTPRequestRunner) starlark.Fn {
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
		if err := runner(ctx, params); err != nil {
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

func parseHTTPRequestArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*HTTPRequestParams, error) {
	a := &HTTPRequestParams{}
	err := lib.UnpackArgs(b.Name(), args, kwargs, "url", &a.URL, "path", &a.Path)
	if err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	return a, nil
}

func HTTPRequest(client *http.Client, fs afero.Fs) HTTPRequestRunner {
	return func(ctx context.Context, p *HTTPRequestParams) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		req, err := http.NewRequest(http.MethodGet, p.URL, nil)
		if err != nil {
			return xerrors.Errorf("Failed to create new http request: %w", err)
		}
		req = req.WithContext(ctx)

		res, err := client.Do(req)
		if err != nil {
			return xerrors.Errorf("Failed to http request: %w", err)
		}
		defer res.Body.Close()

		if err := afero.WriteReader(fs, p.Path, res.Body); err != nil {
			return xerrors.Errorf("Failed to write response body: %w", err)
		}

		return nil
	}
}
