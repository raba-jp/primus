//+build wireinject

package functions

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/internal/backend"
)

func Initialize() ExecFileFn {
	wire.Build(
		backend.NewFs,
		backend.NewExecInterface,
		backend.New,
		NewPredeclaredFunction,
		NewExecFileFn,
	)
	return nil
}
