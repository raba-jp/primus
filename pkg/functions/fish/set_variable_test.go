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

func TestNewSetVariableFunction(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		data      string
		mock      fish.SetVariableRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			data: `test(name="GOPATH", value="$HOME/go", scope="universal", export=True)`,
			mock: func(ctx context.Context, p *fish.SetVariableParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: args",
			data: `test("GOPATH", "$HOME/go", "universal", True)`,
			mock: func(ctx context.Context, p *fish.SetVariableParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: global scope",
			data: `test("GOPATH", "$HOME/go", "global", True)`,
			mock: func(ctx context.Context, p *fish.SetVariableParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: local scope",
			data: `test("GOPATH", "$HOME/go", "local", True)`,
			mock: func(ctx context.Context, p *fish.SetVariableParams) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error: unexpected scope",
			data: `test(name="GOPATH", value="$HOME/go", scope="dummy", export=True)`,
			mock: func(ctx context.Context, p *fish.SetVariableParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: too many arguments",
			data: `test("GOPATH", "$HOME/go", "universal", True, "too many")`,
			mock: func(ctx context.Context, p *fish.SetVariableParams) error {
				return nil
			},
			errAssert: assert.Error,
		},
		{
			name: "error: return handler error",
			data: `test("GOPATH", "$HOME/go", "universal", True)`,
			mock: func(ctx context.Context, p *fish.SetVariableParams) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			_, err := starlark.ExecForTest("test", tt.data, fish.NewSetVariableFunction(tt.mock))
			tt.errAssert(t, err)
		})
	}
}

func TestSetVariable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		params    *fish.SetVariableParams
		mock      command.ExecuteRunner
		errAssert assert.ErrorAssertionFunc
	}{
		{
			name: "success: scope universal",
			params: &fish.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  fish.UniversalScope,
				Export: true,
			},
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: scope global",
			params: &fish.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  fish.GlobalScope,
				Export: true,
			},
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: scope local",
			params: &fish.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  fish.LocalScope,
				Export: true,
			},
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "success: no export",
			params: &fish.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  fish.LocalScope,
				Export: false,
			},
			mock: func(ctx context.Context, p *command.Params) error {
				return nil
			},
			errAssert: assert.NoError,
		},
		{
			name: "error",
			params: &fish.SetVariableParams{
				Name:   "GOPATH",
				Value:  "$HOME/go",
				Scope:  fish.UniversalScope,
				Export: true,
			},
			mock: func(ctx context.Context, p *command.Params) error {
				return xerrors.New("dummy")
			},
			errAssert: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := fish.SetVariable(tt.mock)(context.Background(), tt.params)
			tt.errAssert(t, err)
		})
	}
}
