package command_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/exec"
	"golang.org/x/xerrors"

	"github.com/raba-jp/primus/pkg/functions/command"

	"github.com/stretchr/testify/assert"

	"github.com/raba-jp/primus/pkg/starlark"
	lib "go.starlark.net/starlark"
)

func TestNewExecutableFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      func(ctx context.Context, name string) bool
		want      lib.Value
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: return true",
			data: `v = test("data")`,
			mock: func(ctx context.Context, name string) bool {
				return true
			},
			want:      lib.True,
			errAssert: assert.NoError,
		},
		{
			name: "success: return false",
			data: `v = test("data")`,
			mock: func(ctx context.Context, name string) bool {
				return false
			},
			want:      lib.False,
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `v = test("data", "too many")`,
			mock: func(ctx context.Context, name string) bool {
				return true
			},
			want:      nil,
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			globals, err := starlark.ExecForTest("test", tt.data, command.NewExecutableFunction(tt.mock))
			tt.errAssert(t, err)
			assert.Equal(t, globals["v"], tt.want)
		})
	}
}

func TestExecutable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data string
		mock exec.InterfaceLookPathExpectation
		want bool
	}{
		{
			name: "success",
			data: "cat",
			mock: exec.InterfaceLookPathExpectation{
				Args: exec.InterfaceLookPathArgs{
					File: "cat",
				},
				Returns: exec.InterfaceLookPathReturns{
					Path: "/bin/cat",
					Err:  nil,
				},
			},
			want: true,
		},
		{
			name: "success: not found",
			data: "cat",
			mock: exec.InterfaceLookPathExpectation{
				Args: exec.InterfaceLookPathArgs{
					File: "cat",
				},
				Returns: exec.InterfaceLookPathReturns{
					Path: "",
					Err:  xerrors.New("dummy"),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			exc := new(exec.MockInterface)
			exc.ApplyLookPathExpectation(tt.mock)

			ret := command.Executable(exc)(context.Background(), tt.data)
			assert.Equal(t, tt.want, ret)
		})
	}
}
