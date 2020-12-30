// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package starlark

import (
	"context"
	"github.com/raba-jp/primus/pkg/backend"
	"go.starlark.net/starlark"
)

// Injectors from wire.go:

func Initialize() func(ctx context.Context, path string, predeclared starlark.StringDict) error {
	fs := backend.NewFs()
	v := NewExecFn(fs)
	return v
}
