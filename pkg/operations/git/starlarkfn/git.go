package starlarkfn

import (
	"context"

	"golang.org/x/xerrors"

	lib "go.starlark.net/starlark"

	"github.com/raba-jp/primus/pkg/operations/git/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
)

func Clone(handler handlers.CloneHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		params := &handlers.CloneParams{}
		if err := lib.UnpackArgs(
			b.Name(), args, kwargs,
			"url", &params.URL,
			"path", &params.Path,
			"branch?", &params.Branch,
		); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse argumetns: %w", err)
		}

		if err := handler.Clone(context.Background(), false, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		return lib.None, nil
	}
}
