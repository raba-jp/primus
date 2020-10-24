// Code generated by mockery v1.0.0. DO NOT EDIT.

package handlers

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockArchPkgCheckInstallHandler is an autogenerated mock type for the ArchPkgCheckInstallHandler type
type MockArchPkgCheckInstallHandler struct {
	mock.Mock
}

type ArchPkgCheckInstallHandlerCheckInstallArgs struct {
	Ctx          context.Context
	CtxAnything  bool
	Name         string
	NameAnything bool
}

type ArchPkgCheckInstallHandlerCheckInstallReturns struct {
	Ok bool
}

type ArchPkgCheckInstallHandlerCheckInstallExpectation struct {
	Args    ArchPkgCheckInstallHandlerCheckInstallArgs
	Returns ArchPkgCheckInstallHandlerCheckInstallReturns
}

func (_m *MockArchPkgCheckInstallHandler) ApplyCheckInstallExpectation(e ArchPkgCheckInstallHandlerCheckInstallExpectation) {
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

func (_m *MockArchPkgCheckInstallHandler) ApplyCheckInstallExpectations(expectations []ArchPkgCheckInstallHandlerCheckInstallExpectation) {
	for _, e := range expectations {
		_m.ApplyCheckInstallExpectation(e)
	}
}

// CheckInstall provides a mock function with given fields: ctx, name
func (_m *MockArchPkgCheckInstallHandler) CheckInstall(ctx context.Context, name string) bool {
	ret := _m.Called(ctx, name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}
