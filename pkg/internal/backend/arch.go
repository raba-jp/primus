package backend

import (
	"context"

	"github.com/raba-jp/primus/pkg/internal/exec"
	"golang.org/x/xerrors"
)

var _ Backend = (*ArchLinuxBackend)(nil)

type ArchLinuxBackend struct {
	Backend
	Exec exec.Interface
}

func (p *ArchLinuxBackend) CheckInstall(ctx context.Context, name string) bool {
	err := p.Exec.CommandContext(ctx, "pacman", "-Qg", name).Run()
	return err == nil
}

func (p *ArchLinuxBackend) Install(ctx context.Context, name string, option string) error {
	ctx, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	if err := p.Exec.CommandContext(ctx, "pacman", "-S", "--noconfirm", option, name).Run(); err != nil {
		return xerrors.Errorf("Install package failed: %s: %w", name, err)
	}
	return nil
}

func (p *ArchLinuxBackend) Uninstall(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, installTimeout)
	defer cancel()
	if err := p.Exec.CommandContext(ctx, "pacman", "-R", "--noconfirm", name).Run(); err != nil {
		return xerrors.Errorf("Remove package failed: %s: %w", name, err)
	}
	return nil
}
