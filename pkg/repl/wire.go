//+build wireinject

package repl

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/functions"
)

func Initialize() PromptFunc {
	wire.Build(
		NewPrompt,
		NewState,
		newThread,
		NewREPL,
		NewExecutor,
		NewCompleter,
		functions.New,
	)
	return nil
}
