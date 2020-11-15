package command

import (
	"context"

	"github.com/raba-jp/primus/pkg/exec"

	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type ExecutableRunner func(ctx context.Context, cmd string) bool

func NewExecutableFunction(runner ExecutableRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "function/executable")

		name := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &name); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		logger.Debug("Params", zap.String("name", name))

		ret := runner(ctx, name)
		return starlark.ToBool(ret), nil
	}
}

func Executable(exc exec.Interface) ExecutableRunner {
	return func(ctx context.Context, name string) bool {
		ctx, _ = ctxlib.LoggerWithNamespace(ctx, "os_detector")
		_, logger := ctxlib.LoggerWithNamespace(ctx, "executable_command")

		var executable bool
		path, err := exc.LookPath(name)
		if err != nil {
			executable = false
		} else {
			executable = true
		}

		logger.Debug(
			"Check executable",
			zap.String("name", name),
			zap.String("path", path),
			zap.Bool("ok", executable),
		)
		return executable
	}
}
