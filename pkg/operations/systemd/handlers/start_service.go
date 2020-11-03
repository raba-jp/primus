//go:generate mockery -outpkg=mocks -case=snake -name=StartServiceHandler

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

type StartServiceHandler interface {
	Run(ctx context.Context, name string) (err error)
}

type StartServiceHandlerFunc func(ctx context.Context, name string) error

func (f StartServiceHandlerFunc) Run(ctx context.Context, name string) error {
	return f(ctx, name)
}

func NewStartService(exc exec.Interface) StartServiceHandler {
	return StartServiceHandlerFunc(func(ctx context.Context, name string) error {
		if dryrun := ctxlib.DryRun(ctx); dryrun {
			startServiceDryRun(name)
			return nil
		}
		return startService(ctx, exc, name)
	})
}

func startService(ctx context.Context, exc exec.Interface, name string) error {
	ctx, logger := ctxlib.LoggerWithNamespace(ctx, "start_service")

	buferr := new(bytes.Buffer)

	cmd := exc.CommandContext(ctx, "systemctl", "is-active", name)
	cmd.SetStderr(buferr)
	out, _ := cmd.Output()
	logger.Debug(
		"Command output",
		zap.Strings("command", []string{"systemctl", "is-active", name}),
		zap.String("stdout", string(out)),
		zap.String("stderr", buferr.String()),
	)

	if string(out) == "active\n" {
		logger.Info("Already active", zap.String("name", name))
		return nil
	}

	buf := new(bytes.Buffer)
	buferr.Reset()
	cmd = exc.CommandContext(ctx, "systemctl", "start", name)
	cmd.SetStdout(buf)
	cmd.SetStderr(buferr)
	if err := cmd.Run(); err != nil {
		logger.Error(
			"systemd service start failed",
			zap.String("name", name),
			zap.String("stdout", buf.String()),
			zap.String("stderr", buferr.String()),
		)
		return xerrors.Errorf("systemd service start failed: %w", err)
	}

	logger.Debug(
		"Command output",
		zap.Strings("command", []string{"systemctl", "start", name}),
		zap.String("stdout", buf.String()),
		zap.String("stderr", buferr.String()),
	)
	logger.Info("Finish systemd service start", zap.String("name", name))
	return nil
}

func startServiceDryRun(name string) {
	ui.Printf("systemctl start %s\n", name)
}
