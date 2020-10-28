package handlers

import (
	"context"

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
	Run(ctx context.Context, dryrun bool, p *MultipleInstallParams) (err error)
}

type MultipleInstallHandlerFunc func(ctx context.Context, dryrun bool, p *MultipleInstallParams) error

func (f MultipleInstallHandlerFunc) Run(ctx context.Context, dryrun bool, p *MultipleInstallParams) error {
	return f(ctx, dryrun, p)
}

func NewMultipleInstall(executable command.ExecutableHandler, exc exec.Interface) MultipleInstallHandler {
	return MultipleInstallHandlerFunc(func(ctx context.Context, dryrun bool, p *MultipleInstallParams) error {
		cmd, options := cmdArgs(ctx, executable, install, p.Names)

		if dryrun {
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
