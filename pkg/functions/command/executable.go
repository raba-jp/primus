package command

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/raba-jp/primus/pkg/exec"

	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type ExecutableRunner func(ctx context.Context, cmd string) bool

func NewExecutableFunction(runner ExecutableRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		name := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &name); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		logger.Ctx(ctx).Debug().Str("name", name).Msg("params")

		ret := runner(ctx, name)
		return starlark.ToBool(ret), nil
	}
}

func Executable(exc exec.Interface) ExecutableRunner {
	return func(ctx context.Context, name string) bool {
		var executable bool
		path, err := exc.LookPath(name)
		if err != nil {
			executable = false
		} else {
			executable = true
		}

		log.Ctx(ctx).Debug().Str("name", name).Str("path", path).Bool("ok", executable).Send()
		return executable
	}
}
