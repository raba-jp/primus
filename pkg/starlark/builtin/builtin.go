package builtin

import (
	"github.com/raba-jp/primus/pkg/operations/arch"
	"github.com/raba-jp/primus/pkg/operations/command"
	"github.com/raba-jp/primus/pkg/operations/darwin"
	"github.com/raba-jp/primus/pkg/operations/directory"
	"github.com/raba-jp/primus/pkg/operations/file"
	"github.com/raba-jp/primus/pkg/operations/filepath"
	"github.com/raba-jp/primus/pkg/operations/fish"
	"github.com/raba-jp/primus/pkg/operations/git"
	"github.com/raba-jp/primus/pkg/operations/network"
	"github.com/raba-jp/primus/pkg/operations/os"
	"github.com/raba-jp/primus/pkg/operations/special"
	"github.com/raba-jp/primus/pkg/operations/systemd"
	lib "go.starlark.net/starlark"
)

type Predeclared = lib.StringDict

func NewBuiltinFunction() Predeclared {
	return lib.StringDict{
		"require_previleged_access": lib.NewBuiltin("require_previleged_access", special.RequirePrevilegedAccess()),
		"print_context":             lib.NewBuiltin("print_context", special.PrintContext()),
		"command":                   lib.NewBuiltin("command", command.Command()),
		"executable":                lib.NewBuiltin("executable", command.Executable()),
		"symlink":                   lib.NewBuiltin("symlink", file.Symlink()),
		"http_request":              lib.NewBuiltin("http_request", network.HTTPRequest()),
		"copy_file":                 lib.NewBuiltin("copy_file", file.Copy()),
		"move_file":                 lib.NewBuiltin("move_file", file.Move()),
		"is_darwin":                 lib.NewBuiltin("is_darwin", os.IsDarwin()),
		"is_arch_linux":             lib.NewBuiltin("is_arch_linux", os.IsArchLinux()),
		"fish_set_variable":         lib.NewBuiltin("fish_set_variable", fish.SetVariable()),
		"fish_set_path":             lib.NewBuiltin("fish_set_path", fish.SetPath()),
		"create_directory":          lib.NewBuiltin("create_directory", directory.Create()),
		"enable_service":            lib.NewBuiltin("enable_service", systemd.EnableService()),
		"start_service":             lib.NewBuiltin("start_service", systemd.StartService()),
		"darwin_check_install":      lib.NewBuiltin("darwin_check_install", darwin.CheckInstall()),
		"darwin_install":            lib.NewBuiltin("darwin_install", darwin.Install()),
		"darwin_uninstall":          lib.NewBuiltin("darwin_uninstall", darwin.Uninstall()),
		"arch_check_install":        lib.NewBuiltin("arch_check_install", arch.CheckInstall()),
		"arch_install":              lib.NewBuiltin("arch_install", arch.Install()),
		"arch_uninstall":            lib.NewBuiltin("darwin_install", arch.Uninstall()),
		"arch_multiple_install":     lib.NewBuiltin("arch_multiple_install", arch.MultipleInstall()),
		"current_filepath":          lib.NewBuiltin("current_filepath", filepath.Current()),
		"get_dir":                   lib.NewBuiltin("get_dir", filepath.Dir()),
		"join_filepath":             lib.NewBuiltin("join_filepath", filepath.Join()),
		"git_clone":                 lib.NewBuiltin("git_clone", git.Clone()),
	}
}
