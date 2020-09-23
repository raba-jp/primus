// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package exec

import (
	"context"
	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/starlark/builtin"
)

// Injectors from wire.go:

func Initialize() func(ctx context.Context, dryrun bool, path string) error {
	execInterface := backend.NewExecInterface()
	fs := backend.NewFs()
	backendBackend := backend.New(execInterface, fs)
	stringDict := builtin.NewBuiltinFunction(backendBackend, execInterface, fs)
	v := NewExecFn(stringDict, fs)
	return v
}