package functions

import (
	"bytes"
	"fmt"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/internal/handlers"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/raba-jp/primus/pkg/starlarklib/arguments"
	"go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

// Command execute external command
// Example command(command string, args []string, user string, cwd string)
func Command(handler handlers.CommandHandler) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		dryrun := starlarklib.GetDryRun(thread)
		cmdArgs, err := arguments.NewCommandArgs(b, args, kwargs)
		if err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}
		zap.L().Debug(
			"params",
			zap.String("cmd", cmdArgs.Cmd),
			zap.Strings("args", cmdArgs.Args),
			zap.String("user", cmdArgs.User),
			zap.String("cwd", cmdArgs.Cwd),
		)

		buf := new(bytes.Buffer)
		for _, arg := range cmdArgs.Args {
			fmt.Fprintf(buf, " %s", arg)
		}
		fmt.Fprint(buf, "\n")
		ui.Infof("Executing command: %s%s", cmdArgs.Cmd, buf.String())
		if err := handler.Command(ctx, dryrun, &handlers.CommandParams{
			CmdName: cmdArgs.Cmd,
			CmdArgs: cmdArgs.Args,
			User:    cmdArgs.User,
			Cwd:     cmdArgs.Cwd,
		}); err != nil {
			return retValue, xerrors.Errorf(": %w", err)
		}
		return retValue, nil
	}
}
