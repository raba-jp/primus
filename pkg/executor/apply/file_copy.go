package apply

import (
	"context"
	"io"
	"os"

	"github.com/raba-jp/primus/pkg/executor"
	"golang.org/x/xerrors"
)

func (e *applyExecutor) FileCopy(ctx context.Context, p *executor.FileCopyParams) (bool, error) {
	srcFile, err := e.Fs.Open(p.Src)
	if err != nil {
		return false, xerrors.Errorf("Failed to open src file: %w", err)
	}
	// permission
	destFile, err := e.Fs.OpenFile(p.Dest, os.O_WRONLY|os.O_CREATE, p.Permission)
	if err != nil {
		return false, xerrors.Errorf("Failed to open dest file: %w", err)
	}
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return false, xerrors.Errorf("Failed to copy src to dest: %w", err)
	}
	return true, nil
}
