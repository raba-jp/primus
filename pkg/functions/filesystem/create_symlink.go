package filesystem

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type CreateSymlinkParams struct {
	Src  string
	Dest string
}

type CreateSymlinkRunner func(ctx context.Context, params *CreateSymlinkParams) error

func NewCreateSymlinkFunction(runner CreateSymlinkRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params := &CreateSymlinkParams{}
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "src", &params.Src, "dest", &params.Dest); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		log.Debug().
			Str("src", params.Src).
			Str("dest", params.Dest).
			Msg("params")

		ui.Infof("Creating symbolic link. Src: %s, Dest: %s\n", params.Src, params.Dest)
		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func CreateSymlink(fs afero.Fs) CreateSymlinkRunner {
	return func(ctx context.Context, params *CreateSymlinkParams) error {
		existFile := ExistsFile(fs)
		if ret := existFile(ctx, params.Src); !ret {
			return xerrors.New("Source file not found")
		}
		if existFile(ctx, params.Dest) {
			return xerrors.New("Destination file already exists")
		}

		linker, ok := fs.(afero.Symlinker)
		if !ok {
			return xerrors.New("This filesystem does not support symlink")
		}
		if err := linker.SymlinkIfPossible(params.Src, params.Dest); err != nil {
			return xerrors.Errorf("Failed to create symbolic link: %w", err)
		}

		log.Info().
			Str("source", params.Src).
			Str("destination", params.Dest).
			Msg("Create symbolic link")

		return nil
	}
}
