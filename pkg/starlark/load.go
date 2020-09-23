package starlark

import (
	"path/filepath"

	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type StarlarkLoadFn = func(thread *lib.Thread, module string) (lib.StringDict, error)

func Load(fs afero.Fs, predeclared lib.StringDict) StarlarkLoadFn {
	return func(thread *lib.Thread, module string) (lib.StringDict, error) {
		var modulePath string
		if filepath.IsAbs(module) {
			modulePath = module
		} else {
			path := GetCurrentFilePath(thread)
			modulePath = filepath.Join(filepath.Dir(path), module)
		}

		data, err := afero.ReadFile(fs, modulePath)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		child := NewThread(module, withTakeOverParent(thread))

		return lib.ExecFile(child, modulePath, data, predeclared)
	}
}
