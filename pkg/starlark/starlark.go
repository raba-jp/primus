package starlark

import (
	lib "go.starlark.net/starlark"
)

type Fn = func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kargs []lib.Tuple) (lib.Value, error)

func ExecForTest(name string, data string, fn Fn) (lib.StringDict, error) {
	predeclared := lib.StringDict{
		name: lib.NewBuiltin(name, fn),
	}
	return lib.ExecFile(NewThread("test"), "/sym/test.star", data, predeclared)
}
