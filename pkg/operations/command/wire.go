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
		handlers.New,
		starlarkfn.Command,
	)
	return nil
}
