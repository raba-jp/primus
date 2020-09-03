package apply

import (
	"context"

	"github.com/raba-jp/primus/executor"
	"golang.org/x/xerrors"
)

func (e *applyExecutor) FileMove(ctx context.Context, p *executor.FileMoveParams) (bool, error) {
	if err := e.Fs.Rename(p.Src, p.Dest); err != nil {
		return false, xerrors.Errorf("Failed to move file: %s => %s: %w", p.Src, p.Dest, err)
	}
	return true, nil
}
