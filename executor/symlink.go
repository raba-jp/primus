package executor

import (
	"context"

	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type SymlinkParams struct {
	Src  string
	Dest string
	User string
}

func (e *executor) Symlink(ctx context.Context, p *SymlinkParams) (bool, error) {
	if ext := e.fileExists(p.Dest); ext {
		return false, xerrors.New("File already exists")
	}

	linker, ok := e.Fs.(afero.Symlinker)
	if !ok {
		return false, xerrors.New("This filesystem does not support symlink")
	}
	if err := linker.SymlinkIfPossible(p.Src, p.Dest); err != nil {
		return false, xerrors.Errorf("Failed to create symbolic link: %w", err)
	}

	return true, nil
}

func (e *executor) fileExists(path string) bool {
	_, err := e.Fs.Stat(path)
	if err == nil {
		zap.L().Info("Already exists file")
		return true
	}
	return false
}
