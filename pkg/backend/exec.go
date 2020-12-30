package backend

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"

	"github.com/rs/zerolog/log"
	"golang.org/x/xerrors"
)

type ExecuteParams struct {
	Cmd    string
	Args   []string
	Cwd    string
	User   string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Executable func(ctx context.Context, cmd string) bool

type Execute func(ctx context.Context, params *ExecuteParams) error

func NewExecutable() Executable {
	return func(ctx context.Context, name string) bool {
		e := log.Debug().Str("name", name)

		var executable bool
		path, err := exec.LookPath(name)
		if err != nil {
			e = e.Err(err).Bool("ok", false)
			executable = false
		} else {
			e = e.Str("path", path).Bool("ok", true)
			executable = true
		}

		e.Send()
		return executable
	}
}

func NewExecute() Execute {
	return func(ctx context.Context, params *ExecuteParams) error {
		e := log.Info().
			Str("cmd", params.Cmd).
			Strs("args", params.Args).
			Str("user", params.User).
			Str("cwd", params.Cwd)

		cmd := exec.CommandContext(ctx, params.Cmd, params.Args...) //nolint:gosec

		if params.Stdin != nil {
			cmd.Stdin = params.Stdin
		}
		bufout := new(bytes.Buffer)
		cmd.Stdout = bufout
		if params.Stdout != nil {
			log.Debug().Msg("set stdout")
			cmd.Stdout = io.MultiWriter(params.Stdout, bufout)
		}
		buferr := new(bytes.Buffer)
		cmd.Stderr = buferr
		if params.Stderr != nil {
			log.Debug().Msg("set stderr")
			cmd.Stderr = io.MultiWriter(params.Stderr, buferr)
		}

		if params.Cwd != "" {
			log.Debug().Str("cwd", params.Cwd).Msg("set directory")
			cmd.Dir = params.Cwd
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

			e = e.Uint32("uid", uid).Uint32("gid", gid)
			proc := &syscall.SysProcAttr{}
			proc.Credential = &syscall.Credential{Uid: uid, Gid: gid}
			cmd.SysProcAttr = proc
		}

		if err := cmd.Run(); err != nil {
			return xerrors.Errorf("Failed to execute command '%s': %w", params, err)
		}
		e = e.Str("stdout", bufout.String()).Str("stderr", buferr.String())

		e.Msg("executed command")
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
