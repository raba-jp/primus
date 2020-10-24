package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/packages/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/packages/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
)

func TestDarwinPkgCheckInstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.DarwinPkgCheckInstallHandlerCheckInstallExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: true",
			data: `test(name="base-devel")`,
			mock: mocks.DarwinPkgCheckInstallHandlerCheckInstallExpectation{
				Args: mocks.DarwinPkgCheckInstallHandlerCheckInstallArgs{
					CtxAnything: true,
					Name:        "base-devel",
				},
				Returns: mocks.DarwinPkgCheckInstallHandlerCheckInstallReturns{
					Ok: true,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: false",
			data: `test(name="base-devel")`,
			mock: mocks.DarwinPkgCheckInstallHandlerCheckInstallExpectation{
				Args: mocks.DarwinPkgCheckInstallHandlerCheckInstallArgs{
					CtxAnything: true,
					Name:        "base-devel",
				},
				Returns: mocks.DarwinPkgCheckInstallHandlerCheckInstallReturns{
					Ok: false,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "too many")`,
			mock:      mocks.DarwinPkgCheckInstallHandlerCheckInstallExpectation{},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.DarwinPkgCheckInstallHandler)
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
		mock      mocks.ArchPkgCheckInstallHandlerCheckInstallExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: true",
			data: `test(name="base-devel")`,
			mock: mocks.ArchPkgCheckInstallHandlerCheckInstallExpectation{
				Args: mocks.ArchPkgCheckInstallHandlerCheckInstallArgs{
					CtxAnything: true,
					Name:        "base-devel",
				},
				Returns: mocks.ArchPkgCheckInstallHandlerCheckInstallReturns{
					Ok: true,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: false",
			data: `test(name="base-devel")`,
			mock: mocks.ArchPkgCheckInstallHandlerCheckInstallExpectation{
				Args: mocks.ArchPkgCheckInstallHandlerCheckInstallArgs{
					CtxAnything: true,
					Name:        "base-devel",
				},
				Returns: mocks.ArchPkgCheckInstallHandlerCheckInstallReturns{
					Ok: false,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "too many")`,
			mock:      mocks.ArchPkgCheckInstallHandlerCheckInstallExpectation{},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.ArchPkgCheckInstallHandler)
			handler.ApplyCheckInstallExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.ArchPkgCheckInstall(handler))
			tt.errAssert(t, err)
		})
	}
}
