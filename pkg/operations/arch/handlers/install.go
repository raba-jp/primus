//go:generate mockery -outpkg=mocks -case=snake -name=InstallHandler

package handlers

import (
	"context"
	"time"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"
)

const installTimeout = 5 * time.Minute

type InstallParams struct {
	Name   string
	Option string
	Cmd    string
}

type InstallHandler interface {
	Run(ctx context.Context, dryrun bool, p *InstallParams) (err error)
}

type InstallHandlerFunc func(ctx context.Context, dryrun bool, p *InstallParams) error

func (f InstallHandlerFunc) Run(ctx context.Context, dryrun bool, p *InstallParams) error {
	return f(ctx, dryrun, p)
}

func NewInstall(checkInstall CheckInstallHandler, execIF exec.Interface) InstallHandler {
	return InstallHandlerFunc(func(ctx context.Context, dryrun bool, p *InstallParams) error {
		if dryrun {
			ui.Printf("pacman -S --noconfirm %s %s\n", p.Option, p.Name)
			return nil
		}

		if installed := checkInstall.Run(ctx, p.Name); installed {
			return nil
		}

		ctx, cancel := context.WithTimeout(ctx, installTimeout)
		defer cancel()
		if err := execIF.CommandContext(ctx, "pacman", "-S", "--noconfirm", p.Option, p.Name).Run(); err != nil {
			return xerrors.Errorf("Install package failed: %s: %w", p.Name, err)
		}
		return nil
	})
}
