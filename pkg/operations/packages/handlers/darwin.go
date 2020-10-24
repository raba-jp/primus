package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

type DarwinPkgCheckInstallHandler interface {
	CheckInstall(ctx context.Context, name string) (ok bool)
}

type DarwinPkgInstallParams struct {
	Name   string
	Option string
	Cask   bool
	Cmd    string
}

type DarwinPkgInstallHandler interface {
	Install(ctx context.Context, dryrun bool, p *DarwinPkgInstallParams) (err error)
}

type DarwinPkgUninstallParams struct {
	Name string
	Cask bool
	Cmd  string
}

type DarwinPkgUninstallHandler interface {
	Uninstall(ctx context.Context, dryrun bool, p *DarwinPkgUninstallParams) (err error)
}

type darwin struct {
	DarwinPkgCheckInstallHandler
	DarwinPkgInstallHandler
	DarwinPkgUninstallHandler
	Exec exec.Interface
	Fs   afero.Fs
}

func NewDarwinPkgCheckInstallHandler(execIF exec.Interface, fs afero.Fs) DarwinPkgCheckInstallHandler {
	return newDarwin(execIF, fs)
}

func NewDarwinPkgInstallHandler(execIF exec.Interface, fs afero.Fs) DarwinPkgInstallHandler {
	return newDarwin(execIF, fs)
}

func NewDarwinPkgUninstallHandler(execIF exec.Interface, fs afero.Fs) DarwinPkgUninstallHandler {
	return newDarwin(execIF, fs)
}

func newDarwin(execIF exec.Interface, fs afero.Fs) *darwin {
	return &darwin{Exec: execIF, Fs: fs}
}

func (d *darwin) CheckInstall(ctx context.Context, name string) bool {
	installed := false
	walkFn := func(path string, info os.FileInfo, err error) error {
		installed = installed || strings.Contains(path, name)
		return nil
	}

	// brew list
	res, _ := d.Exec.CommandContext(ctx, "brew", "--prefix").Output()
	prefix := strings.ReplaceAll(string(res), "\n", "")
	_ = afero.Walk(d.Fs, fmt.Sprintf("%s/Celler", prefix), walkFn)

	// brew cask list
	_ = afero.Walk(d.Fs, "/opt/homebrew-cask/Caskroom", walkFn)
	_ = afero.Walk(d.Fs, "/usr/local/Caskroom", walkFn)

	return installed
}

func (d *darwin) Install(ctx context.Context, dryrun bool, p *DarwinPkgInstallParams) error {
	if dryrun {
		ui.Printf("brew install %s %s\n", p.Option, p.Name)
		return nil
	}

	if installed := d.CheckInstall(ctx, p.Name); installed {
		return nil
	}

	if err := d.Exec.CommandContext(ctx, "brew", "install", p.Option, p.Name).Run(); err != nil {
		return xerrors.Errorf("Install package failed: %s: %w", p.Name, err)
	}
	return nil
}

func (d *darwin) Uninstall(ctx context.Context, dryrun bool, p *DarwinPkgUninstallParams) error {
	if dryrun {
		ui.Printf("brew uninstall %s\n", p.Name)
		return nil
	}

	if installed := d.CheckInstall(ctx, p.Name); !installed {
		return nil
	}

	if err := d.Exec.CommandContext(ctx, "brew", "uninstall", p.Name).Run(); err != nil {
		return xerrors.Errorf("Remove package failed: %w", err)
	}
	return nil
}
