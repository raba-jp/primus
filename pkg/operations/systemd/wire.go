//+build wireinject

package systemd

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/systemd/handlers"
	"github.com/raba-jp/primus/pkg/operations/systemd/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func EnableService() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewEnableService,
		starlarkfn.EnableService,
	)
	return nil
}

func StartService() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewStartService,
		starlarkfn.StartService,
	)
	return nil
}
