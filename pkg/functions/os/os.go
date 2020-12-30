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
	fnNum := 5
	dict := lib.NewDict(fnNum)
	_ = dict.SetKey(lib.String("is_arch_linux"), lib.NewBuiltin("is_arch_linux", NewIsArchFunction(checker)))
	if checker(context.Background()) {
		_ = dict.SetKey(lib.String("installed"), lib.NewBuiltin("installed", NewArchInstalledFunction(installed)))
		_ = dict.SetKey(lib.String("install"), lib.NewBuiltin("install", NewArchInstallFunction(install)))
		_ = dict.SetKey(lib.String("multiple_install"), lib.NewBuiltin("multiple_install", NewArchMultipleInstallFunction(multipleInstall)))
		_ = dict.SetKey(lib.String("uninstall"), lib.NewBuiltin("uninstall", NewArchUninstallFunction(uninstall)))
	} else {
		_ = dict.SetKey(lib.String("installed"), lib.NewBuiltin("installed", dummyFunction()))
		_ = dict.SetKey(lib.String("install"), lib.NewBuiltin("install", dummyFunction()))
		_ = dict.SetKey(lib.String("multiple_install"), lib.NewBuiltin("multiple_install", dummyFunction()))
		_ = dict.SetKey(lib.String("uninstall"), lib.NewBuiltin("uninstall", dummyFunction()))
	}
	return dict
}

func newDarwinFunctions(
	checker backend.DarwinChecker,
	installed DarwinInstalledRunner,
	install DarwinInstallRunner,
	uninstall DarwinUninstallRunner,
) lib.Value {
	fnNum := 4
	dict := lib.NewDict(fnNum)
	_ = dict.SetKey(lib.String("is_darwin"), lib.NewBuiltin("is_darwin", NewIsDarwinFunction(checker)))
	if checker(context.Background()) {
		_ = dict.SetKey(lib.String("installed"), lib.NewBuiltin("installed", NewDarwinInstalledFunction(installed)))
		_ = dict.SetKey(lib.String("install"), lib.NewBuiltin("install", NewDarwinInstallFunction(install)))
		_ = dict.SetKey(lib.String("uninstall"), lib.NewBuiltin("uninstall", NewDarwinUninstallFunction(uninstall)))
	} else {
		_ = dict.SetKey(lib.String("installed"), lib.NewBuiltin("installed", dummyFunction()))
		_ = dict.SetKey(lib.String("install"), lib.NewBuiltin("install", dummyFunction()))
		_ = dict.SetKey(lib.String("uninstall"), lib.NewBuiltin("uninstall", dummyFunction()))
	}
	return dict
}

func newFilePathFunctions() lib.Value {
	fnNum := 3
	dict := lib.NewDict(fnNum)
	_ = dict.SetKey(lib.String("get_current_path"), lib.NewBuiltin("get_current_path", GetCurrentPath()))
	_ = dict.SetKey(lib.String("get_dir"), lib.NewBuiltin("get_dir", GetDir()))
	_ = dict.SetKey(lib.String("join_path"), lib.NewBuiltin("join_path", JoinPath()))
	return dict
}
