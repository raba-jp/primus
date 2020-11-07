package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/spf13/afero"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type CopyFileParams struct {
	Src        string
	Dest       string
	Permission os.FileMode
	Cwd        string
}

type CopyFileRunner func(ctx context.Context, p *CopyFileParams) error

func NewCopyFileFunction(runner CopyFileRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "function/copy_file")

		params, err := parseCopyArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		params.Cwd = filepath.Dir(starlark.GetCurrentFilePath(thread))

		logger.Debug(
			"Params",
			zap.String("src", params.Src),
			zap.String("dest", params.Dest),
			zap.String("permission", params.Permission.String()),
			zap.String("cwd", params.Cwd),
		)
		ui.Infof(
			"Coping file. Src: %s, Dest: %s, Permission: %v, Cwd: %s\n",
			params.Src,
			params.Dest,
			params.Permission.String(),
			params.Cwd,
		)
		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseCopyArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*CopyFileParams, error) {
	p := &CopyFileParams{}

	var perm = 0o777
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "src", &p.Src, "dest", &p.Dest, "permission?", &perm); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}

	p.Permission = os.FileMode(perm)
	return p, nil
}

func CopyFile(fs afero.Fs) CopyFileRunner {
	return func(ctx context.Context, params *CopyFileParams) error {
		_, logger := ctxlib.LoggerWithNamespace(ctx, "copy_file")

		if !filepath.IsAbs(params.Src) {
			params.Src = filepath.Join(params.Cwd, params.Src)
		}
		if !filepath.IsAbs(params.Dest) {
			params.Dest = filepath.Join(params.Cwd, params.Dest)
		}

		existsFile := ExistsFile(fs)
		if ret := existsFile(ctx, params.Src); !ret {
			return xerrors.New("Source file not exists")
		}
		if ret := existsFile(ctx, params.Dest); ret {
			return xerrors.New("Destination file already exists")
		}

		srcFile, err := fs.Open(params.Src)
		if err != nil {
			return xerrors.Errorf("Failed to open src file: %w", err)
		}
		destFile, err := fs.OpenFile(params.Dest, os.O_WRONLY|os.O_CREATE, params.Permission)
		if err != nil {
			return xerrors.Errorf("Failed to open dest file: %w", err)
		}
		if _, err := io.Copy(destFile, srcFile); err != nil {
			return xerrors.Errorf("Failed to copy src to dest: %w", err)
		}
		logger.Info(
			"Copied file",
			zap.String("source", params.Src),
			zap.String("destination", params.Dest),
			zap.String("permission", fmt.Sprintf("%v", params.Permission)),
		)
		return nil
	}
}
