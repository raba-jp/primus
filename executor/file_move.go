package executor

import (
	"context"

	"golang.org/x/xerrors"
)

type FileMoveParams struct {
	Src  string
	Dest string
}

func (e *executor) FileMove(ctx context.Context, p *FileMoveParams) (bool, error) {
	if err := e.Fs.Rename(p.Src, p.Dest); err != nil {
		return false, xerrors.Errorf("Failed to move file: %s => %s: %w", p.Src, p.Dest, err)
	}
	return true, nil
}
