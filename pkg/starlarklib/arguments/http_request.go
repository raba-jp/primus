package arguments

import (
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

var _ Arguments = (*HTTPRequestArguments)(nil)

type HTTPRequestArguments struct {
	Arguments
	URL  string
	Path string
}

func NewHTTPRequestArguments(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (*HTTPRequestArguments, error) {
	a := HTTPRequestArguments{}
	err := a.Parse(b, args, kwargs)
	return &a, err
}

func (a *HTTPRequestArguments) Parse(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) error {
	err := starlark.UnpackArgs(b.Name(), args, kwargs, "url", &a.URL, "path", &a.Path)
	if err != nil {
		return xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	return nil
}
