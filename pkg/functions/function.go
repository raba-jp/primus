package functions

import (
	"go.starlark.net/starlark"
)

type StarlarkFn = func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error)

func toStarlarkBool(v bool) starlark.Value {
	if v {
		return starlark.True
	}
	return starlark.False
}
