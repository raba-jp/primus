package backend

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/spf13/afero"
	"golang.org/x/xerrors"
)

var _ Backend = (*DarwinBackend)(nil)

type DarwinBackend struct {
	Backend
	Exec exec.Interface
	Fs   afero.Fs
}

func (p *DarwinBackend) CheckInstall(ctx context.Context, name string) bool {
	installed := false
	walkFn := func(path string, info os.FileInfo, err error) error {
		installed = installed || strings.Contains(path, name)
		return nil
	}

	// brew list
	res, _ := p.Exec.CommandContext(ctx, "brew", "--prefix").Output()
	prefix := strings.ReplaceAll(string(res), "\n", "")
	_ = afero.Walk(p.Fs, fmt.Sprintf("%s/Celler", prefix), walkFn)

	// brew cask list
	_ = afero.Walk(p.Fs, "/opt/homebrew-cask/Caskroom", walkFn)
	_ = afero.Walk(p.Fs, "/usr/local/Caskroom", walkFn)

	return installed
}

func (p *DarwinBackend) Install(ctx context.Context, name string, option string) error {
	if err := p.Exec.CommandContext(ctx, "brew", "install", option, name).Run(); err != nil {
		return xerrors.Errorf("Install package failed: %s: %w", name, err)
	}
	return nil
}

func (p *DarwinBackend) Uninstall(ctx context.Context, name string) error {
	if err := p.Exec.CommandContext(ctx, "brew", "uninstall", name).Run(); err != nil {
		return xerrors.Errorf("Remove package failed: %w", err)
	}
	return nil
}
