package command

import (
	"fmt"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/rs/zerolog/log"

	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

func NewExecuteFunction(exc backend.Execute) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params, err := parseArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		log.Ctx(ctx).Debug().
			Str("cmd", params.Cmd).
			Strs("args", params.Args).
			Str("user", params.User).
			Str("cwd", params.Cwd).
			Msg("params")

		ui.Infof("Executing command: %s\n", params)
		if err := exc(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*backend.ExecuteParams, error) {
	a := &backend.ExecuteParams{}

	cmdArgs := &lib.List{}
	err := lib.UnpackArgs(b.Name(), args, kwargs, "cmd", &a.Cmd, "args?", &cmdArgs, "user?", &a.User, "cwd?", &a.Cwd)
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
	a.Args = cmdArgsStr
	return a, nil
}
