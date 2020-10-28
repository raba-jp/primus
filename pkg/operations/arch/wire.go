//+build wireinject

package arch

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/arch/handlers"
	"github.com/raba-jp/primus/pkg/operations/arch/starlarkfn"
	commandHandlers "github.com/raba-jp/primus/pkg/operations/command/handlers"
	"github.com/raba-jp/primus/pkg/starlark"
)

func CheckInstall() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewCheckInstall,
		starlarkfn.CheckInstall,
	)
	return nil
}

func Install() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		backend.NewFs,
		commandHandlers.NewExecutable,
		handlers.NewCheckInstall,
		handlers.NewInstall,
		starlarkfn.Install,
	)
	return nil
}

func Uninstall() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewCheckInstall,
		handlers.NewUninstall,
		starlarkfn.Uninstall,
	)
	return nil
}

func MultipleInstall() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		backend.NewFs,
		commandHandlers.NewExecutable,
		handlers.NewMultipleInstall,
		starlarkfn.MultipleInstall,
	)
	return nil
}
