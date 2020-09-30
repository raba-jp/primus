//+build wireinject

package packages

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/packages/handlers"
	"github.com/raba-jp/primus/pkg/operations/packages/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func DarwinPkgCheckInstall() starlark.Fn {
	wire.Build(
		backend.NewFs,
		backend.NewExecInterface,
		handlers.NewDarwinPkgCheckInstallHandler,
		starlarkfn.DarwinPkgCheckInstall,
	)
	return nil
}

func DarwinPkgInstall() starlark.Fn {
	wire.Build(
		backend.NewFs,
		backend.NewExecInterface,
		handlers.NewDarwinPkgInstallHandler,
		starlarkfn.DarwinPkgInstall,
	)
	return nil
}

func DarwinPkgUninstall() starlark.Fn {
	wire.Build(
		backend.NewFs,
		backend.NewExecInterface,
		handlers.NewDarwinPkgUninstallHandler,
		starlarkfn.DarwinPkgUninstall,
	)
	return nil
}

func ArchPkgCheckInstall() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewArchPkgCheckInstallHandler,
		starlarkfn.ArchPkgCheckInstall,
	)
	return nil
}

func ArchPkgInstall() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewArchPkgInstallHandler,
		starlarkfn.ArchPkgInstall,
	)
	return nil
}

func ArchPkgUninstall() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewArchPkgUninstallHandler,
		starlarkfn.ArchPkgUninstall,
	)
	return nil
}
