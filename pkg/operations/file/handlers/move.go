package handlers

import (
	"context"
	"path/filepath"

	"github.com/k0kubun/pp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type MoveParams struct {
	Src  string
	Dest string
	Cwd  string
}

func (p *MoveParams) String() string {
	return pp.Sprintf("%v\n", p)
}

type MoveHandler interface {
	Move(ctx context.Context, dryrun bool, p *MoveParams) (err error)
}

type MoveHandlerFunc func(ctx context.Context, dryrun bool, p *MoveParams) error

func (f MoveHandlerFunc) Move(ctx context.Context, dryrun bool, p *MoveParams) error {
	return f(ctx, dryrun, p)
}

func NewMove(fs afero.Fs) MoveHandler {
	return MoveHandlerFunc(func(ctx context.Context, dryrun bool, p *MoveParams) error {
		if dryrun {
			ui.Printf("mv %s %s\n", p.Src, p.Dest)
			return nil
		}

		if !filepath.IsAbs(p.Src) {
			p.Src = filepath.Join(p.Cwd, p.Src)
		}
		if !filepath.IsAbs(p.Dest) {
			p.Dest = filepath.Join(p.Cwd, p.Dest)
		}

		if err := fs.Rename(p.Src, p.Dest); err != nil {
			return xerrors.Errorf("Failed to move file: %s => %s: %w", p.Src, p.Dest, err)
		}
		zap.L().Info(
			"moved file",
			zap.String("source", p.Src),
			zap.String("destination", p.Dest),
		)
		return nil
	})
}
