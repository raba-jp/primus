package command

import (
	"bytes"
	"fmt"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

// Command execute external command
func Command(handler handlers.CommandHandler) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.GetCtx(thread)
		dryrun := starlark.GetDryRunMode(thread)
		params, err := parseArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		zap.L().Debug(
			"params",
			zap.String("cmd", params.CmdName),
			zap.Strings("args", params.CmdArgs),
			zap.String("user", params.User),
			zap.String("cwd", params.Cwd),
		)

		buf := new(bytes.Buffer)
		for _, arg := range params.CmdArgs {
			fmt.Fprintf(buf, " %s", arg)
		}
		fmt.Fprint(buf, "\n")
		ui.Infof("Executing command: %s%s", params.CmdName, buf.String())
		if err := handler.Command(ctx, dryrun, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*handlers.CommandParams, error) {
	a := &handlers.CommandParams{}

	cmdArgs := &lib.List{}
	err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &a.CmdName, "args?", &cmdArgs, "user?", &a.User, "cwd?", &a.Cwd)
	if err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	iter := cmdArgs.Iterate()
	defer iter.Done()

	cmdArgsStr := make([]string, cmdArgs.Len())
	index := 0
	var val lib.Value

	for iter.Next(&val) {
		switch val.Type() {
		case "bool":
			eq, _ := lib.Equal(val, lib.True)
			if eq {
				cmdArgsStr[index] = "true"
			} else {
				cmdArgsStr[index] = "false"
			}
		case "string":
			str, _ := lib.AsString(val)
			cmdArgsStr[index] = str
		case "int":
			i32, err := lib.AsInt32(val)
			if err != nil {
				return nil, xerrors.Errorf("Failed parse int32 command function arguments: %s: cause(%w)", cmdArgs.String(), err)
			}
			cmdArgsStr[index] = fmt.Sprintf("%d", i32)
		}

		index++
	}
	a.CmdArgs = cmdArgsStr
	return a, nil
}
