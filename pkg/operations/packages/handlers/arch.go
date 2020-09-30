//go:generate mockgen -destination mock/arch.go . ArchPkgCheckInstallHandler,ArchPkgInstallHandler,ArchPkgUninstallHandler

package handlers

import (
	"context"

	"github.com/k0kubun/pp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"
)

type ArchPkgCheckInstallHandler interface {
	CheckInstall(ctx context.Context, name string) bool
}

type ArchPkgInstallParams struct {
	Name   string
	Option string
	Cmd    string
}

func (p *ArchPkgInstallParams) String() string {
	return pp.Sprintf("%v\n", p)
}

type ArchPkgInstallHandler interface {
	Install(ctx context.Context, dryrun bool, p *ArchPkgInstallParams) error
}

type ArchPkgUninstallParams struct {
	Name string
	Cmd  string
}

func (p *ArchPkgUninstallParams) String() string {
	return pp.Sprintf("%v\n", p)
}

type ArchPkgUninstallHandler interface {
	Uninstall(ctx context.Context, dryrun bool, p *ArchPkgUninstallParams) error
}

type archLinux struct {
	ArchPkgCheckInstallHandler
	ArchPkgInstallHandler
	ArchPkgUninstallHandler
	Exec exec.Interface
}

func NewArchPkgCheckInstallHandler(execIF exec.Interface) ArchPkgCheckInstallHandler {
	return newArchLinux(execIF)
}

func NewArchPkgInstallHandler(execIF exec.Interface) ArchPkgInstallHandler {
	return newArchLinux(execIF)
}

func NewArchPkgUninstallHandler(execIF exec.Interface) ArchPkgUninstallHandler {
	return newArchLinux(execIF)
}

func newArchLinux(execIF exec.Interface) *archLinux {
	return &archLinux{Exec: execIF}
}

func (a *archLinux) CheckInstall(ctx context.Context, name string) bool {
	err := a.Exec.CommandContext(ctx, "pacman", "-Qg", name).Run()
	return err == nil
}

func (a *archLinux) Install(ctx context.Context, dryrun bool, p *ArchPkgInstallParams) error {
	if dryrun {
		ui.Printf("pacman -S --noconfirm %s %s\n", p.Option, p.Name)
		return nil
	}

	if installed := a.CheckInstall(ctx, p.Name); installed {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	if err := a.Exec.CommandContext(ctx, "pacman", "-S", "--noconfirm", p.Option, p.Name).Run(); err != nil {
		return xerrors.Errorf("Install package failed: %s: %w", p.Name, err)
	}
	return nil
}

func (a *archLinux) Uninstall(ctx context.Context, dryrun bool, p *ArchPkgUninstallParams) error {
	if dryrun {
		ui.Printf("pacman -R --noconfirm %s\n", p.Name)
		return nil
	}

	if installed := a.CheckInstall(ctx, p.Name); !installed {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	if err := a.Exec.CommandContext(ctx, "pacman", "-R", "--noconfirm", p.Name).Run(); err != nil {
		return xerrors.Errorf("Remove package failed: %s: %w", p.Name, err)
	}
	return nil
}
