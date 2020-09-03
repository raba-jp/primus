package apply

import (
	"context"
	"os/user"
	"strconv"
	"syscall"

	"github.com/raba-jp/primus/pkg/executor"
	"golang.org/x/xerrors"
)

func (e *applyExecutor) Command(ctx context.Context, p *executor.CommandParams) (bool, error) {
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

func newSysProcAttr(name string) (*syscall.SysProcAttr, error) {
	u, err := getUser(name)
	if err != nil {
		return nil, err
	}
	uid, err := getUID(u)
	if err != nil {
		return nil, err
	}
	gid, err := getGID(u)
	if err != nil {
		return nil, err
	}
	proc := &syscall.SysProcAttr{}
	proc.Credential = &syscall.Credential{Uid: uid, Gid: gid}

	return proc, nil
}
