package starlarkfn

import (
	"github.com/raba-jp/primus/pkg/operations/command/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func Executable(handler handlers.ExecutableHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		name := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &name); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		if ret := handler.Executable(ctx, name); ret {
			return lib.True, nil
		}
		return lib.False, nil
	}
}