package arguments

import (
	"fmt"

	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

var _ Arguments = (*CommandArgs)(nil)

type CommandArgs struct {
	Arguments
	Cmd  string
	Args []string
	User string
	Cwd  string
}

func NewCommandArgs(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (*CommandArgs, error) {
	a := &CommandArgs{}
	err := a.Parse(b, args, kwargs)
	return a, err
}

func (a *CommandArgs) Parse(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) error {
	cmdArgs := &starlark.List{}
	err := starlark.UnpackArgs(b.Name(), args, kwargs, "name", &a.Cmd, "args?", &cmdArgs, "user?", &a.User, "cwd?", &a.Cwd)
	if err != nil {
		return xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	iter := cmdArgs.Iterate()
	defer iter.Done()

	cmdArgsStr := make([]string, cmdArgs.Len())
	index := 0
	var val starlark.Value

	for iter.Next(&val) {
		switch val.Type() {
		case "bool":
			eq, _ := starlark.Equal(val, starlark.True)
			if eq {
				cmdArgsStr[index] = "true"
			} else {
				cmdArgsStr[index] = "false"
			}
		case "string":
			str, _ := starlark.AsString(val)
			cmdArgsStr[index] = str
		case "int":
			i32, err := starlark.AsInt32(val)
			if err != nil {
				return xerrors.Errorf("Failed parse int32 execute function arguments: %s: cause(%w)", cmdArgs.String(), err)
			}
			cmdArgsStr[index] = fmt.Sprintf("%d", i32)
		}

		index++
	}
	a.Args = cmdArgsStr
	return nil
}
