//+build wireinject

package repl

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/internal/backend"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
)

func Initialize() REPL {
	wire.Build(
		NewState,
		newThread,
		NewREPL,
		backend.NewFs,
		backend.NewExecInterface,
		backend.New,
		functions.NewPredeclaredFunction,
	)
	return nil
}
