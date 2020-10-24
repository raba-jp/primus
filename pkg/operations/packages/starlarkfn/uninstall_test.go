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

func TestDarwinPkgUninstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.DarwinPkgUninstallHandlerUninstallExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="base-devel")`,
			mock: mocks.DarwinPkgUninstallHandlerUninstallExpectation{
				Args: mocks.DarwinPkgUninstallHandlerUninstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.DarwinPkgUninstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.DarwinPkgUninstallHandlerUninstallReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "too many")`,
			mock:      mocks.DarwinPkgUninstallHandlerUninstallExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to uninstall",
			data: `test(name="base-devel")`,
			mock: mocks.DarwinPkgUninstallHandlerUninstallExpectation{
				Args: mocks.DarwinPkgUninstallHandlerUninstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.DarwinPkgUninstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.DarwinPkgUninstallHandlerUninstallReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.DarwinPkgUninstallHandler)
			handler.ApplyUninstallExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.DarwinPkgUninstall(handler))
			tt.errAssert(t, err)
		})
	}
}

func TestArchPkgUninstall(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		mock   mocks.ArchPkgUninstallHandlerUninstallExpectation
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="base-devel")`,
			mock: mocks.ArchPkgUninstallHandlerUninstallExpectation{
				Args: mocks.ArchPkgUninstallHandlerUninstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.ArchPkgUninstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.ArchPkgUninstallHandlerUninstallReturns{
					Err: nil,
				},
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("base-devel", "yay", "too many")`,
			mock:   mocks.ArchPkgUninstallHandlerUninstallExpectation{},
			hasErr: true,
		},
		{
			name: "error: failed to uninstall",
			data: `test(name="base-devel")`,
			mock: mocks.ArchPkgUninstallHandlerUninstallExpectation{
				Args: mocks.ArchPkgUninstallHandlerUninstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.ArchPkgUninstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.ArchPkgUninstallHandlerUninstallReturns{
					Err: xerrors.New("dummy"),
				},
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.ArchPkgUninstallHandler)
			handler.ApplyUninstallExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.ArchPkgUninstall(handler))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
