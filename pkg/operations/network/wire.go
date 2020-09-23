//+build wireinject

package network

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/network/handlers"
	"github.com/raba-jp/primus/pkg/operations/network/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func HTTPRequest() starlark.Fn {
	wire.Build(
		backend.NewFs,
		backend.NewHTTPClient,
		handlers.NewHTTPRequest,
		starlarkfn.HTTPRequest,
	)
	return nil
}
