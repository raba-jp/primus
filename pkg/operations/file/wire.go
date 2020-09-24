//+build wireinject

package file

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/file/handlers"
	"github.com/raba-jp/primus/pkg/operations/file/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func Copy() starlark.Fn {
	wire.Build(
		backend.NewFs,
		handlers.NewCopy,
		starlarkfn.Copy,
	)
	return nil
}

func Move() starlark.Fn {
	wire.Build(
		backend.NewFs,
		handlers.NewMove,
		starlarkfn.Move,
	)
	return nil
}

func Symlink() starlark.Fn {
	wire.Build(
		backend.NewFs,
		handlers.NewSymlink,
		starlarkfn.Symlink,
	)
	return nil
}
