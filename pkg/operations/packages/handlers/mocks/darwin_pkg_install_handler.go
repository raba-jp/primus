// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	handlers "github.com/raba-jp/primus/pkg/operations/packages/handlers"
	mock "github.com/stretchr/testify/mock"
)

// DarwinPkgInstallHandler is an autogenerated mock type for the DarwinPkgInstallHandler type
type DarwinPkgInstallHandler struct {
	mock.Mock
}

type DarwinPkgInstallHandlerInstallArgs struct {
	Ctx            context.Context
	CtxAnything    bool
	Dryrun         bool
	DryrunAnything bool
	P              *handlers.DarwinPkgInstallParams
	PAnything      bool
}

type DarwinPkgInstallHandlerInstallReturns struct {
	Err error
}

type DarwinPkgInstallHandlerInstallExpectation struct {
	Args    DarwinPkgInstallHandlerInstallArgs
	Returns DarwinPkgInstallHandlerInstallReturns
}

func (_m *DarwinPkgInstallHandler) ApplyInstallExpectation(e DarwinPkgInstallHandlerInstallExpectation) {
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
	_m.On("Install", args...).Return(e.Returns.Err)
}

func (_m *DarwinPkgInstallHandler) ApplyInstallExpectations(expectations []DarwinPkgInstallHandlerInstallExpectation) {
	for _, e := range expectations {
		_m.ApplyInstallExpectation(e)
	}
}

// Install provides a mock function with given fields: ctx, dryrun, p
func (_m *DarwinPkgInstallHandler) Install(ctx context.Context, dryrun bool, p *handlers.DarwinPkgInstallParams) error {
	ret := _m.Called(ctx, dryrun, p)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, bool, *handlers.DarwinPkgInstallParams) error); ok {
		r0 = rf(ctx, dryrun, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
