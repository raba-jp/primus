package functions

import (
	"context"
	"net/http"

	"github.com/spf13/afero"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func HttpRequest(ctx context.Context, client *http.Client, fs afero.Fs) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		url, path, err := parseHttpRequestFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		res, err := client.Get(url)
		if err != nil {
			return starlark.False, xerrors.Errorf("Failed to http request: %w", err)
		}
		defer res.Body.Close()

		if err := afero.WriteReader(fs, path, res.Body); err != nil {
			return starlark.False, xerrors.Errorf("Failed to write response body: %w", err)
		}

		return starlark.True, nil
	}
}

func parseHttpRequestFnArgs(
	b *starlark.Builtin,
	args starlark.Tuple,
	kargs []starlark.Tuple,
) (url string, path string, err error) {
	err = starlark.UnpackArgs(b.Name(), args, kargs, "url", &url, "path", &path)
	if err != nil {
		return "", "", xerrors.Errorf("Failed to parse http_request function args: %w", err)
	}
	return
}
