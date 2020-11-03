// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	handlers "github.com/raba-jp/primus/pkg/operations/git/handlers"
	mock "github.com/stretchr/testify/mock"
)

// CloneHandler is an autogenerated mock type for the CloneHandler type
type CloneHandler struct {
	mock.Mock
}

type CloneHandlerRunArgs struct {
	Ctx         context.Context
	CtxAnything bool
	P           *handlers.CloneParams
	PAnything   bool
}

type CloneHandlerRunReturns struct {
	Err error
}

type CloneHandlerRunExpectation struct {
	Args    CloneHandlerRunArgs
	Returns CloneHandlerRunReturns
}

func (_m *CloneHandler) ApplyRunExpectation(e CloneHandlerRunExpectation) {
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

func (_m *CloneHandler) ApplyRunExpectations(expectations []CloneHandlerRunExpectation) {
	for _, e := range expectations {
		_m.ApplyRunExpectation(e)
	}
}

// Run provides a mock function with given fields: ctx, p
func (_m *CloneHandler) Run(ctx context.Context, p *handlers.CloneParams) error {
	ret := _m.Called(ctx, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *handlers.CloneParams) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
