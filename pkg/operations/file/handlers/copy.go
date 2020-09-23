//go:generate mockgen -destination mock/copy.go . CopyHandler

package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/k0kubun/pp"
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

func (p *CopyParams) String() string {
	return pp.Sprintf("%v\n", p)
}

type CopyHandler interface {
	Copy(ctx context.Context, dryrun bool, p *CopyParams) error
}

type CopyHandlerFunc func(ctx context.Context, dryrun bool, p *CopyParams) error

func (f CopyHandlerFunc) Copy(ctx context.Context, dryrun bool, p *CopyParams) error {
	return f(ctx, dryrun, p)
}

func NewCopy(fs afero.Fs) CopyHandler {
	return CopyHandlerFunc(func(ctx context.Context, dryrun bool, p *CopyParams) error {
		if dryrun {
			ui.Printf("cp %s %s\n", p.Src, p.Dest)
			return nil
		}

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
		zap.L().Info(
			"copied file",
			zap.String("source", p.Src),
			zap.String("destination", p.Dest),
			zap.String("permission", fmt.Sprintf("%v", p.Permission)),
		)
		return nil
	})
}
