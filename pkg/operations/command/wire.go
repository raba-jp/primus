//+build wireinject

package command

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/command/handlers"
	"github.com/raba-jp/primus/pkg/operations/command/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func Command() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewCommand,
		starlarkfn.Command,
	)
	return nil
}

func Executable() starlark.Fn {
	wire.Build(
		backend.NewFs,
		handlers.NewExecutable,
		starlarkfn.Executable,
	)
	return nil
}
