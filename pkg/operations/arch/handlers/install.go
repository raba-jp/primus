//go:generate mockery -outpkg=mocks -case=snake -name=InstallHandler

package handlers

import (
	"context"
	"time"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	command "github.com/raba-jp/primus/pkg/operations/command/handlers"
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

		ctx, cancel := context.WithTimeout(ctx, installTimeout)
		defer cancel()
		if err := exc.CommandContext(ctx, cmd, options...).Run(); err != nil {
			return xerrors.Errorf("Install package failed: %w", err)
		}
		return nil
	})
}
