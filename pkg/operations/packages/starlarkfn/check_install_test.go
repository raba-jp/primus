package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/packages/handlers"
	"github.com/raba-jp/primus/pkg/operations/packages/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func TestDarwinPkgCheckInstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      handlers.DarwinPkgCheckInstallHandlerCheckInstallExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: true",
			data: `test(name="base-devel")`,
			mock: handlers.DarwinPkgCheckInstallHandlerCheckInstallExpectation{
				Args: handlers.DarwinPkgCheckInstallHandlerCheckInstallArgs{
					CtxAnything: true,
					Name:        "base-devel",
				},
				Returns: handlers.DarwinPkgCheckInstallHandlerCheckInstallReturns{
					Ok: true,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: false",
			data: `test(name="base-devel")`,
			mock: handlers.DarwinPkgCheckInstallHandlerCheckInstallExpectation{
				Args: handlers.DarwinPkgCheckInstallHandlerCheckInstallArgs{
					CtxAnything: true,
					Name:        "base-devel",
				},
				Returns: handlers.DarwinPkgCheckInstallHandlerCheckInstallReturns{
					Ok: false,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "too many")`,
			mock:      handlers.DarwinPkgCheckInstallHandlerCheckInstallExpectation{},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(handlers.MockDarwinPkgCheckInstallHandler)
			handler.ApplyCheckInstallExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.DarwinPkgCheckInstall(handler))
			tt.errAssert(t, err)
		})
	}
}

func TestArchPkgCheckInstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      handlers.ArchPkgCheckInstallHandlerCheckInstallExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: true",
			data: `test(name="base-devel")`,
			mock: handlers.ArchPkgCheckInstallHandlerCheckInstallExpectation{
				Args: handlers.ArchPkgCheckInstallHandlerCheckInstallArgs{
					CtxAnything: true,
					Name:        "base-devel",
				},
				Returns: handlers.ArchPkgCheckInstallHandlerCheckInstallReturns{
					Ok: true,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: false",
			data: `test(name="base-devel")`,
			mock: handlers.ArchPkgCheckInstallHandlerCheckInstallExpectation{
				Args: handlers.ArchPkgCheckInstallHandlerCheckInstallArgs{
					CtxAnything: true,
					Name:        "base-devel",
				},
				Returns: handlers.ArchPkgCheckInstallHandlerCheckInstallReturns{
					Ok: false,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "too many")`,
			mock:      handlers.ArchPkgCheckInstallHandlerCheckInstallExpectation{},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(handlers.MockArchPkgCheckInstallHandler)
			handler.ApplyCheckInstallExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.ArchPkgCheckInstall(handler))
			tt.errAssert(t, err)
		})
	}
}
