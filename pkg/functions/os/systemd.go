package os

import (
	"bytes"
	"context"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/rs/zerolog/log"
	lib "go.starlark.net/starlark"
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

func SystemdEnable(execute backend.Execute) SystemdEnableRunner {
	return func(ctx context.Context, name string) error {
		bufout := new(bytes.Buffer)
		buferr := new(bytes.Buffer)
		if err := execute(ctx, &backend.ExecuteParams{
			Cmd:    "systemctl",
			Args:   []string{"is-enabled", name},
			Stdout: bufout,
			Stderr: buferr,
		}); err != nil {
			return xerrors.Errorf("systemctl is-enabled %s failed", name)
		}

		log.Ctx(ctx).Debug().
			Strs("command", []string{"systemctl", "is-enabled", name}).
			Str("stdout", bufout.String()).
			Str("stderr", buferr.String()).
			Msg("command output")
		if bufout.String() == "enabled\n" {
			log.Ctx(ctx).Info().Str("name", name).Msg("already enabled")
			return nil
		}

		bufout.Reset()
		buferr.Reset()
		if err := execute(ctx, &backend.ExecuteParams{
			Cmd:    "systemctl",
			Args:   []string{"enable", name},
			Stdout: bufout,
			Stderr: buferr,
		}); err != nil {
			log.Ctx(ctx).Error().
				Str("name", name).
				Str("stdout", bufout.String()).
				Str("stderr", buferr.String()).
				Msg("systemd service enable failed")
			return xerrors.Errorf("systemd service enable failed: %w", err)
		}

		log.Ctx(ctx).Debug().
			Strs("command", []string{"systemctl", "enable", name}).
			Str("stdout", bufout.String()).
			Str("stderr", buferr.String()).
			Msg("command output")
		log.Ctx(ctx).Info().Str("name", name).Msg("finish systemd service enable")
		return nil
	}
}

func SystemdStart(execute backend.Execute) SystemdStartRunner {
	return func(ctx context.Context, name string) error {
		bufout := new(bytes.Buffer)
		buferr := new(bytes.Buffer)
		if err := execute(ctx, &backend.ExecuteParams{
			Cmd:    "systemctl",
			Args:   []string{"is-active", name},
			Stdout: bufout,
			Stderr: buferr,
		}); err != nil {
			return xerrors.Errorf("systemctl is-active %s failed", name)
		}
		log.Ctx(ctx).Debug().
			Strs("command", []string{"systemctl", "is-active", name}).
			Str("stdout", bufout.String()).
			Str("stderr", buferr.String()).
			Msg("command output")
		if bufout.String() == "active\n" {
			log.Ctx(ctx).Info().Str("name", name).Msg("already active")
			return nil
		}

		bufout.Reset()
		buferr.Reset()
		if err := execute(ctx, &backend.ExecuteParams{
			Cmd:    "systemctl",
			Args:   []string{"start", name},
			Stdout: bufout,
			Stderr: buferr,
		}); err != nil {
			log.Ctx(ctx).Error().
				Str("name", name).
				Str("stdout", bufout.String()).
				Str("stderr", buferr.String()).
				Err(err).
				Msg("systemd service start failed")
			return xerrors.Errorf("systemd service start failed: %w", err)
		}

		log.Ctx(ctx).Debug().
			Strs("command", []string{"systemctl", "start", name}).
			Str("stdout", bufout.String()).
			Str("stderr", buferr.String()).
			Msg("command output")

		log.Ctx(ctx).Info().Str("name", name).Msg("finish systemd service start")
		return nil
	}
}
