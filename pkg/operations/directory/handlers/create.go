//go:generate mockgen -destination mock/handler.go . CreateHandler

package handlers

import (
	"context"
	"os"
	"path/filepath"

	"github.com/k0kubun/pp"
	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func init() {
	pp.ColoringEnabled = false
}

type CreateParams struct {
	Path       string
	Permission os.FileMode
	Cwd        string
}

func (p *CreateParams) String() string {
	return pp.Sprint(p)
}

type CreateHandler interface {
	Create(ctx context.Context, dryrun bool, p *CreateParams) error
}

type CreateHandlerFunc func(ctx context.Context, dryrun bool, p *CreateParams) error

func (f CreateHandlerFunc) Create(ctx context.Context, dryrun bool, p *CreateParams) error {
	return f(ctx, dryrun, p)
}

func New(fs afero.Fs) CreateHandler {
	return CreateHandlerFunc(func(ctx context.Context, dryrun bool, p *CreateParams) error {
		if dryrun {
			ui.Printf("mkdir -p %s\n", p.Path)
			ui.Printf("chmod %o %s\n", p.Permission, p.Path)
			return nil
		}

		if !filepath.IsAbs(p.Path) {
			p.Path = filepath.Join(p.Cwd, p.Path)
		}

		if err := fs.MkdirAll(p.Path, p.Permission); err != nil {
			return xerrors.Errorf("Create directory fialed: %w", err)
		}
		zap.L().Info("create directory", zap.String("path", p.Path), zap.String("permission", p.Permission.String()))
		return nil
	})
}
