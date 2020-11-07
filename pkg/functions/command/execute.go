package command

import (
	"bytes"
	"context"
	"fmt"
	"os/user"
	"strconv"
	"syscall"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/exec"
	"go.uber.org/zap"

	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type Params struct {
	Name string
	Args []string
	Cwd  string
	User string
}

func (p *Params) String() string {
	return fmt.Sprintf("Cwd: %s, User: %s, %s %v", p.Cwd, p.Name, p.Name, p.Args)
}

type ExecuteRunner func(ctx context.Context, params *Params) error

func NewExecuteFunction(exc ExecuteRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "function/execute")

		params, err := parseArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		logger.Debug(
			"Params",
			zap.String("cmd", params.Name),
			zap.Strings("args", params.Args),
			zap.String("user", params.User),
			zap.String("cwd", params.Cwd),
		)

		ui.Infof("Executing command: %s\n", params)
		if err := exc(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*Params, error) {
	a := &Params{}

	cmdArgs := &lib.List{}
	err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &a.Name, "args?", &cmdArgs, "user?", &a.User, "cwd?", &a.Cwd)
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

func Execute(exc exec.Interface) ExecuteRunner {
	return func(ctx context.Context, params *Params) error {
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "execute")

		cmd := exc.CommandContext(ctx, params.Name, params.Args...)

		buf := new(bytes.Buffer)
		errbuf := new(bytes.Buffer)
		cmd.SetStdout(buf)
		cmd.SetStderr(errbuf)

		if params.Cwd != "" {
			logger.Debug("Set directory", zap.String("cwd", params.Cwd))
			cmd.SetDir(params.Cwd)
		}

		if params.User != "" {
			user, err := getUser(params.User)
			if err != nil {
				return xerrors.Errorf("Failed to get User: %w", err)
			}

			uid, err := getUID(user)
			if err != nil {
				return xerrors.Errorf("Failed to get UID: %w", err)
			}

			gid, err := getGID(user)
			if err != nil {
				return xerrors.Errorf("Failed to get GID: %w", err)
			}

			logger.Debug("Set UID and GID", zap.Uint32("uid", uid), zap.Uint32("gid", gid))
			proc := &syscall.SysProcAttr{}
			proc.Credential = &syscall.Credential{Uid: uid, Gid: gid}
			cmd.SetSysProcAttr(proc)
		}

		if err := cmd.Run(); err != nil {
			return xerrors.Errorf("Failed to execute command '%s': %w", params, err)
		}
		logger.Debug(
			"Command output",
			zap.String("stdout", buf.String()),
			zap.String("stderr", errbuf.String()),
		)
		logger.Info(
			"Executed command",
			zap.String("cmd", params.Name),
			zap.Strings("args", params.Args),
			zap.String("user", params.User),
			zap.String("cwd", params.Cwd),
		)
		return nil
	}
}

func getUser(name string) (*user.User, error) {
	u, err := user.Lookup(name)
	if err != nil {
		return nil, xerrors.Errorf("Failed to lookup user: %w", err)
	}
	return u, nil
}

func getUID(u *user.User) (uint32, error) {
	uid, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return 0, xerrors.Errorf("%w", err)
	}
	return uint32(uid), nil
}

func getGID(u *user.User) (uint32, error) {
	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		return 0, xerrors.Errorf("%w", err)
	}
	return uint32(gid), nil
}
