package functions

import (
	"github.com/raba-jp/primus/pkg/internal/backend"
	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"
)

func IsDarwin(execIF exec.Interface) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ret := backend.DetectDarwin(execIF)
		if ret {
			return starlark.True, nil
		}
		return starlark.False, nil
	}
}

func IsArchLinux(fs afero.Fs) StarlarkFn {
	return func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		ret := backend.DetectArchLinux(fs)
		if ret {
			return starlark.True, nil
		}
		return starlark.False, nil
	}
}
