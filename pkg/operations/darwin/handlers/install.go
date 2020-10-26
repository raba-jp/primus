//go:generate mockery -outpkg=mocks -case=snake -name=InstallHandler

package handlers

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"
)

type InstallParams struct {
	Name   string
	Option string
	Cask   bool
	Cmd    string
}

type InstallHandler interface {
	Run(ctx context.Context, dryrun bool, p *InstallParams) (err error)
}

type InstallHandlerFunc func(ctx context.Context, dryrun bool, p *InstallParams) error

func (f InstallHandlerFunc) Run(ctx context.Context, dryrun bool, p *InstallParams) error {
	return f(ctx, dryrun, p)
}

func NewInstall(checkInstall CheckInstallHandler, exc exec.Interface) InstallHandler {
	return InstallHandlerFunc(func(ctx context.Context, dryrun bool, p *InstallParams) error {
		if dryrun {
			ui.Printf("brew install %s %s\n", p.Option, p.Name)
			return nil
		}

		if installed := checkInstall.Run(ctx, p.Name); installed {
			return nil
		}

		if err := exc.CommandContext(ctx, "brew", "install", p.Option, p.Name).Run(); err != nil {
			return xerrors.Errorf("Install package failed: %s: %w", p.Name, err)
		}
		return nil
	})
}
