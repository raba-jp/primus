//+build wireinject

package directory

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/directory/handlers"
	"github.com/raba-jp/primus/pkg/operations/directory/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func Create() starlark.Fn {
	wire.Build(
		backend.NewFs,
		handlers.New,
		starlarkfn.Create,
	)
	return nil
}
