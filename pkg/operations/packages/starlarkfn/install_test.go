package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/packages/handlers"
	"github.com/raba-jp/primus/pkg/operations/packages/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/packages/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestDarwinPkgInstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.DarwinPkgInstallHandlerInstallExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="base-devel", option="option", cask=True, cmd="brew")`,
			mock: mocks.DarwinPkgInstallHandlerInstallExpectation{
				Args: mocks.DarwinPkgInstallHandlerInstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.DarwinPkgInstallParams{
						Name:   "base-devel",
						Option: "option",
						Cask:   true,
						Cmd:    "brew",
					},
				},
				Returns: mocks.DarwinPkgInstallHandlerInstallReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "option", True, "brew", "too many")`,
			mock:      mocks.DarwinPkgInstallHandlerInstallExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: package install failed",
			data: `test(name="base-devel")`,
			mock: mocks.DarwinPkgInstallHandlerInstallExpectation{
				Args: mocks.DarwinPkgInstallHandlerInstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.DarwinPkgInstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.DarwinPkgInstallHandlerInstallReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.DarwinPkgInstallHandler)
			handler.ApplyInstallExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.DarwinPkgInstall(handler))
			tt.errAssert(t, err)
		})
	}
}

func TestArchPkgInstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.ArchPkgInstallHandlerInstallExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="base-devel", option="option", cmd="yay")`,
			mock: mocks.ArchPkgInstallHandlerInstallExpectation{
				Args: mocks.ArchPkgInstallHandlerInstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.ArchPkgInstallParams{
						Name:   "base-devel",
						Option: "option",
						Cmd:    "yay",
					},
				},
				Returns: mocks.ArchPkgInstallHandlerInstallReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "option", "yay", "too many")`,
			mock:      mocks.ArchPkgInstallHandlerInstallExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: package install failed",
			data: `test(name="base-devel")`,
			mock: mocks.ArchPkgInstallHandlerInstallExpectation{
				Args: mocks.ArchPkgInstallHandlerInstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.ArchPkgInstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.ArchPkgInstallHandlerInstallReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.ArchPkgInstallHandler)
			handler.ApplyInstallExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.ArchPkgInstall(handler))
			tt.errAssert(t, err)
		})
	}
}
