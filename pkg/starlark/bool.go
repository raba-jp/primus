package starlark

import (
	lib "go.starlark.net/starlark"
)

func ToBool(v bool) lib.Bool {
	if v {
		return lib.True
	}
	return lib.False
}
