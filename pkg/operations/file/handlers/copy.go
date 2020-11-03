//go:generate mockery -outpkg=mocks -case=snake -name=CopyHandler

package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/ctxlib"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type CopyParams struct {
	Src        string
	Dest       string
	Permission os.FileMode
	Cwd        string
}

type CopyHandler interface {
	Run(ctx context.Context, p *CopyParams) (err error)
}

type CopyHandlerFunc func(ctx context.Context, p *CopyParams) error

func (f CopyHandlerFunc) Run(ctx context.Context, p *CopyParams) error {
	return f(ctx, p)
}

func NewCopy(fs afero.Fs) CopyHandler {
	return CopyHandlerFunc(func(ctx context.Context, p *CopyParams) error {
		if dryrun := ctxlib.DryRun(ctx); dryrun {
			copyDryRun(p)
			return nil
		}
		return copy(ctx, fs, p)
	})
}

func copy(ctx context.Context, fs afero.Fs, p *CopyParams) error {
	_, logger := ctxlib.LoggerWithNamespace(ctx, "file_copy")
	if !filepath.IsAbs(p.Src) {
		p.Src = filepath.Join(p.Cwd, p.Src)
	}
	if !filepath.IsAbs(p.Dest) {
		p.Dest = filepath.Join(p.Cwd, p.Dest)
	}

	srcFile, err := fs.Open(p.Src)
	if err != nil {
		return xerrors.Errorf("Failed to open src file: %w", err)
	}
	destFile, err := fs.OpenFile(p.Dest, os.O_WRONLY|os.O_CREATE, p.Permission)
	if err != nil {
		return xerrors.Errorf("Failed to open dest file: %w", err)
	}
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return xerrors.Errorf("Failed to copy src to dest: %w", err)
	}
	logger.Info(
		"copied file",
		zap.String("source", p.Src),
		zap.String("destination", p.Dest),
		zap.String("permission", fmt.Sprintf("%v", p.Permission)),
	)
	return nil
}

func copyDryRun(p *CopyParams) {
	ui.Printf("cp %s %s\n", p.Src, p.Dest)
}
