//go:generate mockery -outpkg=mocks -case=snake -name=MultipleInstallHandler

package handlers

import (
	"context"

	"github.com/raba-jp/primus/pkg/ctxlib"

	command "github.com/raba-jp/primus/pkg/operations/command/handlers"
	"golang.org/x/xerrors"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
)

type cmdType int

const (
	install cmdType = iota + 1
	uninstall
)

type MultipleInstallParams struct {
	Names []string
}

type MultipleInstallHandler interface {
	Run(ctx context.Context, p *MultipleInstallParams) (err error)
}

type MultipleInstallHandlerFunc func(ctx context.Context, p *MultipleInstallParams) error

func (f MultipleInstallHandlerFunc) Run(ctx context.Context, p *MultipleInstallParams) error {
	return f(ctx, p)
}

func NewMultipleInstall(executable command.ExecutableHandler, exc exec.Interface) MultipleInstallHandler {
	return MultipleInstallHandlerFunc(func(ctx context.Context, p *MultipleInstallParams) error {
		cmd, options := cmdArgs(ctx, executable, install, p.Names)

		if dryrun := ctxlib.DryRun(ctx); dryrun {
			cmdStr := sprintCmd(cmd, options)
			ui.Printf("%s", cmdStr)
			return nil
		}

		if err := exc.CommandContext(ctx, cmd, options...).Run(); err != nil {
			return xerrors.Errorf("Install multiple package failed: %w", err)
		}
		return nil
	})
}
