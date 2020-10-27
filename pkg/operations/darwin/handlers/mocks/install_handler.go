// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	handlers "github.com/raba-jp/primus/pkg/operations/darwin/handlers"
	mock "github.com/stretchr/testify/mock"
)

// InstallHandler is an autogenerated mock type for the InstallHandler type
type InstallHandler struct {
	mock.Mock
}

type InstallHandlerRunArgs struct {
	Ctx            context.Context
	CtxAnything    bool
	Dryrun         bool
	DryrunAnything bool
	P              *handlers.InstallParams
	PAnything      bool
}

type InstallHandlerRunReturns struct {
	Err error
}

type InstallHandlerRunExpectation struct {
	Args    InstallHandlerRunArgs
	Returns InstallHandlerRunReturns
}

func (_m *InstallHandler) ApplyRunExpectation(e InstallHandlerRunExpectation) {
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

func (_m *InstallHandler) ApplyRunExpectations(expectations []InstallHandlerRunExpectation) {
	for _, e := range expectations {
		_m.ApplyRunExpectation(e)
	}
}

// Run provides a mock function with given fields: ctx, dryrun, p
func (_m *InstallHandler) Run(ctx context.Context, dryrun bool, p *handlers.InstallParams) error {
	ret := _m.Called(ctx, dryrun, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, bool, *handlers.InstallParams) error); ok {
		r0 = rf(ctx, dryrun, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
