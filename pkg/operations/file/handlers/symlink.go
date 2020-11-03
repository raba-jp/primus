//go:generate mockery -outpkg=mocks -case=snake -name=SymlinkHandler

package handlers

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type SymlinkParams struct {
	Src  string
	Dest string
	User string
}

type SymlinkHandler interface {
	Run(ctx context.Context, p *SymlinkParams) (err error)
}

type SymlinkHandlerFunc func(ctx context.Context, p *SymlinkParams) error

func (f SymlinkHandlerFunc) Run(ctx context.Context, p *SymlinkParams) error {
	return f(ctx, p)
}

func NewSymlink(fs afero.Fs) SymlinkHandler {
	return SymlinkHandlerFunc(func(ctx context.Context, p *SymlinkParams) error {
		if dryrun := ctxlib.DryRun(ctx); dryrun {
			symlinkDryRun(p)
			return nil
		}
		return symlink(ctx, fs, p)
	})
}

func symlink(ctx context.Context, fs afero.Fs, p *SymlinkParams) error {
	_, logger := ctxlib.LoggerWithNamespace(ctx, "symlink")
	if ext := fileExists(fs, p.Dest); ext {
		return xerrors.New("File already exists")
	}

	linker, ok := fs.(afero.Symlinker)
	if !ok {
		return xerrors.New("This filesystem does not support symlink")
	}
	if err := linker.SymlinkIfPossible(p.Src, p.Dest); err != nil {
		return xerrors.Errorf("Failed to create symbolic link: %w", err)
	}

	logger.Info(
		"create symbolic link",
		zap.String("source", p.Src),
		zap.String("destination", p.Dest),
	)

	return nil
}

func symlinkDryRun(p *SymlinkParams) {
	ui.Printf("ln -s %s %s\n", p.Src, p.Dest)
}

func fileExists(fs afero.Fs, path string) bool {
	_, err := fs.Stat(path)
	if err == nil {
		zap.L().Info("Already exists file")
		return true
	}
	return false
}
