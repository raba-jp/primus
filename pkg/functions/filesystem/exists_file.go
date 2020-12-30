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

type ExistsFileRunner func(ctx context.Context, path string) (exists bool)

func NewExistsFileFunction(runner ExistsFileRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		path := ""
		if err := lib.UnpackArgs(b.Name(), args, kwargs, "path", &path); err != nil {
			return lib.False, xerrors.Errorf("Failed to parse arguments: %w", err)
		}
		log.Debug().Str("path", path).Msg("Params")

		ui.Infof("Check existence file. Path: %s\n", path)
		return starlark.ToBool(runner(ctx, path)), nil
	}
}

func ExistsFile(fs afero.Fs) ExistsFileRunner {
	return func(ctx context.Context, path string) bool {
		_, err := fs.Stat(path)
		if err == nil {
			log.Info().Msg("already exists file")
			return true
		}
		return false
	}
}
