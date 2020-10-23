package starlarkfn_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/operations/command/handlers"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/operations/command/starlarkfn"
	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
)

func TestExecutable(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      handlers.ExecutableHandlerExecutableExpectation
		want      lib.Value
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: return true",
			data: `v = test("data")`,
			mock: handlers.ExecutableHandlerExecutableExpectation{
				Args: handlers.ExecutableHandlerExecutableArgs{
					CtxAnything: true,
					Name:        "data",
				},
				Returns: handlers.ExecutableHandlerExecutableReturns{Ok: true},
			},
			want:      lib.True,
			errAssert: assert.NoError,
		},
		{
			name: "success: return false",
			data: `v = test("data")`,
			mock: handlers.ExecutableHandlerExecutableExpectation{
				Args: handlers.ExecutableHandlerExecutableArgs{
					CtxAnything: true,
					Name:        "data",
				},
				Returns: handlers.ExecutableHandlerExecutableReturns{
					Ok: false,
				},
			},
			want:      lib.False,
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `v = test("data", "too many")`,
			mock:      handlers.ExecutableHandlerExecutableExpectation{},
			want:      nil,
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := new(handlers.MockExecutableHandler)
			handler.ApplyExecutableExpectation(tt.mock)

			globals, err := starlark.ExecForTest("test", tt.data, starlarkfn.Executable(handler))
			tt.errAssert(t, err)
			assert.Equal(t, globals["v"], tt.want)
		})
	}
}
