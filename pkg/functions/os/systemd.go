package os

import (
	"bytes"
	"context"

	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type SystemdEnableRunner func(ctx context.Context, name string) error

type SystemdStartRunner func(ctx context.Context, name string) error

func NewSystemdEnableFunction(runner SystemdEnableRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		name := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &name); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		if err := runner(ctx, name); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func NewSystemdStartFunction(runner SystemdStartRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		name := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &name); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		if err := runner(ctx, name); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func SystemdEnable(execute command.ExecuteRunner) SystemdEnableRunner {
	return func(ctx context.Context, name string) error {
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "enable_service")

		bufout := new(bytes.Buffer)
		buferr := new(bytes.Buffer)
		if err := execute(ctx, &command.Params{
			Cmd:    "systemctl",
			Args:   []string{"is-enabled", name},
			Stdout: bufout,
			Stderr: buferr,
		}); err != nil {
			return xerrors.Errorf("systemctl is-enabled %s failed", name)
		}

		logger.Debug(
			"Command output",
			zap.Strings("command", []string{"systemctl", "is-enabled", name}),
			zap.String("stdout", bufout.String()),
			zap.String("stderr", buferr.String()),
		)
		if bufout.String() == "enabled\n" {
			logger.Info("Already enabled", zap.String("name", name))
			return nil
		}

		bufout.Reset()
		buferr.Reset()
		if err := execute(ctx, &command.Params{
			Cmd:    "systemctl",
			Args:   []string{"enable", name},
			Stdout: bufout,
			Stderr: buferr,
		}); err != nil {
			logger.Error("systemd service enable failed",
				zap.String("name", name),
				zap.String("stdout", bufout.String()),
				zap.String("stderr", buferr.String()),
			)
			return xerrors.Errorf("systemd service enable failed: %w", err)
		}

		logger.Debug(
			"Command output",
			zap.Strings("command", []string{"systemctl", "enable", name}),
			zap.String("stdout", bufout.String()),
			zap.String("stderr", buferr.String()),
		)
		logger.Info("Finish systemd service enable", zap.String("name", name))
		return nil
	}
}

func SystemdStart(execute command.ExecuteRunner) SystemdStartRunner {
	return func(ctx context.Context, name string) error {
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "start_service")

		bufout := new(bytes.Buffer)
		buferr := new(bytes.Buffer)
		if err := execute(ctx, &command.Params{
			Cmd:    "systemctl",
			Args:   []string{"is-active", name},
			Stdout: bufout,
			Stderr: buferr,
		}); err != nil {
			return xerrors.Errorf("systemctl is-active %s failed", name)
		}
		logger.Debug(
			"Command output",
			zap.Strings("command", []string{"systemctl", "is-active", name}),
			zap.String("stdout", bufout.String()),
			zap.String("stderr", buferr.String()),
		)

		if bufout.String() == "active\n" {
			logger.Info("Already active", zap.String("name", name))
			return nil
		}

		bufout.Reset()
		buferr.Reset()
		if err := execute(ctx, &command.Params{
			Cmd:    "systemctl",
			Args:   []string{"start", name},
			Stdout: bufout,
			Stderr: buferr,
		}); err != nil {
			logger.Error(
				"systemd service start failed",
				zap.String("name", name),
				zap.String("stdout", bufout.String()),
				zap.String("stderr", buferr.String()),
			)
			return xerrors.Errorf("systemd service start failed: %w", err)
		}

		logger.Debug(
			"Command output",
			zap.Strings("command", []string{"systemctl", "start", name}),
			zap.String("stdout", bufout.String()),
			zap.String("stderr", buferr.String()),
		)
		logger.Info("Finish systemd service start", zap.String("name", name))
		return nil
	}
}
