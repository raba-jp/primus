//+build wireinject

package os

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/modules"
	lib "go.starlark.net/starlark"
)

func newArchFunctions(
	installed ArchInstalledRunner,
	install ArchInstallRunner,
	multipleInstall ArchMultipleInstallRunner,
	uninstall ArchUninstallRunner,
) lib.Value {
	dict := lib.NewDict(4)
	dict.SetKey(lib.String("installed"), lib.NewBuiltin("installed", NewArchInstalledFunction(installed)))
	dict.SetKey(lib.String("install"), lib.NewBuiltin("install", NewArchInstallFunction(install)))
	dict.SetKey(lib.String("multiple_install"), lib.NewBuiltin("multiple_install", NewArchMultipleInstallFunction(multipleInstall)))
	dict.SetKey(lib.String("uninstall"), lib.NewBuiltin("uninstall", NewArchUninstallFunction(uninstall)))
	return dict
}

func newDarwinFunctions(
	installed DarwinInstalledRunner,
	install DarwinInstallRunner,
	uninstall DarwinUninstallRunner,
) lib.Value {
	dict := lib.NewDict(3)
	dict.SetKey(lib.String("installed"), lib.NewBuiltin("installed", NewDarwinInstalledFunction(installed)))
	dict.SetKey(lib.String("install"), lib.NewBuiltin("install", NewDarwinInstallFunction(install)))
	dict.SetKey(lib.String("uninstall"), lib.NewBuiltin("uninstall", NewDarwinUninstallFunction(uninstall)))
	return dict
}

func newFilePathFunctions() lib.Value {
	dict := lib.NewDict(3)
	dict.SetKey(lib.String("get_current_path"), lib.NewBuiltin("get_current_path", GetCurrentPath()))
	dict.SetKey(lib.String("get_dir"), lib.NewBuiltin("get_dir", GetDir()))
	dict.SetKey(lib.String("join_path"), lib.NewBuiltin("join_path", JoinPath()))
	return dict
}

func NewArchFunctions() lib.Value {
	wire.Build(
		modules.NewExecInterface,
		command.Executable,
		command.Execute,
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
		modules.NewFs,
		modules.NewExecInterface,
		command.Execute,
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
