package starlarkfn_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/operations/darwin/handlers/mocks"
	"github.com/raba-jp/primus/pkg/operations/darwin/starlarkfn"

	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/stretchr/testify/assert"
)

func TestCheckInstall(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.CheckInstallHandlerRunExpectation
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: true",
			data: `test(name="base-devel")`,
			mock: mocks.CheckInstallHandlerRunExpectation{
				Args: mocks.CheckInstallHandlerRunArgs{
					CtxAnything: true,
					Name:        "base-devel",
				},
				Returns: mocks.CheckInstallHandlerRunReturns{
					Ok: true,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: false",
			data: `test(name="base-devel")`,
			mock: mocks.CheckInstallHandlerRunExpectation{
				Args: mocks.CheckInstallHandlerRunArgs{
					CtxAnything: true,
					Name:        "base-devel",
				},
				Returns: mocks.CheckInstallHandlerRunReturns{
					Ok: false,
				},
			},
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `test("base-devel", "too many")`,
			mock:      mocks.CheckInstallHandlerRunExpectation{},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(mocks.CheckInstallHandler)
			handler.ApplyRunExpectation(tt.mock)

			_, err := starlark.ExecForTest("test", tt.data, starlarkfn.CheckInstall(handler))
			tt.errAssert(t, err)
		})
	}
}
