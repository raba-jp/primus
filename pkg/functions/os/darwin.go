package os

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type DarwinInstallParams struct {
	Name   string
	Option string
	Cask   bool
}

type DarwinInstalledRunner func(ctx context.Context, name string) bool

type DarwinInstallRunner func(ctx context.Context, p *DarwinInstallParams) error

type DarwinUninstallRunner func(ctx context.Context, name string) error

func NewIsDarwinFunction(checker backend.DarwinChecker) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		return starlark.ToBool(checker(ctx)), nil
	}
}

func NewDarwinInstalledFunction(runner DarwinInstalledRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		name := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &name); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		if ret := runner(ctx, name); ret {
			return lib.True, nil
		}
		return lib.False, nil
	}
}

func NewDarwinInstallFunction(runner DarwinInstallRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params := &DarwinInstallParams{}
		if err := lib.UnpackArgs(
			b.Name(), args, kwargs,
			"name", &params.Name,
			"option?", &params.Option,
			"cask?", &params.Cask,
		); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		log.Debug().
			Str("name", params.Name).
			Str("option", params.Option).
			Bool("cask", params.Cask).
			Msg("params")

		ui.Infof("Installing package. Name: %s, Option: %s, Cask: %v\n", params.Name, params.Option, params.Cask)
		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func NewDarwinUninstallFunction(runner DarwinUninstallRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		name := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "name", &name); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}

		ui.Printf("Uninstalling package. Name: %s\n", name)
		if err := runner(ctx, name); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func DarwinInstalled(execute backend.Execute, fs afero.Fs) DarwinInstalledRunner {
	return func(ctx context.Context, name string) bool {
		installed := false
		walkFn := func(path string, info os.FileInfo, err error) error {
			installed = installed || strings.Contains(path, name)
			return nil
		}

		// brew list
		out := new(bytes.Buffer)
		if err := execute(ctx, &backend.ExecuteParams{
			Cmd:    "brew",
			Args:   []string{"--prefix"},
			Stdout: out,
		}); err != nil {
			return false
		}
		prefix := strings.ReplaceAll(out.String(), "\n", "")
		_ = afero.Walk(fs, fmt.Sprintf("%s/Celler", prefix), walkFn)

		// brew cask list
		_ = afero.Walk(fs, "/opt/homebrew-cask/Caskroom", walkFn)
		_ = afero.Walk(fs, "/usr/local/Caskroom", walkFn)

		return installed
	}
}

func DarwinInstall(execute backend.Execute, fs afero.Fs) DarwinInstallRunner {
	return func(ctx context.Context, p *DarwinInstallParams) error {
		if installed := DarwinInstalled(execute, fs)(ctx, p.Name); installed {
			return nil
		}

		args := []string{"install", p.Option, p.Name}
		if p.Cask {
			args = []string{"cask", "install", p.Option, p.Name}
		}

		if err := execute(ctx, &backend.ExecuteParams{Cmd: "brew", Args: args}); err != nil {
			return xerrors.Errorf("Install package failed: %s: %w", p.Name, err)
		}
		return nil
	}
}

func DarwinUninstall(execute backend.Execute, fs afero.Fs) DarwinUninstallRunner {
	return func(ctx context.Context, name string) error {
		if installed := DarwinInstalled(execute, fs)(ctx, name); !installed {
			return nil
		}

		if err := execute(ctx, &backend.ExecuteParams{
			Cmd:  "brew",
			Args: []string{"uninstall", name},
		}); err != nil {
			return xerrors.Errorf("Remove package failed: %w", err)
		}
		return nil
	}
}
