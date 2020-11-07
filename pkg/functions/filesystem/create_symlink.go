package filesystem

import (
	"context"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/ctxlib"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
	"go.uber.org/zap"
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
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "function/create_symlink")

		params := &CreateSymlinkParams{}
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "src", &params.Src, "dest", &params.Dest); err != nil {
			return lib.None, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		logger.Debug("Params", zap.String("src", params.Src), zap.String("dest", params.Dest))

		ui.Infof("Creating symbolic link. Src: %s, Dest: %s\n", params.Src, params.Dest)
		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func CreateSymlink(fs afero.Fs) CreateSymlinkRunner {
	return func(ctx context.Context, params *CreateSymlinkParams) error {
		ctx, logger := ctxlib.LoggerWithNamespace(ctx, "create_symlink")

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

		logger.Info(
			"Create symbolic link",
			zap.String("source", params.Src),
			zap.String("destination", params.Dest),
		)

		return nil
	}
}
