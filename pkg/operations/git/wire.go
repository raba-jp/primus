//+build wireinject

package git

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/operations/git/handlers"
	"github.com/raba-jp/primus/pkg/operations/git/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func Clone() starlark.Fn {
	wire.Build(
		handlers.SetFileSystemStore,
		handlers.NewClone,
		starlarkfn.Clone,
	)
	return nil
}
