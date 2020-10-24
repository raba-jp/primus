package starlarkfn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/packages/handlers"
	"github.com/raba-jp/primus/pkg/operations/packages/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"golang.org/x/xerrors"
)

func TestDarwinPkgUninstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      handlers.DarwinPkgUninstallHandlerUninstallExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="base-devel")`,
			mock: handlers.DarwinPkgUninstallHandlerUninstallExpectation{
				Args: handlers.DarwinPkgUninstallHandlerUninstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.DarwinPkgUninstallParams{
						Name: "base-devel",
					},
				},
				Returns: handlers.DarwinPkgUninstallHandlerUninstallReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "too many")`,
			mock:      handlers.DarwinPkgUninstallHandlerUninstallExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to uninstall",
			data: `test(name="base-devel")`,
			mock: handlers.DarwinPkgUninstallHandlerUninstallExpectation{
				Args: handlers.DarwinPkgUninstallHandlerUninstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.DarwinPkgUninstallParams{
						Name: "base-devel",
					},
				},
				Returns: handlers.DarwinPkgUninstallHandlerUninstallReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(handlers.MockDarwinPkgUninstallHandler)
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
		mock   handlers.ArchPkgUninstallHandlerUninstallExpectation
		hasErr bool
	}{
		{
			name: "success",
			data: `test(name="base-devel")`,
			mock: handlers.ArchPkgUninstallHandlerUninstallExpectation{
				Args: handlers.ArchPkgUninstallHandlerUninstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.ArchPkgUninstallParams{
						Name: "base-devel",
					},
				},
				Returns: handlers.ArchPkgUninstallHandlerUninstallReturns{
					Err: nil,
				},
			},
			hasErr: false,
		},
		{
			name:   "error: too many arguments",
			data:   `test("base-devel", "yay", "too many")`,
			mock:   handlers.ArchPkgUninstallHandlerUninstallExpectation{},
			hasErr: true,
		},
		{
			name: "error: failed to uninstall",
			data: `test(name="base-devel")`,
			mock: handlers.ArchPkgUninstallHandlerUninstallExpectation{
				Args: handlers.ArchPkgUninstallHandlerUninstallArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.ArchPkgUninstallParams{
						Name: "base-devel",
					},
				},
				Returns: handlers.ArchPkgUninstallHandlerUninstallReturns{
					Err: xerrors.New("dummy"),
				},
			},
			hasErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(handlers.MockArchPkgUninstallHandler)
			handler.ApplyUninstallExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.ArchPkgUninstall(handler))
			if !tt.hasErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
