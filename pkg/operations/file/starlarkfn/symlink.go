package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/operations/file/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

// Symlink create symbolic link
func Symlink(symlink handlers.SymlinkHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		params, err := parseSymlinkArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		zap.L().Debug(
			"params",
			zap.String("source", params.Src),
			zap.String("destination", params.Dest),
		)
		ui.Infof("Creating symbolic link. Source: %s, Destination: %s\n", params.Src, params.Dest)
		if err := symlink.Run(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseSymlinkArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.SymlinkParams, error) {
	a := &handlers.SymlinkParams{}
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "src", &a.Src, "dest", &a.Dest); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	return a, nil
}
