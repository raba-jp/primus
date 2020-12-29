package filesystem

import (
	"context"
	"path/filepath"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type MoveFileParams struct {
	Src  string
	Dest string
	Cwd  string
}

type MoveFileRunner func(ctx context.Context, params *MoveFileParams) error

func NewMoveFileFunction(runner MoveFileRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params := &MoveFileParams{}
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "src", &params.Src, "dest", &params.Dest); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		params.Cwd = filepath.Dir(starlark.GetCurrentFilePath(thread))

		log.Ctx(ctx).Debug().
			Str("src", params.Src).
			Str("dest", params.Dest).
			Msg("params")

		ui.Infof("Moving file. Source: %s, Destination: %s\n", params.Src, params.Dest)
		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func MoveFile(fs afero.Fs) MoveFileRunner {
	return func(ctx context.Context, params *MoveFileParams) error {
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
		if existsFile(ctx, params.Dest) {
			return xerrors.New("Destination file already exists")
		}

		if err := fs.Rename(params.Src, params.Dest); err != nil {
			return xerrors.Errorf("Failed to move file: %s => %s: %w", params.Src, params.Dest, err)
		}

		log.Ctx(ctx).Info().
			Str("src", params.Src).
			Str("dest", params.Dest).
			Msg("moved file")
		return nil
	}
}
