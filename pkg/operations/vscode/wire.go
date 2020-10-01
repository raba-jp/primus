//+build wireinject

package vscode

import (
	"github.com/google/wire"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/operations/vscode/handlers"
	"github.com/raba-jp/primus/pkg/operations/vscode/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func InstallExtension() starlark.Fn {
	wire.Build(
		backend.NewExecInterface,
		handlers.NewInstallExtension,
		starlarkfn.InstallExtension,
	)
	return nil
}
