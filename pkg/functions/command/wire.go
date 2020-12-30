//+build wireinject

package command

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	lib "go.starlark.net/starlark"
)

func newFunctions(execute backend.Execute, executable backend.Executable) lib.Value {
	dict := lib.NewDict(2)
	dict.SetKey(lib.String("execute"), lib.NewBuiltin("execute", NewExecuteFunction(execute)))
	dict.SetKey(lib.String("executable"), lib.NewBuiltin("executable", NewExecutableFunction(executable)))
	return dict
}

func NewFunctions() lib.Value {
	wire.Build(
		backend.NewExecutable,
		backend.NewExecute,
		newFunctions,
	)
	return nil
}
