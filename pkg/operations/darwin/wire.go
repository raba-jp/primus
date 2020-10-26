//+build wireinject

package darwin

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/darwin/handlers"
	"github.com/raba-jp/primus/pkg/operations/darwin/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func CheckInstall() starlark.Fn {
	wire.Build(
		backend.NewFs,
		backend.NewExecInterface,
		handlers.NewCheckInstall,
		starlarkfn.CheckInstall,
	)
	return nil
}

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

func Uninstall() starlark.Fn {
	wire.Build(
		backend.NewFs,
		backend.NewExecInterface,
		handlers.NewCheckInstall,
		handlers.NewUninstall,
		starlarkfn.Uninstall,
	)
	return nil
}
