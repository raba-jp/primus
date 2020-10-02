//go:generate mockgen -destination mock/command.go . CommandHandler

package handlers

import (
	"bytes"
	"context"
	"fmt"
	"os/user"
	"strconv"
	"syscall"

	"github.com/k0kubun/pp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func init() {
	pp.ColoringEnabled = false
}

type CommandParams struct {
	CmdName string
	CmdArgs []string
	Cwd     string
	User    string
}

func (p *CommandParams) String() string {
	return pp.Sprint(p)
}

type CommandHandler interface {
	Command(ctx context.Context, dryrun bool, p *CommandParams) error
}

type CommandHandlerFunc func(ctx context.Context, dryrun bool, p *CommandParams) error

func (f CommandHandlerFunc) Command(ctx context.Context, dryrun bool, p *CommandParams) error {
	return f(ctx, dryrun, p)
}

func NewCommand(execIF exec.Interface) CommandHandler {
	return CommandHandlerFunc(func(ctx context.Context, dryrun bool, p *CommandParams) error {
		if dryrun {
			buf := new(bytes.Buffer)
			fmt.Fprintf(buf, "%s ", p.CmdName)
			for _, arg := range p.CmdArgs {
				fmt.Fprintf(buf, "%s ", arg)
			}
			fmt.Fprintf(buf, "\n")

			ui.Printf(buf.String())
			return nil
		}

		cmd := execIF.CommandContext(ctx, p.CmdName, p.CmdArgs...)
		buf := new(bytes.Buffer)
		errbuf := new(bytes.Buffer)
		cmd.SetStdout(buf)
		cmd.SetStderr(errbuf)
		if p.Cwd != "" {
			zap.L().Debug("Set directory", zap.String("cwd", p.Cwd))
			cmd.SetDir(p.Cwd)
		}

		if p.User != "" {
			proc, err := newSysProcAttr(p.User)
			if err != nil {
				return err
			}
			cmd.SetSysProcAttr(proc)
		}

		if err := cmd.Run(); err != nil {
			return xerrors.Errorf(
				"Failed to execute command '%s %s': %w",
				p.CmdName,
				p.CmdArgs,
				err,
			)
		}
		zap.L().Debug(
			"command output",
			zap.String("stdout", buf.String()),
			zap.String("stderr", errbuf.String()),
		)
		zap.L().Info(
			"Executed command",
			zap.String("cmd", p.CmdName),
			zap.Strings("args", p.CmdArgs),
			zap.String("user", p.User),
			zap.String("cwd", p.Cwd),
		)
		return nil
	})
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
