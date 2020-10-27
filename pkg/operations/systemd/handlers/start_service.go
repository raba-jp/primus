//go:generate mockery -outpkg=mocks -case=snake -name=StartServiceHandler

package handlers

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"
)

type StartServiceHandler interface {
	Run(ctx context.Context, dryrun bool, name string) (err error)
}

type StartServiceHandlerFunc func(ctx context.Context, dryrun bool, name string) error

func (f StartServiceHandlerFunc) Run(ctx context.Context, dryrun bool, name string) error {
	return f(ctx, dryrun, name)
}

func NewStartService(exc exec.Interface) StartServiceHandler {
	return StartServiceHandlerFunc(func(ctx context.Context, dryrun bool, name string) error {
		if dryrun {
			ui.Printf("systemctl start %s\n", name)
			return nil
		}

		out, err := exc.CommandContext(ctx, "systemctl", "is-active", name).Output()
		if err != nil {
			return xerrors.Errorf("systemd service active check failed: %w", err)
		}

		if string(out) == "active" {
			return nil
		}

		if err := exc.CommandContext(ctx, "systemctl", "start", name).Run(); err != nil {
			return xerrors.Errorf("systemd service start failed: %w", err)
		}
		return nil
	})
}
