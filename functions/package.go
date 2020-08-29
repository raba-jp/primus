package functions

import (
	"bytes"
	"context"

	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
	"k8s.io/utils/exec"
)

func Package(ctx context.Context, execCmd exec.Interface) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
		name, err := parsePackageFnArgs(b, args, kargs)
		if err != nil {
			return starlark.False, xerrors.Errorf(": %w", err)
		}

		cmd := execCmd.CommandContext(ctx, "pacman", "-S", "--noconfirm", name)
		buf := new(bytes.Buffer)
		cmd.SetStdout(buf)
		cmd.SetStderr(buf)
		if err := cmd.Run(); err != nil {
			return starlark.False, err
		}

		return starlark.True, nil
	}
}

func parsePackageFnArgs(b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (name string, err error) {
	err = starlark.UnpackArgs(b.Name(), args, kargs, "name", &name)
	if err != nil {
		return "", xerrors.Errorf("Failed to parse package function args", err)
	}

	return
}
