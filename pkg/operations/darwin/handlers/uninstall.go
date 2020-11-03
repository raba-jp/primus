//go:generate mockery -outpkg=mocks -case=snake -name=UninstallHandler

package handlers

import (
	"context"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"
)

type UninstallParams struct {
	Name string
	Cask bool
	Cmd  string
}

type UninstallHandler interface {
	Run(ctx context.Context, p *UninstallParams) (err error)
}

type UninstallHandlerFunc func(ctx context.Context, p *UninstallParams) error

func (f UninstallHandlerFunc) Run(ctx context.Context, p *UninstallParams) error {
	return f(ctx, p)
}

func NewUninstall(checkInstall CheckInstallHandler, exc exec.Interface) UninstallHandler {
	return UninstallHandlerFunc(func(ctx context.Context, p *UninstallParams) error {
		if dryrun := ctxlib.DryRun(ctx); dryrun {
			uninstallDryRun(p)
			return nil
		}
		return uninstall(ctx, checkInstall, exc, p)
	})
}

func uninstall(ctx context.Context, checkInstall CheckInstallHandler, exc exec.Interface, p *UninstallParams) error {
	if installed := checkInstall.Run(ctx, p.Name); !installed {
		return nil
	}

	if err := exc.CommandContext(ctx, "brew", "uninstall", p.Name).Run(); err != nil {
		return xerrors.Errorf("Remove package failed: %w", err)
	}
	return nil
}

func uninstallDryRun(p *UninstallParams) {
	ui.Printf("brew uninstall %s\n", p.Name)
}
