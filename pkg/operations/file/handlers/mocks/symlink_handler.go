// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	handlers "github.com/raba-jp/primus/pkg/operations/file/handlers"
	mock "github.com/stretchr/testify/mock"
)

// SymlinkHandler is an autogenerated mock type for the SymlinkHandler type
type SymlinkHandler struct {
	mock.Mock
}

type SymlinkHandlerRunArgs struct {
	Ctx         context.Context
	CtxAnything bool
	P           *handlers.SymlinkParams
	PAnything   bool
}

type SymlinkHandlerRunReturns struct {
	Err error
}

type SymlinkHandlerRunExpectation struct {
	Args    SymlinkHandlerRunArgs
	Returns SymlinkHandlerRunReturns
}

func (_m *SymlinkHandler) ApplyRunExpectation(e SymlinkHandlerRunExpectation) {
	var args []interface{}
	if e.Args.CtxAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Ctx)
	}
	if e.Args.PAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.P)
	}
	_m.On("Run", args...).Return(e.Returns.Err)
}

func (_m *SymlinkHandler) ApplyRunExpectations(expectations []SymlinkHandlerRunExpectation) {
	for _, e := range expectations {
		_m.ApplyRunExpectation(e)
	}
}

// Run provides a mock function with given fields: ctx, p
func (_m *SymlinkHandler) Run(ctx context.Context, p *handlers.SymlinkParams) error {
	ret := _m.Called(ctx, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *handlers.SymlinkParams) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
