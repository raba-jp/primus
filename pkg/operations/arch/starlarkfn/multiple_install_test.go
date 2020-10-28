package starlarkfn_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/operations/arch/handlers"
	"github.com/raba-jp/primus/pkg/operations/arch/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/arch/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestMultipleInstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.MultipleInstallHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(names=["arg1", "arg2"])`,
			mock: mocks.MultipleInstallHandlerRunExpectation{
				Args: mocks.MultipleInstallHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.MultipleInstallParams{
						Names: []string{"arg1", "arg2"},
					},
				},
				Returns: mocks.MultipleInstallHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success, includes number",
			data: `test(names=["arg1", 1, "arg2"])`,
			mock: mocks.MultipleInstallHandlerRunExpectation{
				Args: mocks.MultipleInstallHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.MultipleInstallParams{
						Names: []string{"arg1", "arg2"},
					},
				},
				Returns: mocks.MultipleInstallHandlerRunReturns{
					Err: nil,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test(["arg1", "arg2"], "too many")`,
			mock:      mocks.MultipleInstallHandlerRunExpectation{},
			errAssert: assert.Error,
		},
		{
			name: "error: package install failed",
			data: `test(names=["arg1", "arg2"])`,
			mock: mocks.MultipleInstallHandlerRunExpectation{
				Args: mocks.MultipleInstallHandlerRunArgs{
					CtxAnything:    true,
					DryrunAnything: true,
					P: &handlers.MultipleInstallParams{
						Names: []string{"arg1", "arg2"},
					},
				},
				Returns: mocks.MultipleInstallHandlerRunReturns{
					Err: xerrors.New("dummy"),
				},
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.MultipleInstallHandler)
			handler.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.MultipleInstall(handler))
			tt.errAssert(t, err)
		})
	}
}
