package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/operations/fish/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func SetPath(setPath handlers.SetPathHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params, err := parseSetPathArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		if err := setPath.Run(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseSetPathArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.SetPathParams, error) {
	a := &handlers.SetPathParams{}

	list := &lib.List{}
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "values", &list); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	values := make([]string, 0, list.Len())

	iter := list.Iterate()
	defer iter.Done()
	var item lib.Value
	for iter.Next(&item) {
		str, ok := lib.AsString(item)
		if !ok {
			continue
		}
		values = append(values, str)
	}
	a.Values = values

	return a, nil
}
