// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// DarwinPkgCheckInstallHandler is an autogenerated mock type for the DarwinPkgCheckInstallHandler type
type DarwinPkgCheckInstallHandler struct {
	mock.Mock
}

type DarwinPkgCheckInstallHandlerCheckInstallArgs struct {
	Ctx          context.Context
	CtxAnything  bool
	Name         string
	NameAnything bool
}

type DarwinPkgCheckInstallHandlerCheckInstallReturns struct {
	Ok bool
}

type DarwinPkgCheckInstallHandlerCheckInstallExpectation struct {
	Args    DarwinPkgCheckInstallHandlerCheckInstallArgs
	Returns DarwinPkgCheckInstallHandlerCheckInstallReturns
}

func (_m *DarwinPkgCheckInstallHandler) ApplyCheckInstallExpectation(e DarwinPkgCheckInstallHandlerCheckInstallExpectation) {
	var args []interface{}
	if e.Args.CtxAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Ctx)
	}
	if e.Args.NameAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Name)
	}
	_m.On("CheckInstall", args...).Return(e.Returns.Ok)
}

func (_m *DarwinPkgCheckInstallHandler) ApplyCheckInstallExpectations(expectations []DarwinPkgCheckInstallHandlerCheckInstallExpectation) {
	for _, e := range expectations {
		_m.ApplyCheckInstallExpectation(e)
	}
}

// CheckInstall provides a mock function with given fields: ctx, name
func (_m *DarwinPkgCheckInstallHandler) CheckInstall(ctx context.Context, name string) bool {
	ret := _m.Called(ctx, name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
