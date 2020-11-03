//go:generate mockery -outpkg=mocks -case=snake -name=InstallHandler

package handlers

import (
	"bytes"
	"context"
	"syscall"
	"time"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	command "github.com/raba-jp/primus/pkg/operations/command/handlers"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const installTimeout = 5 * time.Minute

type InstallParams struct {
	Name   string
	Option string
}

type InstallHandler interface {
	Run(ctx context.Context, dryrun bool, p *InstallParams) (err error)
}

type InstallHandlerFunc func(ctx context.Context, dryrun bool, p *InstallParams) error

func (f InstallHandlerFunc) Run(ctx context.Context, dryrun bool, p *InstallParams) error {
	return f(ctx, dryrun, p)
}

func NewInstall(checkInstall CheckInstallHandler, executable command.ExecutableHandler, exc exec.Interface) InstallHandler {
	return InstallHandlerFunc(func(ctx context.Context, dryrun bool, p *InstallParams) error {
		cmd, options := cmdArgs(ctx, executable, install, []string{p.Option, p.Name})

		if dryrun {
			cmdStr := sprintCmd(cmd, options)
			ui.Printf("%s", cmdStr)
			return nil
		}

		if installed := checkInstall.Run(ctx, p.Name); installed {
			return nil
		}

		zap.L().Debug("execute command", zap.String("cmd", cmd), zap.Strings("args", options))

		ctx, cancel := context.WithTimeout(ctx, installTimeout)
		defer cancel()

		command := exc.CommandContext(ctx, cmd, options...)

		bufout := new(bytes.Buffer)
		buferr := new(bytes.Buffer)
		command.SetStdout(bufout)
		command.SetStderr(buferr)

		if cmd == "pacman" {
			// pacman command requires root permission
			command.SetSysProcAttr(&syscall.SysProcAttr{
				Credential: &syscall.Credential{Uid: 0, Gid: 0},
			})
		}
		zap.L().Debug("execute command output", zap.String("stdout", bufout.String()), zap.String("stderr", buferr.String()))
		if err := command.Run(); err != nil {
			return xerrors.Errorf("Install package failed: Stderr: %s %w", buferr.String(), err)
		}
		return nil
	})
}
