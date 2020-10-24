package handlers

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"
)

type StartServiceHandler interface {
	StartService(ctx context.Context, dryrun bool, name string) (err error)
}

type StartServiceHandlerFunc func(ctx context.Context, dryrun bool, name string) error

func (f StartServiceHandlerFunc) StartService(ctx context.Context, dryrun bool, name string) error {
	return f(ctx, dryrun, name)
}

func NewStartService(execIF exec.Interface) StartServiceHandler {
	return StartServiceHandlerFunc(func(ctx context.Context, dryrun bool, name string) error {
		if dryrun {
			ui.Printf("systemctl start %s\n", name)
			return nil
		}

		out, err := execIF.CommandContext(ctx, "systemctl", "is-active", name).Output()
		if err != nil {
			return xerrors.Errorf("systemd service active check failed: %w", err)
		}

		if string(out) == "active" {
			return nil
		}

		if err := execIF.CommandContext(ctx, "systemctl", "start", name).Run(); err != nil {
			return xerrors.Errorf("systemd service start failed: %w", err)
		}
		return nil
	})
}
