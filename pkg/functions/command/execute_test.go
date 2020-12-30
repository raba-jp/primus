package command_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/backend"
	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestNewExecuteFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      backend.Execute
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: string array kwargs",
			data: `test(cmd="echo", args=["hello", "world"])`,
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: int kwargs",
			data: `test(cmd="echo", args=[1])`,
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: bigint kwargs",
			data: `test(cmd="echo", args=[9007199254740991])`,
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "success: bool kwargs",
			data: `test(cmd="echo", args=[False, True])`,
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success(unsupported): float kwargs",
			data: `test(cmd="echo", args=[1.111])`,
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "success: no args",
			data: `test("echo")`,
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: with user and cwd",
			data: `test("echo", [], user="testuser", cwd="/home/testuser")`,
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test("echo", [], "testuser", "/home/testuser", "too many")`,
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: execute command failed",
			data: `test("echo")`,
			mock: func(ctx context.Context, p *backend.ExecuteParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, command.NewExecuteFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}
