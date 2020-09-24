package handlers

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"
)

type archLinux struct {
	CheckInstallHandler
	InstallHandler
	UninstallHandler
	Exec exec.Interface
}

func (b *archLinux) CheckInstall(ctx context.Context, name string) bool {
	err := b.Exec.CommandContext(ctx, "pacman", "-Qg", name).Run()
	return err == nil
}

func (b *archLinux) Install(ctx context.Context, dryrun bool, p *InstallParams) error {
	if dryrun {
		ui.Printf("pacman -S --noconfirm %s %s\n", p.Option, p.Name)
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	if err := b.Exec.CommandContext(ctx, "pacman", "-S", "--noconfirm", p.Option, p.Name).Run(); err != nil {
		return xerrors.Errorf("Install package failed: %s: %w", p.Name, err)
	}
	return nil
}

func (b *archLinux) Uninstall(ctx context.Context, dryrun bool, p *UninstallParams) error {
	if dryrun {
		ui.Printf("pacman -R --noconfirm %s\n", p.Name)
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	if err := b.Exec.CommandContext(ctx, "pacman", "-R", "--noconfirm", p.Name).Run(); err != nil {
		return xerrors.Errorf("Remove package failed: %s: %w", p.Name, err)
	}
	return nil
}
