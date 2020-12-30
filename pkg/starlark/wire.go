//+build wireinject

package starlark

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
)

func Initialize() Exec {
	wire.Build(
		backend.NewFs,
		NewExecFn,
	)
	return nil
}
