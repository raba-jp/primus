// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	handlers "github.com/raba-jp/primus/pkg/operations/arch/handlers"
	mock "github.com/stretchr/testify/mock"
)

// UninstallHandler is an autogenerated mock type for the UninstallHandler type
type UninstallHandler struct {
	mock.Mock
}

type UninstallHandlerRunArgs struct {
	Ctx         context.Context
	CtxAnything bool
	P           *handlers.UninstallParams
	PAnything   bool
}

type UninstallHandlerRunReturns struct {
	Err error
}

type UninstallHandlerRunExpectation struct {
	Args    UninstallHandlerRunArgs
	Returns UninstallHandlerRunReturns
}

func (_m *UninstallHandler) ApplyRunExpectation(e UninstallHandlerRunExpectation) {
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

func (_m *UninstallHandler) ApplyRunExpectations(expectations []UninstallHandlerRunExpectation) {
	for _, e := range expectations {
		_m.ApplyRunExpectation(e)
	}
}

// Run provides a mock function with given fields: ctx, p
func (_m *UninstallHandler) Run(ctx context.Context, p *handlers.UninstallParams) error {
	ret := _m.Called(ctx, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *handlers.UninstallParams) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
