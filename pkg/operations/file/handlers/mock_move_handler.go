// Code generated by mockery v1.0.0. DO NOT EDIT.

package handlers

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockMoveHandler is an autogenerated mock type for the MoveHandler type
type MockMoveHandler struct {
	mock.Mock
}

type MoveHandlerMoveArgs struct {
	Ctx            context.Context
	CtxAnything    bool
	Dryrun         bool
	DryrunAnything bool
	P              *MoveParams
	PAnything      bool
}

type MoveHandlerMoveReturns struct {
	Err error
}

type MoveHandlerMoveExpectation struct {
	Args    MoveHandlerMoveArgs
	Returns MoveHandlerMoveReturns
}

func (_m *MockMoveHandler) ApplyMoveExpectation(e MoveHandlerMoveExpectation) {
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
	_m.On("Move", args...).Return(e.Returns.Err)
}

func (_m *MockMoveHandler) ApplyMoveExpectations(expectations []MoveHandlerMoveExpectation) {
	for _, e := range expectations {
		_m.ApplyMoveExpectation(e)
	}
}

// Move provides a mock function with given fields: ctx, dryrun, p
func (_m *MockMoveHandler) Move(ctx context.Context, dryrun bool, p *MoveParams) error {
	ret := _m.Called(ctx, dryrun, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, bool, *MoveParams) error); ok {
		r0 = rf(ctx, dryrun, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
