//go:generate mockery -outpkg=mocks -case=snake -name=MoveHandler

package handlers

import (
	"context"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type MoveParams struct {
	Src  string
	Dest string
	Cwd  string
}

type MoveHandler interface {
	Run(ctx context.Context, p *MoveParams) (err error)
}

type MoveHandlerFunc func(ctx context.Context, p *MoveParams) error

func (f MoveHandlerFunc) Run(ctx context.Context, p *MoveParams) error {
	return f(ctx, p)
}

func NewMove(fs afero.Fs) MoveHandler {
	return MoveHandlerFunc(func(ctx context.Context, p *MoveParams) error {
		if dryrun := ctxlib.DryRun(ctx); dryrun {
			moveDryRun(p)
			return nil
		}
		return move(ctx, fs, p)
	})
}

func move(ctx context.Context, fs afero.Fs, p *MoveParams) error {
	_, logger := ctxlib.LoggerWithNamespace(ctx, "file_move")
	if !filepath.IsAbs(p.Src) {
		p.Src = filepath.Join(p.Cwd, p.Src)
	}
	if !filepath.IsAbs(p.Dest) {
		p.Dest = filepath.Join(p.Cwd, p.Dest)
	}

	if err := fs.Rename(p.Src, p.Dest); err != nil {
		return xerrors.Errorf("Failed to move file: %s => %s: %w", p.Src, p.Dest, err)
	}
	logger.Info(
		"moved file",
		zap.String("source", p.Src),
		zap.String("destination", p.Dest),
	)
	return nil
}

func moveDryRun(p *MoveParams) {
	ui.Printf("mv %s %s\n", p.Src, p.Dest)
}
