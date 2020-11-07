package command_test

import (
	"testing"

	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/modules/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
)

func TestNewExecutableFunction(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		mock      mocks.OSDetectorExecutableCommandExpectation
		want      lib.Value
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: return true",
			data: `v = test("data")`,
			mock: mocks.OSDetectorExecutableCommandExpectation{
				Args: mocks.OSDetectorExecutableCommandArgs{
					CtxAnything: true,
					Name:        "data",
				},
				Returns: mocks.OSDetectorExecutableCommandReturns{Ok: true},
			},
			want:      lib.True,
			errAssert: assert.NoError,
		},
		{
			name: "success: return false",
			data: `v = test("data")`,
			mock: mocks.OSDetectorExecutableCommandExpectation{
				Args: mocks.OSDetectorExecutableCommandArgs{
					CtxAnything: true,
					Name:        "data",
				},
				Returns: mocks.OSDetectorExecutableCommandReturns{Ok: false},
			},
			want:      lib.False,
			errAssert: assert.NoError,
		},
		{
			name:      "error: too many arguments",
			data:      `v = test("data", "too many")`,
			mock:      mocks.OSDetectorExecutableCommandExpectation{},
			want:      nil,
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detector := new(mocks.OSDetector)
			detector.ApplyExecutableCommandExpectation(tt.mock)

			globals, err := starlark.ExecForTest("test", tt.data, command.NewExecutableFunction(detector))
			tt.errAssert(t, err)
			assert.Equal(t, globals["v"], tt.want)
		})
	}
}
