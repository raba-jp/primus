//+build wireinject

package command

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/modules"
	lib "go.starlark.net/starlark"
)

func newFunctions(exc ExecuteRunner, detector modules.OSDetector) lib.Value {
	dict := lib.NewDict(2)
	dict.SetKey(lib.String("execute"), lib.NewBuiltin("execute", NewExecuteFunction(exc)))
	dict.SetKey(lib.String("executable"), lib.NewBuiltin("executable", NewExecutableFunction(detector)))
	return dict
}

func NewFunctions() lib.Value {
	wire.Build(
		modules.NewExecInterface,
		modules.NewFs,
		modules.NewOSDetector,
		Execute,
		newFunctions,
	)
	return nil
}
