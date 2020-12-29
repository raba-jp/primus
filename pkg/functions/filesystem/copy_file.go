package filesystem

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
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

		params, err := parseCopyArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}

		params.Cwd = filepath.Dir(starlark.GetCurrentFilePath(thread))

		log.Ctx(ctx).Debug().
			Str("src", params.Src).
			Str("dest", params.Dest).
			Str("permission", params.Permission.String()).
			Str("cwd", params.Cwd).
			Msg("params")
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
		log.Ctx(ctx).Info().
			Str("source", params.Src).
			Str("destination", params.Dest).
			Str("permission", fmt.Sprintf("%v", params.Permission)).
			Msg("copied file")
		return nil
	}
}
