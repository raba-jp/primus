package fish_test

import (
	"context"
	"testing"

	"github.com/raba-jp/primus/pkg/functions/command"
	"github.com/raba-jp/primus/pkg/functions/fish"
	"github.com/raba-jp/primus/pkg/starlark"
	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestNewSetPathFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      fish.SetPathRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(values=["$GOPATH/bin", "$HOME/.bin"])`,
			mock: func(ctx context.Context, p *fish.SetPathParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: args",
			data: `test(["$GOPATH/bin", "$HOME/.bin"])`,
			mock: func(ctx context.Context, p *fish.SetPathParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: include int and bool",
			data: `test(["$GOPATH/bin", 1, True, "$HOME/.bin"])`,
			mock: func(ctx context.Context, p *fish.SetPathParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: too many arguments",
			data: `test(["$GOPATH/bin", "$HOME/.bin"], "too many")`,
			mock: func(ctx context.Context, p *fish.SetPathParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: return handler error",
			data: `test(["$GOPATH/bin", "$HOME/.bin"])`,
			mock: func(ctx context.Context, p *fish.SetPathParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, fish.NewSetPathFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestSetPath(t *testing.T) {
	tests := []struct {
		name      string
		params    *fish.SetPathParams
		mock      command.ExecuteRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			params: &fish.SetPathParams{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error",
			params: &fish.SetPathParams{
				Values: []string{"$GOPATH/bin", "$HOME/.bin"},
			},
			mock: func(ctx context.Context, p *command.Params) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := fish.SetPath(tt.mock)(context.Background(), tt.params)
			tt.errAssert(t, err)
		})
	}
}
