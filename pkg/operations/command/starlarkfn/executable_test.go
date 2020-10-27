package starlarkfn_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/operations/command/handlers/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/command/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
)

func TestExecutable(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.ExecutableHandlerRunExpectation
		want      lib.Value
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: return true",
			data: `v = test("data")`,
			mock: mocks.ExecutableHandlerRunExpectation{
				Args: mocks.ExecutableHandlerRunArgs{
					CtxAnything: true,
					Name:        "data",
				},
				Returns: mocks.ExecutableHandlerRunReturns{Ok: true},
			},
			want:      lib.True,
			errAssert: assert.NoError,
		},
		{
			name: "success: return false",
			data: `v = test("data")`,
			mock: mocks.ExecutableHandlerRunExpectation{
				Args: mocks.ExecutableHandlerRunArgs{
					CtxAnything: true,
					Name:        "data",
				},
				Returns: mocks.ExecutableHandlerRunReturns{
					Ok: false,
				},
			},
			want:      lib.False,
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `v = test("data", "too many")`,
			mock:      mocks.ExecutableHandlerRunExpectation{},
			want:      nil,
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executable := new(mocks.ExecutableHandler)
			executable.ApplyRunExpectation(tt.mock)

			globals, err := starlark.ExecForTest("test", tt.data, starlarkfn.Executable(executable))
			tt.errAssert(t, err)
			assert.Equal(t, globals["v"], tt.want)
		})
	}
}
