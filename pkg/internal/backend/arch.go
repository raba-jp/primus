package backend

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/raba-jp/primus/pkg/internal/handlers"
	"golang.org/x/xerrors"
)

var _ Backend = (*ArchLinuxBackend)(nil)

type ArchLinuxBackend struct {
	Backend
	*BaseBackend
	Exec exec.Interface
}

func (b *ArchLinuxBackend) CheckInstall(ctx context.Context, name string) bool {
	err := b.Exec.CommandContext(ctx, "pacman", "-Qg", name).Run()
	return err == nil
}

func (b *ArchLinuxBackend) Install(ctx context.Context, dryrun bool, p *handlers.InstallParams) error {
	if dryrun {
		ui.Printf("pacman -S --noconfirm %s %s", p.Option, p.Name)
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	if err := b.Exec.CommandContext(ctx, "pacman", "-S", "--noconfirm", p.Option, p.Name).Run(); err != nil {
		return xerrors.Errorf("Install package failed: %s: %w", p.Name, err)
	}
	return nil
}

func (b *ArchLinuxBackend) Uninstall(ctx context.Context, dryrun bool, p *handlers.UninstallParams) error {
	if dryrun {
		ui.Printf("pacman -R %s", p.Name)
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	if err := b.Exec.CommandContext(ctx, "pacman", "-R", "--noconfirm", p.Name).Run(); err != nil {
		return xerrors.Errorf("Remove package failed: %s: %w", p.Name, err)
	}
	return nil
}

func (b *ArchLinuxBackend) FileCopy(ctx context.Context, dryrun bool, p *handlers.FileCopyParams) error {
	return b.BaseBackend.FileCopy(ctx, dryrun, p)
}

func (b *ArchLinuxBackend) FileMove(ctx context.Context, dryrun bool, p *handlers.FileMoveParams) error {
	return b.BaseBackend.FileMove(ctx, dryrun, p)
}

func (b *ArchLinuxBackend) Symlink(ctx context.Context, dryrun bool, p *handlers.SymlinkParams) error {
	return b.BaseBackend.Symlink(ctx, dryrun, p)
}

func (b *ArchLinuxBackend) HTTPRequest(ctx context.Context, dryrun bool, p *handlers.HTTPRequestParams) error {
	return b.BaseBackend.HTTPRequest(ctx, dryrun, p)
}

func (b *ArchLinuxBackend) Command(ctx context.Context, dryrun bool, p *handlers.CommandParams) error {
	return b.BaseBackend.Command(ctx, dryrun, p)
}

func (b *ArchLinuxBackend) FishSetVariable(ctx context.Context, dryrun bool, p *handlers.FishSetVariableParams) error {
	return b.BaseBackend.FishSetVariable(ctx, dryrun, p)
}
