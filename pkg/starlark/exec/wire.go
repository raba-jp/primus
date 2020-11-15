//+build wireinject

package exec

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/modules"
	"github.com/raba-jp/primus/pkg/starlark/builtin"
)

func Initialize() Fn {
	wire.Build(
		modules.NewFs,
		builtin.NewBuiltinFunction,
		NewExecFn,
	)
	return nil
}
