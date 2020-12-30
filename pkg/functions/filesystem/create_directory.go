package filesystem

import (
	"context"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"

	"github.com/raba-jp/primus/pkg/cli/ui"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type CreateDirectoryParams struct {
	Path       string
	Permission os.FileMode
	Cwd        string
}

type CreateDirectoryRunner func(ctx context.Context, p *CreateDirectoryParams) error

func NewCreateDirectoryFunction(runner CreateDirectoryRunner) starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ctx := starlark.ToContext(thread)

		params, err := parseCreateArgs(b, args, kwargs)
		if err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		params.Cwd = filepath.Dir(starlark.GetCurrentFilePath(thread))
		log.Debug().
			Str("path", params.Path).
			Str("permission", params.Permission.String()).
			Str("cwd", params.Cwd).
			Msg("params")

		ui.Infof("Creating directories: %s\n", params.Path)
		if err := runner(ctx, params); err != nil {
			return lib.None, xerrors.Errorf(": %w", err)
		}
		return lib.None, nil
	}
}

func parseCreateArgs(b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (*CreateDirectoryParams, error) {
	p := &CreateDirectoryParams{}
	var perm = 0o644
	if err := lib.UnpackArgs(b.Name(), args, kwargs, "path", &p.Path, "permission?", &perm, "cwd", &p.Cwd); err != nil {
		return nil, xerrors.Errorf("Failed to parse arguments: %w", err)
	}
	p.Permission = os.FileMode(perm)
	return p, nil
}

func CreateDirectory(fs afero.Fs) CreateDirectoryRunner {
	return func(ctx context.Context, params *CreateDirectoryParams) error {
		if !filepath.IsAbs(params.Path) {
			params.Path = filepath.Join(params.Cwd, params.Path)
		}

		if err := fs.MkdirAll(params.Path, params.Permission); err != nil {
			return xerrors.Errorf("Create directory fialed: %w", err)
		}
		log.Info().
			Str("path", params.Path).
			Str("permission", params.Permission.String()).
			Str("cwd", params.Cwd).
			Msg("create directory")
		return nil
	}
}
