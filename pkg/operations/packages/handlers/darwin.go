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

type darwin struct {
	CheckInstallHandler
	InstallHandler
	UninstallHandler
	Fs   afero.Fs
	Exec exec.Interface
}

func (b *darwin) CheckInstall(ctx context.Context, name string) bool {
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

func (b *darwin) Install(ctx context.Context, dryrun bool, p *InstallParams) error {
	if dryrun {
		ui.Printf("brew install %s %s\n", p.Option, p.Name)
		return nil
	}

	if err := b.Exec.CommandContext(ctx, "brew", "install", p.Option, p.Name).Run(); err != nil {
		return xerrors.Errorf("Install package failed: %s: %w", p.Name, err)
	}
	return nil
}

func (b *darwin) Uninstall(ctx context.Context, dryrun bool, p *UninstallParams) error {
	if dryrun {
		ui.Printf("brew uninstall %s\n", p.Name)
		return nil
	}

	if err := b.Exec.CommandContext(ctx, "brew", "uninstall", p.Name).Run(); err != nil {
		return xerrors.Errorf("Remove package failed: %w", err)
	}
	return nil
}
