//+build wireinject

package os

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/os/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func IsDarwin() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		starlarkfn.IsDarwin,
	)
	return nil
}

func IsArchLinux() starlark.Fn {
	wire.Build(
		backend.NewFs,
		starlarkfn.IsArchLinux,
	)
	return nil
}
