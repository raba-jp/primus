// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package special

import (
	"go.starlark.net/starlark"
	"io"
	"os"
)

// Injectors from wire.go:

func RequirePrevilegedAccess() func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
	ioReader := reader()
	v := NewRequirePrevilegedAccessFunction(ioReader)
	return v
}

func PrintContext() func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kargs []starlark.Tuple) (starlark.Value, error) {
	v := NewPrintContextFunction()
	return v
}

// wire.go:

func reader() io.Reader {
	return os.Stdin
}