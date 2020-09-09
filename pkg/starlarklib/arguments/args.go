package arguments

import "go.starlark.net/starlark"

type Arguments interface {
	Parse(b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) error
}
