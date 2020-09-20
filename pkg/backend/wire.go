//+build wireinject

package backend

import (
	"github.com/google/wire"
)

func Initialize() Backend {
	wire.Build(
		NewFs,
		NewExecInterface,
		New,
	)
	return nil
}
