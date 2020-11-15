package builtin

import (
	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/functions/filesystem"
	"github.com/raba-jp/primus/pkg/functions/os"
	lib "go.starlark.net/starlark"
)

type Predeclared = lib.StringDict

func NewBuiltinFunction() Predeclared {
	return lib.StringDict{
		"command":    command.NewFunctions(),
		"filesystem": filesystem.NewFunctions(),
		"arch":       os.NewArchFunctions(),
		"darwin":     os.NewDarwinFunctions(),
	}
}
