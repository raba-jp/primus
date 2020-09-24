//+build wireinject

package packages

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/packages/handlers"
	"github.com/raba-jp/primus/pkg/operations/packages/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func Install() starlark.Fn {
	wire.Build(
		backend.NewFs,
		backend.NewExecInterface,
		handlers.NewCheckInstall,
		handlers.NewInstall,
		starlarkfn.Install,
	)
	return nil
}
