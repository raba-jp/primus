//go:generate mockery -outpkg=mocks -case=snake -name=CreateHandler

package handlers

import (
	"context"
	"os"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type CreateParams struct {
	Path       string
	Permission os.FileMode
	Cwd        string
}

type CreateHandler interface {
	Run(ctx context.Context, p *CreateParams) (err error)
}

type CreateHandlerFunc func(ctx context.Context, p *CreateParams) error

func (f CreateHandlerFunc) Run(ctx context.Context, p *CreateParams) error {
	return f(ctx, p)
}

func New(fs afero.Fs) CreateHandler {
	return CreateHandlerFunc(func(ctx context.Context, p *CreateParams) error {
		if dryrun := ctxlib.DryRun(ctx); dryrun {
			createDryRun(p)
			return nil
		}
		return create(ctx, fs, p)
	})
}

func create(ctx context.Context, fs afero.Fs, p *CreateParams) error {
	_, logger := ctxlib.LoggerWithNamespace(ctx, "create_directory")
	if !filepath.IsAbs(p.Path) {
		p.Path = filepath.Join(p.Cwd, p.Path)
	}

	if err := fs.MkdirAll(p.Path, p.Permission); err != nil {
		return xerrors.Errorf("Create directory fialed: %w", err)
	}
	logger.Info("create directory", zap.String("path", p.Path), zap.String("permission", p.Permission.String()))
	return nil
}

func createDryRun(p *CreateParams) {
	ui.Printf("mkdir -p %s\n", p.Path)
	ui.Printf("chmod %o %s\n", p.Permission, p.Path)
}
