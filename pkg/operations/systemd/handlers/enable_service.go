//go:generate mockgen -destination mock/enable_service.go . EnableServiceHandler

package handlers

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"
)

type EnableServiceHandler interface {
	EnableService(ctx context.Context, dryrun bool, name string) error
}

type EnableServiceHandlerFunc func(ctx context.Context, dryrun bool, name string) error

func (f EnableServiceHandlerFunc) EnableService(ctx context.Context, dryrun bool, name string) error {
	return f(ctx, dryrun, name)
}

func NewEnableService(execIF exec.Interface) EnableServiceHandler {
	return EnableServiceHandlerFunc(func(ctx context.Context, dryrun bool, name string) error {
		if dryrun {
			ui.Printf("systemctl enable %s\n", name)
			return nil
		}

		check, err := execIF.CommandContext(ctx, "systemctl", "is-enabled", name).Output()
		if err != nil {
			return xerrors.Errorf("systemd service enable check failed: %w", err)
		}
		if string(check) == "enabled" {
			return nil
		}

		if err := execIF.CommandContext(ctx, "systemctl", "enable", name).Run(); err != nil {
			return xerrors.Errorf("systemd service enable failed: %w", err)
		}
		return nil
	})
}
