// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package functions

import (
	"context"
	"github.com/raba-jp/primus/pkg/internal/backend"
)

// Injectors from wire.go:

func Initialize() func(ctx context.Context, dryrun bool, path string) error {
	execInterface := backend.NewExecInterface()
	fs := backend.NewFs()
	backendBackend := backend.New(execInterface, fs)
	stringDict := NewPredeclaredFunction(backendBackend, execInterface, fs)
	v := NewExecFileFn(stringDict, fs)
	return v
}
