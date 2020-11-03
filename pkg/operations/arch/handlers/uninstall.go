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
			ui.Printf("pacman -R --noconfirm %s\n", p.Name)
			return nil
		}

		if installed := checkInstall.Run(ctx, p.Name); !installed {
			return nil
		}

		ctx, cancel := context.WithTimeout(ctx, installTimeout)
		defer cancel()
		if err := exc.CommandContext(ctx, "pacman", "-R", "--noconfirm", p.Name).Run(); err != nil {
			return xerrors.Errorf("Remove package failed: %s: %w", p.Name, err)
		}
		return nil
	})
}
