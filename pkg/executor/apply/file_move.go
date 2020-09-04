package apply

import (
	"context"

	"github.com/raba-jp/primus/pkg/executor"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func (e *applyExecutor) FileMove(ctx context.Context, p *executor.FileMoveParams) (bool, error) {
	if err := e.Fs.Rename(p.Src, p.Dest); err != nil {
		return false, xerrors.Errorf("Failed to move file: %s => %s: %w", p.Src, p.Dest, err)
	}
	zap.L().Info(
		"moved file",
		zap.String("source", p.Src),
		zap.String("destination", p.Dest),
	)
	return true, nil
}
