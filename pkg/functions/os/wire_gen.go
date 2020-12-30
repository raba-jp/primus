// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package os

import (
	"github.com/raba-jp/primus/pkg/backend"
	"go.starlark.net/starlark"
)

// Injectors from wire.go:

func NewArchFunctions() starlark.Value {
	fs := backend.NewFs()
	archLinuxChecker := backend.NewArchLinuxChecker(fs)
	executable := backend.NewExecutable()
	execute := backend.NewExecute()
	archInstalledRunner := ArchInstalled(executable, execute)
	archInstallRunner := ArchInstall(executable, execute)
	archMultipleInstallRunner := ArchMultipleInstall(executable, execute)
	archUninstallRunner := ArchUninstall(executable, execute)
	value := newArchFunctions(archLinuxChecker, archInstalledRunner, archInstallRunner, archMultipleInstallRunner, archUninstallRunner)
	return value
}

func NewDarwinFunctions() starlark.Value {
	execute := backend.NewExecute()
	darwinChecker := backend.NewDarwinChecker(execute)
	fs := backend.NewFs()
	darwinInstalledRunner := DarwinInstalled(execute, fs)
	darwinInstallRunner := DarwinInstall(execute, fs)
	darwinUninstallRunner := DarwinUninstall(execute, fs)
	value := newDarwinFunctions(darwinChecker, darwinInstalledRunner, darwinInstallRunner, darwinUninstallRunner)
	return value
}

func NewFilePathFunctions() starlark.Value {
	value := newFilePathFunctions()
	return value
}
