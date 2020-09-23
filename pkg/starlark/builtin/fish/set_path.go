package fish

import (
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/raba-jp/primus/pkg/starlark/builtin"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func SetPath(handler handlers.FishSetPathHandler) builtin.StarlarkFn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)

		params, err := parseSetPathArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		if err := handler.FishSetPath(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseSetPathArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.FishSetPathParams, error) {
	a := &handlers.FishSetPathParams{}

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
