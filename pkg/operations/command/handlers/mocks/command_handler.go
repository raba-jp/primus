// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	handlers "github.com/raba-jp/primus/pkg/operations/command/handlers"
	mock "github.com/stretchr/testify/mock"
)

// CommandHandler is an autogenerated mock type for the CommandHandler type
type CommandHandler struct {
	mock.Mock
}

type CommandHandlerRunArgs struct {
	Ctx            context.Context
	CtxAnything    bool
	Dryrun         bool
	DryrunAnything bool
	P              *handlers.CommandParams
	PAnything      bool
}

type CommandHandlerRunReturns struct {
	Err error
}

type CommandHandlerRunExpectation struct {
	Args    CommandHandlerRunArgs
	Returns CommandHandlerRunReturns
}

func (_m *CommandHandler) ApplyRunExpectation(e CommandHandlerRunExpectation) {
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
	_m.On("Run", args...).Return(e.Returns.Err)
}

func (_m *CommandHandler) ApplyRunExpectations(expectations []CommandHandlerRunExpectation) {
	for _, e := range expectations {
		_m.ApplyRunExpectation(e)
	}
}

// Run provides a mock function with given fields: ctx, dryrun, p
func (_m *CommandHandler) Run(ctx context.Context, dryrun bool, p *handlers.CommandParams) error {
	ret := _m.Called(ctx, dryrun, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, bool, *handlers.CommandParams) error); ok {
		r0 = rf(ctx, dryrun, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
