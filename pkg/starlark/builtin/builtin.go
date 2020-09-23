package builtin

import (
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
)

type StarlarkFn = func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kargs []lib.Tuple) (lib.Value, error)

func ExecForTest(name string, data string, fn StarlarkFn) (lib.StringDict, error) {
	predeclared := lib.StringDict{
		name: lib.NewBuiltin(name, fn),
	}
	return lib.ExecFile(starlark.NewThread("test"), "test.star", data, predeclared)
}
