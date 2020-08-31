package functions

import (
	"context"
	"fmt"

	"github.com/raba-jp/primus/executor"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

// Command execute external command
// Example command(command string, args []string)
func Command(ctx context.Context, exc executor.Executor) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		cmdName, cmdArgs, err := parseCommandFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}
		// TODO user + cwd
		ret, err := exc.Command(ctx, &executor.CommandParams{CmdName: cmdName, CmdArgs: cmdArgs})
		if err != nil {
			return toStarlarkBool(ret), xerrors.Errorf(": %w", err)
		}
		return toStarlarkBool(ret), nil
	}
}

func parseCommandFnArgs(b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (string, []string, error) {
	var cmdName string
	cmdArgs := &starlark.List{}
	err := starlark.UnpackArgs(b.Name(), args, kargs, "name", &cmdName, "args?", &cmdArgs)
	if err != nil {
		return "", nil, xerrors.Errorf("Failed to parse execute function args: %w", err)
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
				return "", nil, xerrors.Errorf("Faield to parse bool execute function arguments: %s:  (cause: %w)", cmdArgs.String(), err)
			}
			if eq {
				cmdArgsStr[index] = "true"
			} else {
				cmdArgsStr[index] = "false"
			}
		case "string":
			str, ok := starlark.AsString(val)
			if !ok {
				return "", nil, xerrors.Errorf("Faield to parse string execute function arguments: %s:  (cause: %w)", cmdArgs.String(), err)
			}
			cmdArgsStr[index] = str
		case "int":
			i32, err := starlark.AsInt32(val)
			if err != nil {
				return "", nil, xerrors.Errorf("Failed parse int32 execute function arguments: %s: cause(%w)", cmdArgs.String(), err)
			}
			cmdArgsStr[index] = fmt.Sprintf("%d", i32)
		case "float":
			return "", nil, xerrors.New("starlark interpreter does not support floating point")
		}

		index++
	}
	return cmdName, cmdArgsStr, nil
}
