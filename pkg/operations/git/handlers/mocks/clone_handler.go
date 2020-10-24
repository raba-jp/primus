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

type CloneHandlerCloneArgs struct {
	Ctx            context.Context
	CtxAnything    bool
	Dryrun         bool
	DryrunAnything bool
	P              *handlers.CloneParams
	PAnything      bool
}

type CloneHandlerCloneReturns struct {
	Err error
}

type CloneHandlerCloneExpectation struct {
	Args    CloneHandlerCloneArgs
	Returns CloneHandlerCloneReturns
}

func (_m *CloneHandler) ApplyCloneExpectation(e CloneHandlerCloneExpectation) {
	var args []interface{}
	if e.Args.CtxAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Ctx)
	}
	if e.Args.DryrunAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Dryrun)
	}
	if e.Args.PAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.P)
	}
	_m.On("Clone", args...).Return(e.Returns.Err)
}

func (_m *CloneHandler) ApplyCloneExpectations(expectations []CloneHandlerCloneExpectation) {
	for _, e := range expectations {
		_m.ApplyCloneExpectation(e)
	}
}

// Clone provides a mock function with given fields: ctx, dryrun, p
func (_m *CloneHandler) Clone(ctx context.Context, dryrun bool, p *handlers.CloneParams) error {
	ret := _m.Called(ctx, dryrun, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, bool, *handlers.CloneParams) error); ok {
		r0 = rf(ctx, dryrun, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
