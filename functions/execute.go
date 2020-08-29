package functions

import (
	"bytes"
	"context"
	"fmt"

	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
	"k8s.io/utils/exec"
)

// Execute external command
// Example execute(command string, args []string)
func Execute(ctx context.Context, execCmd exec.Interface) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		cmdName, cmdArgs, err := parseExecuteFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		cmd := execCmd.CommandContext(ctx, cmdName, cmdArgs...)
		buf := new(bytes.Buffer)
		cmd.SetStdout(buf)
		cmd.SetStderr(buf)
		if err := cmd.Run(); err != nil {
			return starlark.False, err
		}
		return starlark.True, nil
	}
}

func parseExecuteFnArgs(b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (string, []string, error) {
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
