// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package fish

import (
	"github.com/raba-jp/primus/pkg/backend"
	"go.starlark.net/starlark"
)

// Injectors from wire.go:

func NewFunctions() starlark.Value {
	execute := backend.NewExecute()
	setVariableRunner := SetVariable(execute)
	setPathRunner := SetPath(execute)
	value := newFunctions(setVariableRunner, setPathRunner)
	return value
}

// wire.go:

func newFunctions(
	setVariable SetVariableRunner,
	setPath SetPathRunner,
) starlark.Value {
	dict := starlark.NewDict(2)
	dict.SetKey(starlark.String("set_variable"), starlark.NewBuiltin("set_variable", NewSetVariableFunction(setVariable)))
	dict.SetKey(starlark.String("set_path"), starlark.NewBuiltin("set_path", NewSetPathFunction(setPath)))
	return dict
}
