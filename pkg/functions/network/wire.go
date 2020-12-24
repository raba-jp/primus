//+build wireinject

package network

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/modules"
	lib "go.starlark.net/starlark"
)

func newFunctions(clone GitCloneRunner, httpRequest HTTPRequestRunner) lib.Value {
	dict := lib.NewDict(2)
	dict.SetKey(lib.String("git_clone"), lib.NewBuiltin("git_clone", NewGitCloneFunction(clone)))
	dict.SetKey(lib.String("http_request"), lib.NewBuiltin("http_request", NewHTTPRequestFunction(httpRequest)))
	return dict
}

func NewFunctions() lib.Value {
	wire.Build(
		modules.NewExecInterface,
		modules.NewHTTPClient,
		modules.NewFs,
		command.Execute,
		GitClone,
		HTTPRequest,
		newFunctions,
	)
	return nil
}
