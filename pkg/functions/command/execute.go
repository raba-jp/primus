package command

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/user"
	"strconv"
	"syscall"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/rs/zerolog/log"

	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type Params struct {
	Cmd    string
	Args   []string
	Cwd    string
	User   string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (p *Params) String() string {
	return fmt.Sprintf("Cwd: %s, User: %s, %s %v", p.Cwd, p.User, p.Cmd, p.Args)
}

type ExecuteRunner func(ctx context.Context, params *Params) error

func NewExecuteFunction(exc ExecuteRunner) starlark.Fn {
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

func parseArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*Params, error) {
	a := &Params{}

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

func Execute(exc exec.Interface) ExecuteRunner {
	return func(ctx context.Context, params *Params) error {
		cmd := exc.CommandContext(ctx, params.Cmd, params.Args...)
		logger := log.Ctx(ctx)

		if params.Stdin != nil {
			cmd.SetStdin(params.Stdin)
		}
		bufout := new(bytes.Buffer)
		cmd.SetStdout(params.Stdout)
		if params.Stdout != nil {
			logger.Debug().Msg("set stdout")
			cmd.SetStdout(io.MultiWriter(params.Stdout, bufout))
		}
		buferr := new(bytes.Buffer)
		cmd.SetStderr(params.Stderr)
		if params.Stderr != nil {
			logger.Debug().Msg("set stderr")
			cmd.SetStderr(io.MultiWriter(params.Stderr, buferr))
		}

		if params.Cwd != "" {
			logger.Debug().Str("cwd", params.Cwd).Msg("set directory")
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

			logger.Debug().Uint32("uid", uid).Uint32("gid", gid).Msg("set UID and GID")
			proc := &syscall.SysProcAttr{}
			proc.Credential = &syscall.Credential{Uid: uid, Gid: gid}
			cmd.SetSysProcAttr(proc)
		}

		if err := cmd.Run(); err != nil {
			return xerrors.Errorf("Failed to execute command '%s': %w", params, err)
		}
		log.Ctx(ctx).Debug().Str("stdout", bufout.String()).Str("stderr", buferr.String()).Msg("command output")
		log.Ctx(ctx).Info().
			Str("cmd", params.Cmd).
			Strs("args", params.Args).
			Str("user", params.User).
			Str("cwd", params.Cwd).
			Msg("executed command")
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
