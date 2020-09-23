// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package backend

// Injectors from wire.go:

func Initialize() Backend {
	execInterface := NewExecInterface()
	fs := NewFs()
	backend := New(execInterface, fs)
	return backend
}