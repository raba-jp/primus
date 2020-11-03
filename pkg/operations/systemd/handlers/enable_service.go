//go:generate mockery -outpkg=mocks -case=snake -name=EnableServiceHandler

package handlers

import (
	"bytes"
	"context"

	"github.com/raba-jp/primus/pkg/ctxlib"
	"go.uber.org/zap"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"
)

type EnableServiceHandler interface {
	Run(ctx context.Context, name string) (err error)
}

type EnableServiceHandlerFunc func(ctx context.Context, name string) error

func (f EnableServiceHandlerFunc) Run(ctx context.Context, name string) error {
	return f(ctx, name)
}

func NewEnableService(exc exec.Interface) EnableServiceHandler {
	return EnableServiceHandlerFunc(func(ctx context.Context, name string) error {
		if dryrun := ctxlib.DryRun(ctx); dryrun {
			enableServiceDryRun(name)
			return nil
		}
		return enableService(ctx, exc, name)
	})
}

func enableService(ctx context.Context, exc exec.Interface, name string) error {
	ctx, logger := ctxlib.LoggerWithNamespace(ctx, "enable_service")

	buferr := new(bytes.Buffer)
	cmd := exc.CommandContext(ctx, "systemctl", "is-enabled", name)
	cmd.SetStderr(buferr)

	check, _ := cmd.Output()
	logger.Debug(
		"Command output",
		zap.Strings("command", []string{"systemctl", "is-enabled", name}),
		zap.String("stdout", string(check)),
		zap.String("stderr", buferr.String()),
	)
	if string(check) == "enabled\n" {
		logger.Info("Already enabled", zap.String("name", name))
		return nil
	}

	buf := new(bytes.Buffer)
	buferr.Reset()
	cmd = exc.CommandContext(ctx, "systemctl", "enable", name)
	cmd.SetStdout(buf)
	cmd.SetStderr(buferr)

	if err := cmd.Run(); err != nil {
		logger.Error("systemd service enable failed",
			zap.String("name", name),
			zap.String("stdout", buf.String()),
			zap.String("stderr", buferr.String()),
		)
		return xerrors.Errorf("systemd service enable failed: %w", err)
	}

	logger.Debug(
		"Command output",
		zap.Strings("command", []string{"systemctl", "enable", name}),
		zap.String("stdout", buf.String()),
		zap.String("stderr", buferr.String()),
	)
	logger.Info("Finish systemd service enable", zap.String("name", name))
	return nil
}

func enableServiceDryRun(name string) {
	ui.Printf("systemctl enable %s\n", name)
}
