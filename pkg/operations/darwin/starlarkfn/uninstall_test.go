package starlarkfn_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/operations/darwin/handlers"
	"github.com/raba-jp/primus/pkg/operations/darwin/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/darwin/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestUninstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.UninstallHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="base-devel")`,
			mock: mocks.UninstallHandlerRunExpectation{
				Args: mocks.UninstallHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.UninstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.UninstallHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "too many")`,
			mock:      mocks.UninstallHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: failed to uninstall",
			data: `test(name="base-devel")`,
			mock: mocks.UninstallHandlerRunExpectation{
				Args: mocks.UninstallHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.UninstallParams{
						Name: "base-devel",
					},
				},
				Returns: mocks.UninstallHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.UninstallHandler)
			handler.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.Uninstall(handler))
			tt.errAssert(t, err)
		})
	}
}
