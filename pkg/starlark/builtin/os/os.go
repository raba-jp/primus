package os

import (
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/starlark/builtin"
	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
)

func IsDarwin(execIF exec.Interface) builtin.StarlarkFn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ret := backend.DetectDarwin(execIF)
		if ret {
			return lib.True, nil
		}
		return lib.False, nil
	}
}

func IsArchLinux(fs afero.Fs) builtin.StarlarkFn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		ret := backend.DetectArchLinux(fs)
		if ret {
			return lib.True, nil
		}
		return lib.False, nil
	}
}
