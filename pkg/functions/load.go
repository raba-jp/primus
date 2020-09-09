package functions

import (
	"path/filepath"

	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type StarlarkLoadFn = func(thread *starlark.Thread, module string) (starlark.StringDict, error)

func Load(fs afero.Fs) StarlarkLoadFn {
	return func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		path := starlarklib.GetCurrentFilePath(thread)
		mp := filepath.Join(filepath.Dir(path), module)
		data, err := afero.ReadFile(fs, mp)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		ctx := starlarklib.GetCtx(thread)
		childThread := starlarklib.NewThread(module, starlarklib.WithLoad(Load(fs)), starlarklib.WithContext(ctx))

		return starlark.ExecFile(childThread, mp, data, nil)
	}
}
