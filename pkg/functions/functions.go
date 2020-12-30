package functions

import (
	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/functions/filesystem"
	"github.com/raba-jp/primus/pkg/functions/fish"
	"github.com/raba-jp/primus/pkg/functions/network"
	"github.com/raba-jp/primus/pkg/functions/os"
	lib "go.starlark.net/starlark"
)

func New() lib.StringDict {
	return lib.StringDict{
		"command":    command.NewFunctions(),
		"filesystem": filesystem.NewFunctions(),
		"arch":       os.NewArchFunctions(),
		"darwin":     os.NewDarwinFunctions(),
		"filepath":   os.NewFilePathFunctions(),
		"fish":       fish.NewFunctions(),
		"network":    network.NewFunctions(),
	}
}
