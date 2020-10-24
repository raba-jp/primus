// Code generated by mockery v1.0.0. DO NOT EDIT.

package handlers

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockSetPathHandler is an autogenerated mock type for the SetPathHandler type
type MockSetPathHandler struct {
	mock.Mock
}

type SetPathHandlerSetPathArgs struct {
	Ctx            context.Context
	CtxAnything    bool
	Dryrun         bool
	DryrunAnything bool
	P              *SetPathParams
	PAnything      bool
}

type SetPathHandlerSetPathReturns struct {
	Err error
}

type SetPathHandlerSetPathExpectation struct {
	Args    SetPathHandlerSetPathArgs
	Returns SetPathHandlerSetPathReturns
}

func (_m *MockSetPathHandler) ApplySetPathExpectation(e SetPathHandlerSetPathExpectation) {
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
	_m.On("SetPath", args...).Return(e.Returns.Err)
}

func (_m *MockSetPathHandler) ApplySetPathExpectations(expectations []SetPathHandlerSetPathExpectation) {
	for _, e := range expectations {
		_m.ApplySetPathExpectation(e)
	}
}

// SetPath provides a mock function with given fields: ctx, dryrun, p
func (_m *MockSetPathHandler) SetPath(ctx context.Context, dryrun bool, p *SetPathParams) error {
	ret := _m.Called(ctx, dryrun, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, bool, *SetPathParams) error); ok {
		r0 = rf(ctx, dryrun, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
