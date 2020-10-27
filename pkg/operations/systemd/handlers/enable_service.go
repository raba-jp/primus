//go:generate mockery -outpkg=mocks -case=snake -name=EnableServiceHandler

package handlers

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"
)

type EnableServiceHandler interface {
	Run(ctx context.Context, dryrun bool, name string) (err error)
}

type EnableServiceHandlerFunc func(ctx context.Context, dryrun bool, name string) error

func (f EnableServiceHandlerFunc) Run(ctx context.Context, dryrun bool, name string) error {
	return f(ctx, dryrun, name)
}

func NewEnableService(exc exec.Interface) EnableServiceHandler {
	return EnableServiceHandlerFunc(func(ctx context.Context, dryrun bool, name string) error {
		if dryrun {
			ui.Printf("systemctl enable %s\n", name)
			return nil
		}

		check, err := exc.CommandContext(ctx, "systemctl", "is-enabled", name).Output()
		if err != nil {
			return xerrors.Errorf("systemd service enable check failed: %w", err)
		}
		if string(check) == "enabled" {
			return nil
		}

		if err := exc.CommandContext(ctx, "systemctl", "enable", name).Run(); err != nil {
			return xerrors.Errorf("systemd service enable failed: %w", err)
		}
		return nil
	})
}
