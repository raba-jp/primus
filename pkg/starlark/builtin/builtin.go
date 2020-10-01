package builtin

import (
	"github.com/raba-jp/primus/pkg/operations/command"
	"github.com/raba-jp/primus/pkg/operations/directory"
	"github.com/raba-jp/primus/pkg/operations/file"
	"github.com/raba-jp/primus/pkg/operations/filepath"
	"github.com/raba-jp/primus/pkg/operations/fish"
	"github.com/raba-jp/primus/pkg/operations/network"
	"github.com/raba-jp/primus/pkg/operations/os"
	"github.com/raba-jp/primus/pkg/operations/packages"
	"github.com/raba-jp/primus/pkg/operations/systemd"
	"github.com/raba-jp/primus/pkg/operations/vscode"
	lib "go.starlark.net/starlark"
)

type Predeclared = lib.StringDict

func NewBuiltinFunction() Predeclared {
	return lib.StringDict{
		"command":                  lib.NewBuiltin("command", command.Command()),
		"symlink":                  lib.NewBuiltin("symlink", file.Symlink()),
		"http_request":             lib.NewBuiltin("http_request", network.HTTPRequest()),
		"copy_file":                lib.NewBuiltin("copy_file", file.Copy()),
		"move_file":                lib.NewBuiltin("move_file", file.Move()),
		"is_darwin":                lib.NewBuiltin("is_darwin", os.IsDarwin()),
		"is_arch_linux":            lib.NewBuiltin("is_arch_linux", os.IsArchLinux()),
		"fish_set_variable":        lib.NewBuiltin("fish_set_variable", fish.SetVariable()),
		"fish_set_path":            lib.NewBuiltin("fish_set_path", fish.SetPath()),
		"create_directory":         lib.NewBuiltin("create_directory", directory.Create()),
		"enable_service":           lib.NewBuiltin("enable_service", systemd.EnableService()),
		"start_service":            lib.NewBuiltin("start_service", systemd.StartService()),
		"darwin_pkg_check_install": lib.NewBuiltin("darwin_pkg_check_install", packages.DarwinPkgCheckInstall()),
		"darwin_pkg_install":       lib.NewBuiltin("darwin_pkg_install", packages.DarwinPkgInstall()),
		"darwin_pkg_uninstall":     lib.NewBuiltin("darwin_pkg_uninstall", packages.DarwinPkgUninstall()),
		"arch_pkg_check_install":   lib.NewBuiltin("arch_pkg_check_install", packages.ArchPkgCheckInstall()),
		"arch_pkg_install":         lib.NewBuiltin("arch_pkg_install", packages.ArchPkgInstall()),
		"arch_pkg_uninstall":       lib.NewBuiltin("darwin_pkg_install", packages.ArchPkgUninstall()),
		"current_filepath":         lib.NewBuiltin("current_filepath", filepath.Current()),
		"get_dir":                  lib.NewBuiltin("get_dir", filepath.Dir()),
		"join_filepath":            lib.NewBuiltin("join_filepath", filepath.Join()),
		"vscode_install_extension": lib.NewBuiltin("vscode_install_extension", vscode.InstallExtension()),
	}
}
