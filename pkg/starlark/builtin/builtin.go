package builtin

import (
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/starlark/builtin/command"
	"github.com/raba-jp/primus/pkg/starlark/builtin/directory"
	"github.com/raba-jp/primus/pkg/starlark/builtin/file"
	"github.com/raba-jp/primus/pkg/starlark/builtin/fish"
	"github.com/raba-jp/primus/pkg/starlark/builtin/network"
	"github.com/raba-jp/primus/pkg/starlark/builtin/os"
	"github.com/raba-jp/primus/pkg/starlark/builtin/packages"
	"github.com/spf13/afero"
	lib "go.starlark.net/starlark"
)

type Predeclared = lib.StringDict

func NewBuiltinFunction(be backend.Backend, execIF exec.Interface, fs afero.Fs) Predeclared {
	return lib.StringDict{
		"command":           lib.NewBuiltin("command", command.Command(be)),
		"symlink":           lib.NewBuiltin("symlink", file.Symlink(be)),
		"http_request":      lib.NewBuiltin("http_request", network.HTTPRequest(be)),
		"package":           lib.NewBuiltin("package", packages.Install(be, be)),
		"copy_file":         lib.NewBuiltin("copy_file", file.Copy(be)),
		"move_file":         lib.NewBuiltin("move_file", file.Move(be)),
		"is_darwin":         lib.NewBuiltin("is_darwin", os.IsDarwin(execIF)),
		"is_arch_linux":     lib.NewBuiltin("is_arch_linux", os.IsArchLinux(fs)),
		"fish_set_variable": lib.NewBuiltin("fish_set_variable", fish.SetVariable(be)),
		"fish_set_path":     lib.NewBuiltin("fish_set_path", fish.SetPath(be)),
		"create_directory":  lib.NewBuiltin("create_directory", directory.Create(be)),
	}
}
