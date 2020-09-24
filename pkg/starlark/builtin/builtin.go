package builtin

import (
	"github.com/raba-jp/primus/pkg/operations/command"
	"github.com/raba-jp/primus/pkg/operations/directory"
	"github.com/raba-jp/primus/pkg/operations/file"
	"github.com/raba-jp/primus/pkg/operations/fish"
	"github.com/raba-jp/primus/pkg/operations/network"
	"github.com/raba-jp/primus/pkg/operations/os"
	"github.com/raba-jp/primus/pkg/operations/packages"
	lib "go.starlark.net/starlark"
)

type Predeclared = lib.StringDict

func NewBuiltinFunction() Predeclared {
	return lib.StringDict{
		"command":           lib.NewBuiltin("command", command.Command()),
		"symlink":           lib.NewBuiltin("symlink", file.Symlink()),
		"http_request":      lib.NewBuiltin("http_request", network.HTTPRequest()),
		"pkg_install":       lib.NewBuiltin("pkg_install", packages.Install()),
		"copy_file":         lib.NewBuiltin("copy_file", file.Copy()),
		"move_file":         lib.NewBuiltin("move_file", file.Move()),
		"is_darwin":         lib.NewBuiltin("is_darwin", os.IsDarwin()),
		"is_arch_linux":     lib.NewBuiltin("is_arch_linux", os.IsArchLinux()),
		"fish_set_variable": lib.NewBuiltin("fish_set_variable", fish.SetVariable()),
		"fish_set_path":     lib.NewBuiltin("fish_set_path", fish.SetPath()),
		"create_directory":  lib.NewBuiltin("create_directory", directory.Create()),
	}
}
