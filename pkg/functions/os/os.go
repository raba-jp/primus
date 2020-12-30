package os

import (
	"context"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
)

func dummyFunction() starlark.Fn {
	return func(thread *lib.Thread, b *lib.Builtin, args lib.Tuple, kwargs []lib.Tuple) (lib.Value, error) {
		return lib.None, nil
	}
}

func newArchFunctions(
	checker backend.ArchLinuxChecker,
	installed ArchInstalledRunner,
	install ArchInstallRunner,
	multipleInstall ArchMultipleInstallRunner,
	uninstall ArchUninstallRunner,
) lib.Value {
	dict := lib.NewDict(4)
	if checker(context.Background()) {
		dict.SetKey(lib.String("installed"), lib.NewBuiltin("installed", NewArchInstalledFunction(installed)))
		dict.SetKey(lib.String("install"), lib.NewBuiltin("install", NewArchInstallFunction(install)))
		dict.SetKey(lib.String("multiple_install"), lib.NewBuiltin("multiple_install", NewArchMultipleInstallFunction(multipleInstall)))
		dict.SetKey(lib.String("uninstall"), lib.NewBuiltin("uninstall", NewArchUninstallFunction(uninstall)))
	} else {
		dict.SetKey(lib.String("installed"), lib.NewBuiltin("installed", dummyFunction()))
		dict.SetKey(lib.String("install"), lib.NewBuiltin("install", dummyFunction()))
		dict.SetKey(lib.String("multiple_install"), lib.NewBuiltin("multiple_install", dummyFunction()))
		dict.SetKey(lib.String("uninstall"), lib.NewBuiltin("uninstall", dummyFunction()))
	}
	return dict
}

func newDarwinFunctions(
	checker backend.DarwinChecker,
	installed DarwinInstalledRunner,
	install DarwinInstallRunner,
	uninstall DarwinUninstallRunner,
) lib.Value {
	dict := lib.NewDict(3)
	if checker(context.Background()) {
		dict.SetKey(lib.String("installed"), lib.NewBuiltin("installed", NewDarwinInstalledFunction(installed)))
		dict.SetKey(lib.String("install"), lib.NewBuiltin("install", NewDarwinInstallFunction(install)))
		dict.SetKey(lib.String("uninstall"), lib.NewBuiltin("uninstall", NewDarwinUninstallFunction(uninstall)))
	} else {
		dict.SetKey(lib.String("installed"), lib.NewBuiltin("installed", dummyFunction()))
		dict.SetKey(lib.String("install"), lib.NewBuiltin("install", dummyFunction()))
		dict.SetKey(lib.String("uninstall"), lib.NewBuiltin("uninstall", dummyFunction()))
	}
	return dict
}

func newFilePathFunctions() lib.Value {
	dict := lib.NewDict(3)
	dict.SetKey(lib.String("get_current_path"), lib.NewBuiltin("get_current_path", GetCurrentPath()))
	dict.SetKey(lib.String("get_dir"), lib.NewBuiltin("get_dir", GetDir()))
	dict.SetKey(lib.String("join_path"), lib.NewBuiltin("join_path", JoinPath()))
	return dict
}
