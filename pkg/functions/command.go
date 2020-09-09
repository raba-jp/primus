package functions

import (
	"fmt"

	"github.com/raba-jp/primus/pkg/executor"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

// Command execute external command
// Example command(command string, args []string, user string, cwd string)
func Command(exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		ctx := starlarklib.GetCtx(thread)
		params, err := parseCommandFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}
		ret, err := exc.Command(ctx, params)
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}

func parseCommandFnArgs(b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (*executor.CommandParams, error) {
	var cmdName string
	cmdArgs := &starlark.List{}
	var user string
	var cwd string
	err := starlark.UnpackArgs(b.Name(), args, kargs, "name", &cmdName, "args?", &cmdArgs, "user?", &user, "cwd?", &cwd)
	if err != nil {
		return nil, xerrors.Errorf("Failed to parse execute function args: %w", err)
	}

	iter := cmdArgs.Iterate()
	defer iter.Done()

	cmdArgsStr := make([]string, cmdArgs.Len())
	index := 0
	var val starlark.Value

	for iter.Next(&val) {
		switch val.Type() {
		case "bool":
			eq, err := starlark.Equal(val, starlark.True)
			if err != nil {
				return nil, xerrors.Errorf("Faield to parse bool execute function arguments: %s:  (cause: %w)", cmdArgs.String(), err)
			}
			if eq {
				cmdArgsStr[index] = "true"
			} else {
				cmdArgsStr[index] = "false"
			}
		case "string":
			str, ok := starlark.AsString(val)
			if !ok {
				return nil, xerrors.Errorf("Faield to parse string execute function arguments: %s:  (cause: %w)", cmdArgs.String(), err)
			}
			cmdArgsStr[index] = str
		case "int":
			i32, err := starlark.AsInt32(val)
			if err != nil {
				return nil, xerrors.Errorf("Failed parse int32 execute function arguments: %s: cause(%w)", cmdArgs.String(), err)
			}
			cmdArgsStr[index] = fmt.Sprintf("%d", i32)
		case "float":
			return nil, xerrors.New("starlark interpreter does not support floating point")
		}

		index++
	}
	return &executor.CommandParams{
		CmdName: cmdName,
		CmdArgs: cmdArgsStr,
		User:    user,
		Cwd:     cwd,
	}, nil
}
