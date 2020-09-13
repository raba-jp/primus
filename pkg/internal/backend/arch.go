package backend

import (
	"context"

	"github.com/raba-jp/primus/pkg/internal/exec"
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

func (b *ArchLinuxBackend) Install(ctx context.Context, p *InstallParams) error {
	ctx, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	if err := b.Exec.CommandContext(ctx, "pacman", "-S", "--noconfirm", p.Option, p.Name).Run(); err != nil {
		return xerrors.Errorf("Install package failed: %s: %w", p.Name, err)
	}
	return nil
}

func (b *ArchLinuxBackend) Uninstall(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	if err := b.Exec.CommandContext(ctx, "pacman", "-R", "--noconfirm", name).Run(); err != nil {
		return xerrors.Errorf("Remove package failed: %s: %w", name, err)
	}
	return nil
}

func (b *ArchLinuxBackend) FileCopy(ctx context.Context, p *FileCopyParams) error {
	return b.BaseBackend.FileCopy(ctx, p)
}

func (b *ArchLinuxBackend) FileMove(ctx context.Context, p *FileMoveParams) error {
	return b.BaseBackend.FileMove(ctx, p)
}

func (b *ArchLinuxBackend) Symlink(ctx context.Context, p *SymlinkParams) error {
	return b.BaseBackend.Symlink(ctx, p)
}

func (b *ArchLinuxBackend) HTTPRequest(ctx context.Context, p *HTTPRequestParams) error {
	return b.BaseBackend.HTTPRequest(ctx, p)
}

func (b *ArchLinuxBackend) Command(ctx context.Context, p *CommandParams) error {
	return b.BaseBackend.Command(ctx, p)
}
