//go:generate mockgen -destination mock/symlink.go . SymlinkHandler

package handlers

import (
	"context"

	"github.com/k0kubun/pp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type SymlinkParams struct {
	Src  string
	Dest string
	User string
}

func (p *SymlinkParams) String() string {
	return pp.Sprintf("%v\n", p)
}

type SymlinkHandler interface {
	Symlink(ctx context.Context, dryrun bool, p *SymlinkParams) error
}

type SymlinkHandlerFunc func(ctx context.Context, dryrun bool, p *SymlinkParams) error

func (f SymlinkHandlerFunc) Symlink(ctx context.Context, dryrun bool, p *SymlinkParams) error {
	return f(ctx, dryrun, p)
}

func NewSymlink(fs afero.Fs) SymlinkHandler {
	return SymlinkHandlerFunc(func(ctx context.Context, dryrun bool, p *SymlinkParams) error {
		if dryrun {
			ui.Printf("ln -s %s %s\n", p.Src, p.Dest)
			return nil
		}

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

		zap.L().Info(
			"create symbolic link",
			zap.String("source", p.Src),
			zap.String("destination", p.Dest),
		)

		return nil
	})
}

func fileExists(fs afero.Fs, path string) bool {
	_, err := fs.Stat(path)
	if err == nil {
		zap.L().Info("Already exists file")
		return true
	}
	return false
}
