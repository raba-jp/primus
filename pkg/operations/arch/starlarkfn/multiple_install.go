package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/arch/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func MultipleInstall(multipleInstall handlers.MultipleInstallHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params, err := parseMultipleInstallArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		zap.L().Debug(
			"params",
			zap.Strings("names", params.Names),
		)
		ui.Infof("Installing package. Names: %s\n", params.Names)
		if err := multipleInstall.Run(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseMultipleInstallArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.MultipleInstallParams, error) {
	a := &handlers.MultipleInstallParams{}

	list := &lib.List{}
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "names", &list); err != nil {
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
	a.Names = values

	return a, nil
}
