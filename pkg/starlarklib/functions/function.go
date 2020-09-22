package functions

import (
	"context"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/exec"
	"github.com/raba-jp/primus/pkg/starlarklib"
	"github.com/spf13/afero"
	"go.starlark.net/starlark"
	"golang.org/x/xerrors"
)

type StarlarkFn = func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error)
type Predeclared = starlark.StringDict
type ExecFileFn = func(ctx context.Context, dryrun bool, path string) error

const retValue = starlark.None

func NewPredeclaredFunction(be backend.Backend, execIF exec.Interface, fs afero.Fs) Predeclared {
	return starlark.StringDict{
		"command":           starlark.NewBuiltin("command", Command(be)),
		"symlink":           starlark.NewBuiltin("symlink", Symlink(be)),
		"http_request":      starlark.NewBuiltin("http_request", HTTPRequest(be)),
		"package":           starlark.NewBuiltin("package", Package(be, be)),
		"copy_file":         starlark.NewBuiltin("copy_file", FileCopy(be)),
		"move_file":         starlark.NewBuiltin("move_file", FileMove(be)),
		"is_darwin":         starlark.NewBuiltin("is_darwin", IsDarwin(execIF)),
		"is_arch_linux":     starlark.NewBuiltin("is_arch_linux", IsArchLinux(fs)),
		"fish_set_variable": starlark.NewBuiltin("fish_set_variable", FishSetVariable(be)),
		"fish_set_path":     starlark.NewBuiltin("fish_set_path", FishSetPath(be)),
		"create_directory":  starlark.NewBuiltin("create_directory", CreateDirectory(be)),
	}
}

func NewExecFileFn(predeclared Predeclared, fs afero.Fs) ExecFileFn {
	return func(ctx context.Context, dryrun bool, path string) error {
		data, err := afero.ReadFile(fs, path)
		if err != nil {
			return xerrors.Errorf("Read file failed: %s: %w", path, err)
		}

		thread := &starlark.Thread{
			Name: "main",
			Load: Load(dryrun, fs, predeclared),
		}
		starlarklib.SetCtx(ctx, thread)
		starlarklib.SetDryRun(thread, dryrun)
		if _, err := starlark.ExecFile(thread, path, data, predeclared); err != nil {
			return xerrors.Errorf("Failed exec file: %w", err)
		}

		return nil
	}
}
