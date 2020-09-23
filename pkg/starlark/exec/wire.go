//+build wireinject

package exec

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/starlark/builtin"
)

func Initialize() Fn {
	wire.Build(
		backend.NewFs,
		backend.NewExecInterface,
		backend.New,
		builtin.NewBuiltinFunction,
		NewExecFn,
	)
	return nil
}
