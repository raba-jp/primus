//+build wireinject

package special

import (
	"io"
	"os"

	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/operations/special/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func reader() io.Reader {
	return os.Stdin
}

func RequirePrevilegedAccess() starlark.Fn {
	wire.Build(
		reader,
		starlarkfn.RequirePrevilegedAccess,
	)
	return nil
}

func PrintContext() starlark.Fn {
	wire.Build(
		starlarkfn.PrintContext,
	)
	return nil
}
