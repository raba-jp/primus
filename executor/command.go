package executor

import (
	"context"
	"os/user"
	"strconv"
	"syscall"

	"golang.org/x/xerrors"
)

type CommandParams struct {
	CmdName string
	CmdArgs []string
	Cwd     string
	User    string
}

func (e *executor) Command(ctx context.Context, p *CommandParams) (bool, error) {
	cmd := e.Exec.CommandContext(ctx, p.CmdName, p.CmdArgs...)
	cmd.SetStdout(e.Out)
	cmd.SetStderr(e.Errout)
	if p.Cwd != "" {
		cmd.SetDir(p.Cwd)
	}

	if p.User != "" {
		proc, err := newSysProcAttr(p.User)
		if err != nil {
			return false, err
		}
		cmd.SetSysProcAttr(proc)
	}

	if err := cmd.Run(); err != nil {
		return false, xerrors.Errorf(
			"Failed to execute command '%s %s': %w",
			p.CmdName,
			p.CmdArgs,
			err,
		)
	}
	return true, nil
}

func getUser(name string) (*user.User, error) {
	u, err := user.Lookup(name)
	if err != nil {
		return nil, xerrors.Errorf("Failed to lookup user: %w", err)
	}
	return u, nil
}

func getUid(u *user.User) (uint32, error) {
	uid, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return 0, xerrors.Errorf("%w", err)
	}
	return uint32(uid), nil
}

func getGid(u *user.User) (uint32, error) {
	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		return 0, xerrors.Errorf("%w", err)
	}
	return uint32(gid), nil
}

func newSysProcAttr(name string) (*syscall.SysProcAttr, error) {
	u, err := getUser(name)
	if err != nil {
		return nil, err
	}
	uid, err := getUid(u)
	if err != nil {
		return nil, err
	}
	gid, err := getGid(u)
	if err != nil {
		return nil, err
	}
	proc := &syscall.SysProcAttr{}
	proc.Credential = &syscall.Credential{Uid: uid, Gid: gid}

	return proc, nil
}
