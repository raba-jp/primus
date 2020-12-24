//+build wireinject

package fish

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/modules"
	lib "go.starlark.net/starlark"
)

func newFunctions(
	setVariable SetVariableRunner,
	setPath SetPathRunner,
) lib.Value {
	dict := lib.NewDict(2)
	dict.SetKey(lib.String("set_variable"), lib.NewBuiltin("set_variable", NewSetVariableFunction(setVariable)))
	dict.SetKey(lib.String("set_path"), lib.NewBuiltin("set_path", NewSetPathFunction(setPath)))
	return dict
}

func NewFunctions() lib.Value {
	wire.Build(
		modules.NewExecInterface,
		command.Execute,
		SetVariable,
		SetPath,
		newFunctions,
	)
	return nil
}
