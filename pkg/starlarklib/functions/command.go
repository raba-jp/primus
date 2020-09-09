package functions

import (
	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

// Command execute external command
// Example command(command string, args []string, user string, cwd string)
func Command(exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		cmdArgs, err := arguments.NewCommandArgs(b, args, kwargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}
		zap.L().Debug(
			"params",
			zap.String("cmd", cmdArgs.Cmd),
			zap.Strings("args", cmdArgs.Args),
			zap.String("user", cmdArgs.User),
			zap.String("cwd", cmdArgs.Cwd),
		)
		ret, err := exc.Command(ctx, &executor.CommandParams{
			CmdName: cmdArgs.Cmd,
			CmdArgs: cmdArgs.Args,
			User:    cmdArgs.User,
			Cwd:     cmdArgs.Cmd,
		})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}
