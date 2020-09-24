//+build wireinject

package fish

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/fish/handlers"
	"github.com/raba-jp/primus/pkg/operations/fish/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func SetPath() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewSetPath,
		starlarkfn.SetPath,
	)
	return nil
}

func SetVariable() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewSetVariable,
		starlarkfn.SetVariable,
	)
	return nil
}
