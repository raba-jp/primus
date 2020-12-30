//+build wireinject

package os

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	lib "go.starlark.net/starlark"
)

func NewArchFunctions() lib.Value {
	wire.Build(
		backend.NewFs,
		backend.NewExecutable,
		backend.NewExecute,
		backend.NewArchLinuxChecker,
		ArchInstalled,
		ArchInstall,
		ArchMultipleInstall,
		ArchUninstall,
		newArchFunctions,
	)
	return nil
}

func NewDarwinFunctions() lib.Value {
	wire.Build(
		backend.NewFs,
		backend.NewExecute,
		backend.NewDarwinChecker,
		DarwinInstalled,
		DarwinInstall,
		DarwinUninstall,
		newDarwinFunctions,
	)
	return nil
}

func NewFilePathFunctions() lib.Value {
	wire.Build(newFilePathFunctions)
	return nil
}
