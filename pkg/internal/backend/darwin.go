package backend

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/raba-jp/primus/pkg/internal/handlers"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

var _ Backend = (*DarwinBackend)(nil)

type DarwinBackend struct {
	Backend
	*BaseBackend
	Exec exec.Interface
	Fs   afero.Fs
}

func (b *DarwinBackend) CheckInstall(ctx context.Context, name string) bool {
	installed := false
	walkFn := func(path string, info os.FileInfo, err error) error {
		installed = installed || strings.Contains(path, name)
		return nil
	}

	// brew list
	res, _ := b.Exec.CommandContext(ctx, "brew", "--prefix").Output()
	prefix := strings.ReplaceAll(string(res), "\n", "")
	_ = afero.Walk(b.Fs, fmt.Sprintf("%s/Celler", prefix), walkFn)

	// brew cask list
	_ = afero.Walk(b.Fs, "/opt/homebrew-cask/Caskroom", walkFn)
	_ = afero.Walk(b.Fs, "/usr/local/Caskroom", walkFn)

	return installed
}

func (b *DarwinBackend) Install(ctx context.Context, dryrun bool, p *handlers.InstallParams) error {
	if dryrun {
		ui.Printf("brew install %s %s", p.Option, p.Name)
		return nil
	}

	if err := b.Exec.CommandContext(ctx, "brew", "install", p.Option, p.Name).Run(); err != nil {
		return xerrors.Errorf("Install package failed: %s: %w", p.Name, err)
	}
	return nil
}

func (b *DarwinBackend) Uninstall(ctx context.Context, dryrun bool, p *handlers.UninstallParams) error {
	if dryrun {
		ui.Printf("brew uninstall %s", p.Name)
		return nil
	}

	if err := b.Exec.CommandContext(ctx, "brew", "uninstall", p.Name).Run(); err != nil {
		return xerrors.Errorf("Remove package failed: %w", err)
	}
	return nil
}

func (b *DarwinBackend) FileCopy(ctx context.Context, dryrun bool, p *handlers.FileCopyParams) error {
	return b.BaseBackend.FileCopy(ctx, dryrun, p)
}

func (b *DarwinBackend) FileMove(ctx context.Context, dryrun bool, p *handlers.FileMoveParams) error {
	return b.BaseBackend.FileMove(ctx, dryrun, p)
}

func (b *DarwinBackend) Symlink(ctx context.Context, dryrun bool, p *handlers.SymlinkParams) error {
	return b.BaseBackend.Symlink(ctx, dryrun, p)
}

func (b *DarwinBackend) HTTPRequest(ctx context.Context, dryrun bool, p *handlers.HTTPRequestParams) error {
	return b.BaseBackend.HTTPRequest(ctx, dryrun, p)
}

func (b *DarwinBackend) Command(ctx context.Context, dryrun bool, p *handlers.CommandParams) error {
	return b.BaseBackend.Command(ctx, dryrun, p)
}

func (b *DarwinBackend) FishSetVariable(ctx context.Context, dryrun bool, p *handlers.FishSetVariableParams) error {
	return b.BaseBackend.FishSetVariable(ctx, dryrun, p)
}
