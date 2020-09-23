//+build wireinject

package repl

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/starlark/builtin"
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
		builtin.NewBuiltinFunction,
	)
	return nil
}
