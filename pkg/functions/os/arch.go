package os

import (
	"bytes"
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/wesovilabs/koazee"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/modules"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"

	lib "go.starlark.net/starlark"
)

const (
	timeout                = 5 * time.Minute
	multipleInstallTimeout = 30 * time.Minute
)

type ArchInstallParams struct {
	Name   string
	Option string
}

type ArchInstalledRunner func(ctx context.Context, name string) bool

type ArchInstallRunner func(ctx context.Context, p *ArchInstallParams) error

type ArchMultipleInstallRunner func(ctx context.Context, ps []*ArchInstallParams) error

type ArchUninstallRunner func(ctx context.Context, name string) error

func NewIsArchFunction(detector modules.OSDetector) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		return starlark.ToBool(detector.ArchLinux(ctx)), nil
	}
}

func NewArchInstalledFunction(runner ArchInstalledRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		name := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &name); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		log.Ctx(ctx).Debug().Str("name", name).Msg("params")
		return starlark.ToBool(runner(ctx, name)), nil
	}
}

func NewArchInstallFunction(runner ArchInstallRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params := &ArchInstallParams{}
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &params.Name, "option?", &params.Option); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		log.Ctx(ctx).Debug().
			Str("name", params.Name).
			Str("option", params.Option).
			Msg("params")
		ui.Infof("Installing package. Name: %s, Option: %s\n", params.Name, params.Option)
		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func NewArchMultipleInstallFunction(runner ArchMultipleInstallRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params, err := parseArchMultipleInstallArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseArchMultipleInstallArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) ([]*ArchInstallParams, error) {
	list := &lib.List{}
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "names", &list); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	values := make([]string, 0, list.Len())

	iter := list.Iterate()
	defer iter.Done()
	var item lib.Value
	for iter.Next(&item) {
		str, ok := lib.AsString(item)
		if !ok {
			continue
		}
		values = append(values, str)
	}
	params := koazee.StreamOf(values).Map(func(v string) *ArchInstallParams {
		return &ArchInstallParams{Name: v}
	}).Out().Val().([]*ArchInstallParams)

	return params, nil
}

func NewArchUninstallFunction(runner ArchUninstallRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		name := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &name); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		ui.Printf("Uninstalling package. Name: %s\n", name)
		log.Ctx(ctx).Debug().Str("name", name).Msg("params")
		if err := runner(ctx, name); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func ArchInstalled(execRunner command.ExecuteRunner) ArchInstalledRunner {
	return func(ctx context.Context, name string) bool {
		err := execRunner(ctx, &command.Params{
			Cmd:  "pacman",
			Args: []string{"-Qg", name},
		})
		return err == nil
	}
}

func ArchInstall(executable command.ExecutableRunner, execute command.ExecuteRunner) ArchInstallRunner {
	return func(ctx context.Context, p *ArchInstallParams) error {
		cmd, options := archCmdArgs(ctx, executable, []string{p.Option, p.Name})
		previlegedAccess := ctxlib.PrevilegedAccessKey(ctx)

		if ArchInstalled(execute)(ctx, p.Name) {
			log.Ctx(ctx).Info().Msg("already installed")
			return nil
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		log.Ctx(ctx).Debug().Str("cmd", cmd).Strs("args", options).Msg("params")
		params := &command.Params{
			Cmd:  cmd,
			Args: options,
		}
		if strings.Contains(cmd, "pacman") {
			params.User = "root"
			params.Stdin = bytes.NewBufferString(previlegedAccess)
		}
		if err := execute(ctx, params); err != nil {
			return xerrors.Errorf("Install package failed: Stderr: %w", err)
		}
		return nil
	}
}

func ArchMultipleInstall(executable command.ExecutableRunner, execute command.ExecuteRunner) ArchMultipleInstallRunner {
	return func(ctx context.Context, ps []*ArchInstallParams) error {
		previlegedAccess := ctxlib.PrevilegedAccessKey(ctx)
		ctx, cancel := context.WithTimeout(ctx, multipleInstallTimeout)
		defer cancel()

		names := koazee.StreamOf(ps).Map(func(p *ArchInstallParams) string {
			return p.Name
		}).Do().Out().Val().([]string)
		cmd, options := archCmdArgs(ctx, executable, names)

		log.Ctx(ctx).Debug().Str("cmd", cmd).Strs("options", options).Msg("params")

		params := &command.Params{
			Cmd:  cmd,
			Args: options,
		}
		if strings.Contains(cmd, "pacman") {
			params.User = "root"
			params.Stdin = bytes.NewBufferString(previlegedAccess)
		}

		if err := execute(ctx, params); err != nil {
			return xerrors.Errorf("Install multiple package failed: %w", err)
		}
		return nil
	}
}

func ArchUninstall(execute command.ExecuteRunner) ArchUninstallRunner {
	return func(ctx context.Context, name string) error {
		if installed := ArchInstalled(execute)(ctx, name); !installed {
			return nil
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		if err := execute(ctx, &command.Params{
			Cmd:  "pacman",
			Args: []string{"-R", "--noconfirm", name},
		}); err != nil {
			return xerrors.Errorf("Remove package failed: %s: %w", name, err)
		}
		return nil
	}
}

func archCmdArgs(ctx context.Context, executable command.ExecutableRunner, cmds []string) (string, []string) {
	cmd := "sudo pacman"
	options := []string{"-S", "--noconfirm"}
	yay := usableYay(ctx, executable)
	if yay {
		cmd = "yay"
	}

	powerpill := usablePowerpill(ctx, executable)
	if powerpill && yay {
		opts := []string{"--pacman", "powerpill", "-S", "--noconfirm"}
		options = make([]string, 0, len(opts)+len(cmds))
		options = append(options, opts...)
	}

	for _, opt := range cmds {
		if opt == "" {
			continue
		}
		options = append(options, opt)
	}

	return cmd, options
}

func usableYay(ctx context.Context, executable command.ExecutableRunner) bool {
	return executable(ctx, "yay")
}

func usablePowerpill(ctx context.Context, executable command.ExecutableRunner) bool {
	return executable(ctx, "powerpill")
}
