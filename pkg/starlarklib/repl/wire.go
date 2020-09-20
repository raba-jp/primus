//+build wireinject

package repl

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/starlarklib/functions"
)

func Initialize() PromptFunc {
	wire.Build(
		NewPrompt,
		NewState,
		newThread,
		NewREPL,
		NewExecutor,
		NewCompleter,
		backend.NewFs,
		backend.NewExecInterface,
		backend.New,
		functions.NewPredeclaredFunction,
	)
	return nil
}
